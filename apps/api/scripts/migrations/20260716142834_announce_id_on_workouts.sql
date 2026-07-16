-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts
    ADD COLUMN announce_id UUID REFERENCES announces(id) ON DELETE CASCADE;

CREATE UNIQUE INDEX workouts_user_announce_uidx
    ON workouts (user_id, announce_id)
    WHERE announce_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS workouts_user_announce_uidx;
ALTER TABLE workouts DROP COLUMN IF EXISTS announce_id;
-- +goose StatementEnd
