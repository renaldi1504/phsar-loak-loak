"""Phsar embedder — layanan kecil untuk cari-dengan-foto.

Menerima satu gambar (POST /embed, body = byte gambar) dan mengembalikan vektor CLIP 512 dimensi
yang sudah dinormalisasi. Backend Go yang menyimpan vektor & mencari produk mirip; layanan ini
tidak menyimpan apa pun dan tidak perlu akses database.

Jalankan:  python3 embedder/server.py   (env: EMBEDDER_MODEL, EMBEDDER_ADDR, EMBEDDER_THREADS)
"""

import io
import json
import os
import sys
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

import numpy as np
import onnxruntime as ort
from PIL import Image, ImageOps

MODEL = os.environ.get("EMBEDDER_MODEL", os.path.join(os.path.dirname(__file__), "models", "clip-vit-b32-q8.onnx"))
ADDR = os.environ.get("EMBEDDER_ADDR", "127.0.0.1:8090")
THREADS = int(os.environ.get("EMBEDDER_THREADS", "2"))
MAX_BYTES = 12 * 1024 * 1024
MODEL_NAME = "clip-vit-b32-q8"

MEAN = np.array([0.48145466, 0.4578275, 0.40821073], dtype=np.float32).reshape(3, 1, 1)
STD = np.array([0.26862954, 0.26130258, 0.27577711], dtype=np.float32).reshape(3, 1, 1)

opts = ort.SessionOptions()
opts.intra_op_num_threads = THREADS
opts.inter_op_num_threads = 1
session = ort.InferenceSession(MODEL, sess_options=opts, providers=["CPUExecutionProvider"])
lock = threading.Lock()  # satu inferensi sekaligus: hemat RAM di server kecil
DUMMY_IDS = np.array([[49406, 49407]], dtype=np.int64)  # teks kosong (<start><end>), hanya agar model gabungan mau jalan
DUMMY_MASK = np.ones_like(DUMMY_IDS)


def preprocess(data: bytes) -> np.ndarray:
    img = Image.open(io.BytesIO(data))
    img = ImageOps.exif_transpose(img)  # foto HP sering berputar
    if img.mode in ("RGBA", "LA", "P"):
        img = img.convert("RGBA")
        bg = Image.new("RGBA", img.size, (255, 255, 255, 255))
        img = Image.alpha_composite(bg, img)
    img = img.convert("RGB")
    w, h = img.size
    s = 224 / min(w, h)
    img = img.resize((max(224, round(w * s)), max(224, round(h * s))), Image.BICUBIC)
    w, h = img.size
    left, top = (w - 224) // 2, (h - 224) // 2
    img = img.crop((left, top, left + 224, top + 224))
    x = np.asarray(img, dtype=np.float32).transpose(2, 0, 1) / 255.0
    return ((x - MEAN) / STD)[None].astype(np.float32)


def embed(data: bytes) -> list:
    px = preprocess(data)
    with lock:
        out = session.run(["image_embeds"], {"input_ids": DUMMY_IDS, "attention_mask": DUMMY_MASK, "pixel_values": px})[0][0]
    n = float(np.linalg.norm(out)) or 1.0
    return [round(float(v) / n, 6) for v in out]


class Handler(BaseHTTPRequestHandler):
    server_version = "phsar-embedder"

    def log_message(self, fmt, *args):  # log singkat ke stderr
        sys.stderr.write("embedder: " + (fmt % args) + "\n")

    def _json(self, code, obj):
        body = json.dumps(obj).encode()
        self.send_response(code)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def do_GET(self):
        if self.path == "/health":
            return self._json(200, {"ok": True, "model": MODEL_NAME, "dim": 512})
        self._json(404, {"error": "not_found"})

    def do_POST(self):
        if self.path != "/embed":
            return self._json(404, {"error": "not_found"})
        n = int(self.headers.get("Content-Length") or 0)
        if n <= 0:
            return self._json(400, {"error": "empty"})
        if n > MAX_BYTES:
            return self._json(413, {"error": "too_large"})
        data = self.rfile.read(n)
        try:
            vec = embed(data)
        except Exception as e:  # bukan gambar / rusak
            return self._json(422, {"error": "bad_image", "detail": str(e)[:200]})
        self._json(200, {"model": MODEL_NAME, "vector": vec})


if __name__ == "__main__":
    host, port = ADDR.rsplit(":", 1)
    srv = ThreadingHTTPServer((host, int(port)), Handler)
    sys.stderr.write(f"embedder: {MODEL_NAME} siap di http://{ADDR}\n")
    srv.serve_forever()
