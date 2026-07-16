-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts DROP COLUMN IF EXISTS duration;
ALTER TABLE workouts DROP COLUMN IF EXISTS pace;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS duration TEXT NOT NULL DEFAULT '';
ALTER TABLE workouts ADD COLUMN IF NOT EXISTS pace TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd
