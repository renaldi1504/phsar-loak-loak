# Jalan pintas dari folder utama: make run, make migrate-up, make test, dll.
# Semua perintah diteruskan ke backend/Makefile.
.PHONY: media-backfill translate-test run build test tidy migrate-up migrate-down migrate-status migrate-reset db qa

run build test tidy migrate-up migrate-down migrate-status migrate-reset:
	$(MAKE) -C backend $@

# Nyalakan database Postgres di Docker
db:
	docker compose up -d db

# QA otomatis: format, vet, build, unit test, tes integrasi semua alur, JS, cek rahasia (lihat QA.md)
qa:
	./scripts/qa.sh

# Cek API key & model terjemahan AI (satu kalimat contoh → km & zh)
translate-test:
	$(MAKE) -C backend translate-test

# Antrekan video/foto lama (sebelum ffmpeg terpasang) supaya dikompres
media-backfill:
	$(MAKE) -C backend media-backfill
