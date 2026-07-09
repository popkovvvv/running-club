-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('athlete', 'coach')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE clubs (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    invite_code TEXT NOT NULL UNIQUE,
    accent_hex TEXT NOT NULL DEFAULT '#ff5c22',
    coach_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE memberships (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    club_id UUID NOT NULL REFERENCES clubs(id),
    status TEXT NOT NULL CHECK (status IN ('active', 'left', 'removed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, club_id)
);

CREATE TABLE announces (
    id UUID PRIMARY KEY,
    club_id UUID NOT NULL REFERENCES clubs(id),
    place TEXT NOT NULL,
    day_label TEXT NOT NULL,
    time TEXT NOT NULL,
    group_name TEXT NOT NULL,
    note TEXT NOT NULL DEFAULT '',
    going_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE announce_signups (
    id UUID PRIMARY KEY,
    announce_id UUID NOT NULL REFERENCES announces(id) ON DELETE CASCADE,
    athlete_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (announce_id, athlete_id)
);

CREATE TABLE workouts (
    id UUID PRIMARY KEY,
    club_id UUID REFERENCES clubs(id),
    user_id UUID NOT NULL REFERENCES users(id),
    kind TEXT NOT NULL CHECK (kind IN ('plan', 'own', 'builder')),
    day_label TEXT NOT NULL DEFAULT '',
    tag TEXT NOT NULL DEFAULT '',
    title TEXT NOT NULL,
    dist_km DOUBLE PRECISION NOT NULL DEFAULT 0,
    duration TEXT NOT NULL DEFAULT '',
    pace TEXT NOT NULL DEFAULT '',
    hr TEXT NOT NULL DEFAULT '',
    week_index INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE segments (
    id UUID PRIMARY KEY,
    workout_id UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    kind TEXT NOT NULL,
    title TEXT NOT NULL,
    dist_km DOUBLE PRECISION NOT NULL DEFAULT 0,
    pace TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE activities (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    when_label TEXT NOT NULL,
    dist_km DOUBLE PRECISION NOT NULL,
    duration TEXT NOT NULL,
    pace TEXT NOT NULL,
    hr INT NOT NULL DEFAULT 0,
    kudos INT NOT NULL DEFAULT 0,
    comments INT NOT NULL DEFAULT 0,
    route_svg TEXT NOT NULL DEFAULT '',
    start_x DOUBLE PRECISION NOT NULL DEFAULT 0,
    start_y DOUBLE PRECISION NOT NULL DEFAULT 0,
    end_x DOUBLE PRECISION NOT NULL DEFAULT 0,
    end_y DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE prs (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    distance TEXT NOT NULL,
    time TEXT NOT NULL,
    date_label TEXT NOT NULL,
    pending BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE races (
    id UUID PRIMARY KEY,
    club_id UUID REFERENCES clubs(id),
    user_id UUID REFERENCES users(id),
    name TEXT NOT NULL,
    date_label TEXT NOT NULL,
    dist TEXT NOT NULL,
    goal TEXT NOT NULL DEFAULT '',
    days_left INT NOT NULL DEFAULT 0,
    finished BOOLEAN NOT NULL DEFAULT FALSE,
    result TEXT NOT NULL DEFAULT ''
);

CREATE TABLE month_stats (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    month TEXT NOT NULL,
    km DOUBLE PRECISION NOT NULL,
    tr INT NOT NULL,
    pace TEXT NOT NULL,
    diff TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS month_stats;
DROP TABLE IF EXISTS races;
DROP TABLE IF EXISTS prs;
DROP TABLE IF EXISTS activities;
DROP TABLE IF EXISTS segments;
DROP TABLE IF EXISTS workouts;
DROP TABLE IF EXISTS announce_signups;
DROP TABLE IF EXISTS announces;
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS clubs;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
