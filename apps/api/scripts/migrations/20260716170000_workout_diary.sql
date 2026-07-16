-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts
    ADD COLUMN rpe SMALLINT,
    ADD COLUMN athlete_report TEXT NOT NULL DEFAULT '',
    ADD COLUMN coach_comment TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts
    DROP COLUMN IF EXISTS rpe,
    DROP COLUMN IF EXISTS athlete_report,
    DROP COLUMN IF EXISTS coach_comment;
-- +goose StatementEnd
