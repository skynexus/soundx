BEGIN;

CREATE TABLE sounds (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    genres     TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    credits    JSONB NOT NULL,
    bpm        INTEGER NOT NULL,
    duration   INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE playlists (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    sound_ids  BIGINT[] NOT NULL DEFAULT '{}'::BIGINT[],
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

COMMIT;