-- +goose Up
-- Login web admin juga lewat Telegram/SMS (MFA menyusul): izinkan audience 'admin'.
ALTER TABLE telegram_login_requests DROP CONSTRAINT telegram_login_requests_audience_check;
ALTER TABLE telegram_login_requests ADD CONSTRAINT telegram_login_requests_audience_check
  CHECK (audience IN ('app', 'seller_web', 'admin'));

-- +goose Down
DELETE FROM telegram_login_requests WHERE audience = 'admin';
ALTER TABLE telegram_login_requests DROP CONSTRAINT telegram_login_requests_audience_check;
ALTER TABLE telegram_login_requests ADD CONSTRAINT telegram_login_requests_audience_check
  CHECK (audience IN ('app', 'seller_web'));
