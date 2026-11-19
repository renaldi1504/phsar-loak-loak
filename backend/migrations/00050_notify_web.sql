-- +goose Up
-- Notifikasi web yang berbunyi: seller (Seller Center), pembeli (website), admin (per izin peran).
--   notifications.audience   : untuk siapa baris ini (lonceng & aliran realtime difilter per audience)
--   notifications.push_*     : outbox Web Push (VAPID), dikirim worker seperti telegram_*
--   notification_prefs       : pengaturan per user per audience (grup on/off: suara, push, Telegram; jam tenang)
--   push_subscriptions       : langganan Web Push per browser
--   pg_notify('phsar_notify'): dikirim saat COMMIT → aliran realtime (SSE) di semua proses server
--   notify_admins(izin, …)   : satu baris per admin aktif yang perannya punya izin itu (atau '*')

ALTER TABLE notifications ADD COLUMN audience varchar(8) NOT NULL DEFAULT 'buyer'
  CHECK (audience IN ('seller', 'buyer', 'courier', 'admin'));
UPDATE notifications SET audience = CASE WHEN store_id IS NOT NULL THEN 'seller' WHEN kind LIKE 'courier.%' THEN 'courier' ELSE 'buyer' END;
ALTER TABLE notifications ADD COLUMN push_status varchar(12) NOT NULL DEFAULT 'pending'
  CHECK (push_status IN ('pending', 'sending', 'sent', 'failed', 'skipped'));
UPDATE notifications SET push_status = 'skipped'; -- notifikasi lama tidak dikirim ulang
ALTER TABLE notifications ADD COLUMN push_attempts integer NOT NULL DEFAULT 0;
ALTER TABLE notifications ADD COLUMN push_next_at timestamptz;
ALTER TABLE notifications ADD COLUMN push_error text;
CREATE INDEX notifications_push_outbox ON notifications (created_at) WHERE push_status = 'pending';
CREATE INDEX notifications_audience_idx ON notifications (user_id, audience, id DESC);
CREATE INDEX notifications_unread_aud_idx ON notifications (user_id, audience) WHERE read_at IS NULL;

CREATE TABLE notification_prefs (
  user_id    uuid NOT NULL REFERENCES users(id),
  audience   varchar(8) NOT NULL CHECK (audience IN ('seller', 'buyer', 'courier', 'admin')),
  prefs      jsonb NOT NULL DEFAULT '{}',   -- {"groups":{"orders":{"sound":true,"push":true,"telegram":true}},"quiet_from":"23:00","quiet_to":"07:00","volume":80}
  updated_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, audience)
);

CREATE TABLE push_subscriptions (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         uuid NOT NULL REFERENCES users(id),
  audience        varchar(8) NOT NULL CHECK (audience IN ('seller', 'buyer', 'courier', 'admin')),
  endpoint        text NOT NULL UNIQUE,
  p256dh          text NOT NULL,             -- kunci publik browser (base64url)
  auth            text NOT NULL,             -- rahasia auth browser (base64url)
  user_agent      varchar(200),
  created_at      timestamptz NOT NULL DEFAULT now(),
  last_success_at timestamptz,
  failures        integer NOT NULL DEFAULT 0
);
CREATE INDEX push_subscriptions_user ON push_subscriptions (user_id, audience);

-- Realtime: notifikasi baru & "sudah dibaca" (sekali per user per pernyataan, mis. "tandai semua dibaca").
-- +goose StatementBegin
CREATE FUNCTION notifications_realtime_new() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  PERFORM pg_notify('phsar_notify', json_build_object('u', NEW.user_id, 'id', NEW.id, 'a', NEW.audience, 'ev', 'new')::text);
  RETURN NULL;
END $$;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE FUNCTION notifications_realtime_read() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE r record;
BEGIN
  FOR r IN SELECT DISTINCT n.user_id, n.audience FROM changed n JOIN unchanged o ON o.id = n.id
           WHERE n.read_at IS NOT NULL AND o.read_at IS NULL LOOP
    PERFORM pg_notify('phsar_notify', json_build_object('u', r.user_id, 'id', 0, 'a', r.audience, 'ev', 'read')::text);
  END LOOP;
  RETURN NULL;
END $$;
-- +goose StatementEnd
CREATE TRIGGER notifications_realtime_new AFTER INSERT ON notifications FOR EACH ROW EXECUTE FUNCTION notifications_realtime_new();
CREATE TRIGGER notifications_realtime_read AFTER UPDATE ON notifications REFERENCING OLD TABLE AS unchanged NEW TABLE AS changed
  FOR EACH STATEMENT EXECUTE FUNCTION notifications_realtime_read();

