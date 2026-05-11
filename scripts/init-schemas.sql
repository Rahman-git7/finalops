-- Initialize schemas and roles for FinalOps
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS exercise;
CREATE SCHEMA IF NOT EXISTS workout;

-- Per-schema users (least privilege)
DO $$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'auth_user') THEN
    CREATE USER auth_user WITH PASSWORD 'auth_pass';
  END IF;
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'exercise_user') THEN
    CREATE USER exercise_user WITH PASSWORD 'exercise_pass';
  END IF;
  IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'workout_user') THEN
    CREATE USER workout_user WITH PASSWORD 'workout_pass';
  END IF;
END
$$;

GRANT ALL PRIVILEGES ON SCHEMA auth TO auth_user;
GRANT ALL PRIVILEGES ON SCHEMA exercise TO exercise_user;
GRANT ALL PRIVILEGES ON SCHEMA workout TO workout_user;

-- workout_user needs read on exercise schema for exercise_id validation queries
GRANT USAGE ON SCHEMA exercise TO workout_user;
GRANT SELECT ON ALL TABLES IN SCHEMA exercise TO workout_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA exercise GRANT SELECT ON TABLES TO workout_user;
