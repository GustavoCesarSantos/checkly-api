CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE urls (
    id bigserial PRIMARY KEY,
    external_id uuid NOT NULL DEFAULT uuid_generate_v4(),
    address text NOT NULL,
    interval integer NOT NULL,
    retry_limit integer NOT NULL,
    retry_count integer NOT NULL DEFAULT 0,
    stability_count integer NOT NULL DEFAULT 0,
    contact text NOT NULL,
    next_check timestamp(0) with time zone NOT NULL,
    status integer NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NULL
);

CREATE INDEX idx_urls_next_check_status ON urls (next_check, status);