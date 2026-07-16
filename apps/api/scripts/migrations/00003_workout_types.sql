-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts
    ADD COLUMN workout_type TEXT NOT NULL DEFAULT 'easy',
    ADD COLUMN description TEXT NOT NULL DEFAULT '',
    ADD COLUMN scheduled_date DATE,
    ADD COLUMN status TEXT NOT NULL DEFAULT 'planned',
    ADD COLUMN completed_activity_id UUID REFERENCES activities(id),
    ADD COLUMN assigned_by UUID REFERENCES users(id),
    ADD COLUMN is_club_template BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE workouts ADD CONSTRAINT workouts_workout_type_check CHECK (
    workout_type IN ('easy', 'long', 'tempo', 'interval', 'fartlek', 'recovery', 'hills', 'race', 'cross', 'rest')
);

ALTER TABLE workouts ADD CONSTRAINT workouts_status_check CHECK (
    status IN ('planned', 'completed', 'skipped')
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts DROP CONSTRAINT IF EXISTS workouts_status_check;
ALTER TABLE workouts DROP CONSTRAINT IF EXISTS workouts_workout_type_check;
ALTER TABLE workouts
    DROP COLUMN IF EXISTS is_club_template,
    DROP COLUMN IF EXISTS assigned_by,
    DROP COLUMN IF EXISTS completed_activity_id,
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS scheduled_date,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS workout_type;
-- +goose StatementEnd
