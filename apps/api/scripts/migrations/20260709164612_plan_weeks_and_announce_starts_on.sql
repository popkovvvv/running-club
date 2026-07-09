-- +goose Up
-- +goose StatementBegin
ALTER TABLE announces ADD COLUMN starts_on DATE;

CREATE TABLE plan_weeks (
    id UUID PRIMARY KEY,
    club_id UUID NOT NULL REFERENCES clubs(id),
    week_index INT NOT NULL,
    range_label TEXT NOT NULL,
    plan_label TEXT NOT NULL,
    UNIQUE (club_id, week_index)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS plan_weeks;
ALTER TABLE announces DROP COLUMN IF EXISTS starts_on;
-- +goose StatementEnd
