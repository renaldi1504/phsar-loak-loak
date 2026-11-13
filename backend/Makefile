# Jalankan dari folder backend/. Variabel DATABASE_URL dibaca dari ../.env bila ada.
-include ../.env
export

.PHONY: tidy run build test migrate-up migrate-down migrate-status migrate-reset

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
