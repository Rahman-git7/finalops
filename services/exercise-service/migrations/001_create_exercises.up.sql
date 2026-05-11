CREATE TABLE IF NOT EXISTS exercise.muscle_groups (
    id    SMALLSERIAL PRIMARY KEY,
    name  VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise.categories (
    id    SMALLSERIAL PRIMARY KEY,
    name  VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS exercise.exercises (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name             VARCHAR(255) UNIQUE NOT NULL,
    category_id      SMALLINT NOT NULL REFERENCES exercise.categories(id),
    primary_muscle   SMALLINT NOT NULL REFERENCES exercise.muscle_groups(id),
    secondary_muscle SMALLINT REFERENCES exercise.muscle_groups(id),
    equipment        VARCHAR(100),
    description      TEXT,
    is_custom        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_exercises_category ON exercise.exercises(category_id);
CREATE INDEX IF NOT EXISTS idx_exercises_muscle   ON exercise.exercises(primary_muscle);
CREATE INDEX IF NOT EXISTS idx_exercises_name     ON exercise.exercises USING gin(to_tsvector('french', name));

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA exercise TO exercise_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA exercise TO exercise_user;
GRANT SELECT ON ALL TABLES IN SCHEMA exercise TO workout_user;
