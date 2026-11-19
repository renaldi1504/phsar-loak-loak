#!/usr/bin/env bash
# Unduh model CLIP (ViT-B/32, kuantisasi 8-bit, ±150 MB) untuk cari-dengan-foto.
# Sumber utama: Hugging Face (Xenova/clip-vit-base-patch32). Cadangan: repo GitHub website Phsar pertama (Git LFS).
set -euo pipefail
cd "$(dirname "$0")"
mkdir -p models
OUT=models/clip-vit-b32-q8.onnx
if [ -s "$OUT" ] && [ "$(wc -c < "$OUT")" -gt 100000000 ]; then
  echo "Model sudah ada: $OUT"; exit 0
fi
for URL in \
  "https://huggingface.co/Xenova/clip-vit-base-patch32/resolve/main/onnx/model_quantized.onnx" \
  "https://media.githubusercontent.com/media/xudodigital/phsar/main/assets/models/clip-vit-base-patch32/onnx/model_quantized.onnx"; do
  echo "Mengunduh dari $URL …"
  if curl -fL --retry 2 -o "$OUT.part" "$URL" && [ "$(wc -c < "$OUT.part")" -gt 100000000 ]; then
    mv "$OUT.part" "$OUT"; echo "✅ Model tersimpan di embedder/$OUT"; exit 0
  fi
done
rm -f "$OUT.part"
echo "❌ Gagal mengunduh model." >&2
exit 1
