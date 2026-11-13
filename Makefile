# Jalan pintas dari folder utama: make run, make migrate-up, make test, dll.
# Semua perintah diteruskan ke backend/Makefile.
.PHONY: run build test tidy migrate-up migrate-down migrate-status migrate-reset db

run build test tidy migrate-up migrate-down migrate-status migrate-reset:
	$(MAKE) -C backend $@

# Nyalakan database Postgres di Docker
db:
	docker compose up -d db
