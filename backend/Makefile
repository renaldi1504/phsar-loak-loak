# Jalankan dari folder backend/. Variabel DATABASE_URL dibaca dari ../.env bila ada.
-include ../.env
export

.PHONY: media-backfill translate-test tidy run build test qa migrate-up migrate-down migrate-status migrate-reset

tidy:
	go mod tidy

run:
	go run ./cmd/phsar serve

build:
	CGO_ENABLED=0 go build -o bin/phsar ./cmd/phsar

test:
	go test ./...

migrate-up:
	go run ./cmd/phsar migrate up

migrate-down:
	go run ./cmd/phsar migrate down

migrate-status:
	go run ./cmd/phsar migrate status

migrate-reset:
	go run ./cmd/phsar migrate reset

qa:
	../scripts/qa.sh

# Cek API key & model terjemahan AI (satu kalimat contoh → km & zh)
translate-test:
	go run ./cmd/phsar translate-test

# Antrekan video/foto lama (sebelum ffmpeg terpasang) supaya dikompres
media-backfill:
	go run ./cmd/phsar media-backfill
