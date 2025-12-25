CREATE TABLE sent_alerts (
    idempotency_key TEXT PRIMARY KEY,
    status integer NOT NULL DEFAULT 10,
    sent_at TIMESTAMP NULL
);