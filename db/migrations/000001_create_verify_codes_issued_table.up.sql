CREATE TABLE IF NOT EXISTS verify_codes_issued(
    id           bigserial PRIMARY KEY,
    user_id      uuid        NOT NULL,
    code_hash    text        NOT NULL,
    purpose      smallint     NOT NULL,
    expires_at   timestamptz NOT NULL,
    attempts     int         NOT NULL DEFAULT 0,
    max_attempts int         NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now(),
    consumed_at  timestamptz NULL
);
