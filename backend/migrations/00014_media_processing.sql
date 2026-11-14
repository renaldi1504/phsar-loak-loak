-- +goose Up
-- Kompres video (H.264 ≤1280 px, faststart) + gambar sampul, dan perkecil foto (≤1600 px, tanpa EXIF/GPS).
ALTER TABLE media_processing_jobs ADD COLUMN run_after timestamptz;
CREATE INDEX media_processing_jobs_queue ON media_processing_jobs (created_at) WHERE status = 'queued';
ALTER TABLE product_media ADD COLUMN original_bytes bigint;   -- ukuran file asli dari HP
ALTER TABLE product_media ADD COLUMN processed_bytes bigint;  -- ukuran setelah kompres

-- +goose Down
ALTER TABLE product_media DROP COLUMN processed_bytes;
ALTER TABLE product_media DROP COLUMN original_bytes;
DROP INDEX media_processing_jobs_queue;
ALTER TABLE media_processing_jobs DROP COLUMN run_after;