-- Notifikasi admin per izin peran.
-- +goose StatementBegin
CREATE FUNCTION notify_admins(p_perm text, p_kind text, p_params jsonb, p_link text) RETURNS void LANGUAGE sql AS $$
  INSERT INTO notifications (user_id, kind, params, link, audience)
  SELECT au.user_id, p_kind, coalesce(p_params, '{}'), p_link, 'admin'
  FROM admin_users au JOIN users u ON u.id = au.user_id AND u.deleted_at IS NULL
  WHERE au.active AND au.deleted_at IS NULL
    AND EXISTS (SELECT 1 FROM role_permissions rp WHERE rp.role_code = au.role_code AND rp.permission IN (p_perm, '*'));
$$;
-- +goose StatementEnd

-- Pemicu antrean admin. DEFERRABLE INITIALLY DEFERRED = dijalankan saat COMMIT, jadi nama produk (ditulis
-- sesudah baris produk dalam transaksi yang sama) sudah ada, dan perubahan yang dibatalkan tidak mengirim apa pun.
-- +goose StatementBegin
CREATE FUNCTION admin_queue_notify() RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE
  is_new boolean := TG_OP = 'INSERT';
  v_store text;
  v_name text;
  v_order text;
BEGIN
  IF TG_TABLE_NAME = 'stores' THEN
    IF NEW.status = 'pending_review' AND (is_new OR OLD.status IS DISTINCT FROM 'pending_review') AND NEW.deleted_at IS NULL THEN
      PERFORM notify_admins('seller.verify', 'admin.store_review', jsonb_build_object('store', NEW.name), '#/store/' || NEW.id);
    END IF;

  ELSIF TG_TABLE_NAME = 'products' THEN
    IF NEW.moderation_status = 'pending_review' AND NEW.status <> 'draft' AND NEW.deleted_at IS NULL
       AND (is_new OR OLD.moderation_status IS DISTINCT FROM 'pending_review' OR OLD.status = 'draft') THEN
      SELECT s.name, coalesce(t.name, '') INTO v_store, v_name FROM stores s
        LEFT JOIN product_translations t ON t.product_id = NEW.id AND t.lang = NEW.source_lang WHERE s.id = NEW.store_id;
      PERFORM notify_admins('product.moderate', 'admin.product_review', jsonb_build_object('store', v_store, 'product', v_name), '#/product/' || NEW.id);
    END IF;

  ELSIF TG_TABLE_NAME = 'orders' THEN
    IF NEW.refund_status = 'pending' AND (is_new OR OLD.refund_status IS DISTINCT FROM 'pending') THEN
      PERFORM notify_admins('refund.decide', 'admin.refund_pending', jsonb_build_object('order', NEW.order_no), '#/refunds');
    END IF;

  ELSIF TG_TABLE_NAME = 'reports' THEN
    -- satu notifikasi per target (laporan berikutnya untuk target yang sama hanya menambah hitungan di antrean)
    IF NEW.status = 'open' AND NOT EXISTS (SELECT 1 FROM reports r WHERE r.status = 'open' AND r.id <> NEW.id
         AND r.target_type = NEW.target_type AND coalesce(r.product_id, r.store_id) = coalesce(NEW.product_id, NEW.store_id)) THEN
      SELECT s.name INTO v_store FROM stores s WHERE s.id = NEW.store_id;
      PERFORM notify_admins('report.handle', 'admin.report_new', jsonb_build_object('store', v_store, 'note', NEW.reason), '#/reports');
    END IF;

  ELSIF TG_TABLE_NAME = 'tickets' THEN
    IF is_new THEN
      PERFORM notify_admins('support.handle', 'admin.ticket_new', jsonb_build_object('note', NEW.subject, 'order', NEW.ticket_no), '#/support');
    ELSIF NEW.last_user_message_at IS DISTINCT FROM OLD.last_user_message_at AND NEW.status = 'open'
          AND (NEW.last_admin_reply_at IS NULL OR NEW.last_user_message_at > NEW.last_admin_reply_at) THEN
      PERFORM notify_admins('support.handle', 'admin.ticket_reply', jsonb_build_object('note', NEW.subject, 'order', NEW.ticket_no), '#/support');
    END IF;

  ELSIF TG_TABLE_NAME = 'returns' THEN
    SELECT o.order_no INTO v_order FROM orders o WHERE o.id = NEW.order_id;
    IF NEW.status = 'disputed' AND (is_new OR OLD.status IS DISTINCT FROM 'disputed') THEN
      PERFORM notify_admins('return.decide', 'admin.return_dispute', jsonb_build_object('order', v_order), '#/returns');
    ELSIF NEW.status = 'approved' AND (is_new OR OLD.status IS DISTINCT FROM 'approved') THEN
      PERFORM notify_admins('refund.decide', 'admin.return_refund', jsonb_build_object('order', v_order), '#/returns?s=approved');
    END IF;

  ELSIF TG_TABLE_NAME = 'payouts' THEN
    IF NEW.status = 'pending' AND (is_new OR OLD.status IS DISTINCT FROM 'pending') THEN
      SELECT s.name INTO v_store FROM stores s WHERE s.id = NEW.store_id;
      PERFORM notify_admins('payout.approve', 'admin.payout_pending', jsonb_build_object('store', v_store, 'product', '$' || NEW.amount::text), '#/payouts');
    END IF;

  ELSIF TG_TABLE_NAME = 'deliveries' THEN
    SELECT o.order_no INTO v_order FROM orders o WHERE o.id = NEW.order_id;
    IF NEW.status = 'unassigned' AND (is_new OR OLD.status IS DISTINCT FROM 'unassigned') THEN
      PERFORM notify_admins('delivery.dispatch', 'admin.delivery_new', jsonb_build_object('order', v_order), '#/delivery');
    ELSIF NEW.status = 'failed' AND (is_new OR OLD.status IS DISTINCT FROM 'failed') THEN
      PERFORM notify_admins('delivery.dispatch', 'admin.delivery_failed', jsonb_build_object('order', v_order, 'note', coalesce(NEW.fail_reason, '')), '#/delivery');
    END IF;

  ELSIF TG_TABLE_NAME = 'store_pickup_points' THEN
    IF is_new AND NEW.verified_at IS NULL AND NEW.deleted_at IS NULL THEN
      SELECT s.name INTO v_store FROM stores s WHERE s.id = NEW.store_id AND s.status = 'active';
      IF FOUND THEN
        PERFORM notify_admins('seller.verify', 'admin.pickup_review', jsonb_build_object('store', v_store, 'point', NEW.label), '#/store/' || NEW.store_id);
      END IF;
    END IF;
  END IF;
  RETURN NULL;
