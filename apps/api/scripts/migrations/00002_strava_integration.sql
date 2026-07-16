-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_integrations (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    status TEXT NOT NULL,
    external_athlete_id TEXT NOT NULL DEFAULT '',
    access_token TEXT NOT NULL DEFAULT '',
    refresh_token TEXT NOT NULL DEFAULT '',
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    scopes TEXT[] NOT NULL DEFAULT '{}',
    last_synced_at TIMESTAMPTZ,
    last_webhook_at TIMESTAMPTZ,
    last_error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, provider)
);

ALTER TABLE activities
    ADD COLUMN source TEXT NOT NULL DEFAULT '',
    ADD COLUMN external_id TEXT NOT NULL DEFAULT '',
    ADD COLUMN sport_type TEXT NOT NULL DEFAULT '',
    ADD COLUMN distance_meters DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN moving_seconds INT NOT NULL DEFAULT 0,
    ADD COLUMN elapsed_seconds INT NOT NULL DEFAULT 0,
    ADD COLUMN average_heartrate INT NOT NULL DEFAULT 0,
    ADD COLUMN max_heartrate INT NOT NULL DEFAULT 0,
    ADD COLUMN elevation_gain DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD COLUMN visibility TEXT NOT NULL DEFAULT '',
    ADD COLUMN polyline TEXT NOT NULL DEFAULT '',
    ADD COLUMN started_at TIMESTAMPTZ,
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE UNIQUE INDEX activities_user_source_external_id_idx
    ON activities (user_id, source, external_id)
    WHERE source <> '' AND external_id <> '';

CREATE TABLE activity_streams (
    id UUID PRIMARY KEY,
    activity_id UUID NOT NULL REFERENCES activities(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    data_json JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (activity_id, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activity_streams;
DROP INDEX IF EXISTS activities_user_source_external_id_idx;
ALTER TABLE activities
    DROP COLUMN IF EXISTS source,
    DROP COLUMN IF EXISTS external_id,
    DROP COLUMN IF EXISTS sport_type,
    DROP COLUMN IF EXISTS distance_meters,
    DROP COLUMN IF EXISTS moving_seconds,
    DROP COLUMN IF EXISTS elapsed_seconds,
    DROP COLUMN IF EXISTS average_heartrate,
    DROP COLUMN IF EXISTS max_heartrate,
    DROP COLUMN IF EXISTS elevation_gain,
    DROP COLUMN IF EXISTS visibility,
    DROP COLUMN IF EXISTS polyline,
    DROP COLUMN IF EXISTS started_at,
    DROP COLUMN IF EXISTS updated_at;
DROP TABLE IF EXISTS user_integrations;
-- +goose StatementEnd
