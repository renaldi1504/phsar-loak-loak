# Layanan cari-dengan-foto (CLIP ViT-B/32 ONNX). Model diunduh saat build supaya container siap tanpa internet.
FROM python:3.12-slim
RUN apt-get update && apt-get install -y --no-install-recommends curl ca-certificates \
    && rm -rf /var/lib/apt/lists/* && useradd --system --uid 10002 --create-home embedder
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY download-model.sh server.py ./
RUN bash download-model.sh && chown -R embedder /app
USER embedder
ENV EMBEDDER_ADDR=0.0.0.0:8090 EMBEDDER_THREADS=2
EXPOSE 8090
CMD ["python3", "server.py"]