END $$;
-- +goose StatementEnd

CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT OR UPDATE ON stores DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();
CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT OR UPDATE ON products DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();
CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT OR UPDATE ON orders DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();
CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT ON reports DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();
CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT OR UPDATE ON tickets DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();
CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT OR UPDATE ON returns DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();
CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT OR UPDATE ON payouts DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();
CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT OR UPDATE ON deliveries DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();
CREATE CONSTRAINT TRIGGER admin_queue_notify AFTER INSERT ON store_pickup_points DEFERRABLE INITIALLY DEFERRED FOR EACH ROW EXECUTE FUNCTION admin_queue_notify();

-- +goose Down
DROP TRIGGER IF EXISTS admin_queue_notify ON stores;
DROP TRIGGER IF EXISTS admin_queue_notify ON products;
DROP TRIGGER IF EXISTS admin_queue_notify ON orders;
DROP TRIGGER IF EXISTS admin_queue_notify ON reports;
DROP TRIGGER IF EXISTS admin_queue_notify ON tickets;
DROP TRIGGER IF EXISTS admin_queue_notify ON returns;
DROP TRIGGER IF EXISTS admin_queue_notify ON payouts;
DROP TRIGGER IF EXISTS admin_queue_notify ON deliveries;
DROP TRIGGER IF EXISTS admin_queue_notify ON store_pickup_points;
DROP FUNCTION IF EXISTS admin_queue_notify();
DROP FUNCTION IF EXISTS notify_admins(text, text, jsonb, text);
DROP TRIGGER IF EXISTS notifications_realtime_new ON notifications;
DROP TRIGGER IF EXISTS notifications_realtime_read ON notifications;
DROP FUNCTION IF EXISTS notifications_realtime_new();
DROP FUNCTION IF EXISTS notifications_realtime_read();
DROP TABLE push_subscriptions;
DROP TABLE notification_prefs;
DELETE FROM notifications WHERE audience = 'admin';
DROP INDEX notifications_unread_aud_idx;
DROP INDEX notifications_audience_idx;
DROP INDEX notifications_push_outbox;
ALTER TABLE notifications DROP COLUMN push_error, DROP COLUMN push_next_at, DROP COLUMN push_attempts, DROP COLUMN push_status, DROP COLUMN audience;
