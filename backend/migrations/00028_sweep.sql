-- +goose Up
-- Commission sweep: catatan pemindahan komisi Phsar dari rekening penampung PayWay ke rekening perusahaan.
-- Tidak pernah dihapus; salah catat → dibatalkan (void) dengan alasan.
CREATE TABLE commission_sweeps (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  amount      numeric(12,2) NOT NULL CHECK (amount > 0),
  to_account  varchar(120) NOT NULL,                 -- mis. "ABA Phsar Co. ••• 1234"
  reference   varchar(80) NOT NULL,                  -- nomor transfer bank
  swept_on    date NOT NULL,
  note        varchar(300),
  created_by  uuid NOT NULL REFERENCES users(id),
  created_at  timestamptz NOT NULL DEFAULT now(),
  voided_at   timestamptz,
  voided_by   uuid REFERENCES users(id),
  void_note   varchar(300),
  CHECK ((voided_at IS NULL) = (voided_by IS NULL))
);
CREATE INDEX commission_sweeps_date_idx ON commission_sweeps (swept_on DESC, created_at DESC);
CREATE TRIGGER commission_sweeps_audit AFTER INSERT OR UPDATE OR DELETE ON commission_sweeps FOR EACH ROW EXECUTE FUNCTION audit_row_change();

-- +goose Down
DROP TABLE commission_sweeps;
