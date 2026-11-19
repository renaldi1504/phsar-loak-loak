FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/phsar ./cmd/phsar

# Runtime dengan ffmpeg (kompres video + thumbnail, perkecil foto).
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg ca-certificates tzdata curl \
    && rm -rf /var/lib/apt/lists/* && useradd --system --uid 10001 --create-home phsar
COPY --from=build /out/phsar /usr/local/bin/phsar
USER phsar
WORKDIR /home/phsar
RUN mkdir -p /home/phsar/data/uploads
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 CMD curl -fsS http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT ["phsar"]
CMD ["serve"]
