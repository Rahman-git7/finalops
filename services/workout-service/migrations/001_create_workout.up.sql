CREATE TABLE IF NOT EXISTS workout.program_types (
    id    SMALLSERIAL PRIMARY KEY,
    name  VARCHAR(50) UNIQUE NOT NULL
);

INSERT INTO workout.program_types (name) VALUES
  ('Split'), ('PPL'), ('Upper/Lower'), ('Full Body')
ON CONFLICT (name) DO NOTHING;

CREATE TABLE IF NOT EXISTS workout.programs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    name            VARCHAR(255) NOT NULL,
    program_type_id SMALLINT NOT NULL REFERENCES workout.program_types(id),
    notes           TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS workout.sessions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL,
    program_id       UUID REFERENCES workout.programs(id) ON DELETE SET NULL,
    session_date     DATE NOT NULL,
    duration_minutes SMALLINT,
    calories_burned  SMALLINT,
    notes            TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_date ON workout.sessions(user_id, session_date DESC);

CREATE TABLE IF NOT EXISTS workout.session_exercises (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id  UUID NOT NULL REFERENCES workout.sessions(id) ON DELETE CASCADE,
    exercise_id UUID NOT NULL,
    order_index SMALLINT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_session_exercises_session ON workout.session_exercises(session_id);

CREATE TABLE IF NOT EXISTS workout.sets (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_exercise_id UUID NOT NULL REFERENCES workout.session_exercises(id) ON DELETE CASCADE,
    set_number          SMALLINT NOT NULL,
    weight_kg           NUMERIC(6,2),
    reps                SMALLINT,
    rpe                 NUMERIC(3,1),
    rest_seconds        SMALLINT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sets_session_exercise ON workout.sets(session_exercise_id);
