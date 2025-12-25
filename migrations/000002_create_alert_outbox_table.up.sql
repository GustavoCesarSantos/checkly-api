CREATE TABLE alert_outbox (
    id bigserial PRIMARY KEY,
    url_id bigserial NOT NULL,
    idempotency_key TEXT NULL,
    payload JSONB NOT NULL,
    sent_at TIMESTAMP NULL,
    processing_at TIMESTAMP NULL,
    retry_count INT NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMP NOT NULL DEFAULT NOW(),
    locked_at TIMESTAMP NULL,
    locked_by TEXT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NULL
);

CREATE INDEX idx_alert_outbox_pending
ON alert_outbox (next_retry_at, created_at)
WHERE sent_at IS NULL;

CREATE INDEX idx_alert_outbox_locked_at
ON alert_outbox (locked_at)
WHERE sent_at IS NULL;

CREATE UNIQUE INDEX idx_alert_outbox_url_status
ON alert_outbox (url_id, status)
WHERE sent_at IS NULL;

CREATE UNIQUE INDEX idx_alert_outbox_idempotency_key
ON alert_outbox (idempotency_key)
WHERE sent_at IS NULL;
