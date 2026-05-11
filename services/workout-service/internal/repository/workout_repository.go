package repository

import (
	"context"
	"fmt"

	"github.com/finalops/workout-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutRepository struct {
	db *pgxpool.Pool
}

func NewWorkoutRepository(db *pgxpool.Pool) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

// Sessions

func (r *WorkoutRepository) CreateSession(ctx context.Context, userID uuid.UUID, req *model.CreateSessionRequest) (*model.Session, error) {
	s := &model.Session{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO workout.sessions (user_id, program_id, session_date, duration_minutes, calories_burned, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, program_id, session_date::text, duration_minutes, calories_burned, notes, created_at, updated_at`,
		userID, req.ProgramID, req.SessionDate, req.DurationMinutes, req.CaloriesBurned, req.Notes,
	).Scan(&s.ID, &s.UserID, &s.ProgramID, &s.SessionDate, &s.DurationMinutes, &s.CaloriesBurned, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *WorkoutRepository) ListSessions(ctx context.Context, userID uuid.UUID, from, to string, limit, offset int) ([]model.Session, error) {
	query := `
		SELECT id, user_id, program_id, session_date::text, duration_minutes, calories_burned, notes, created_at, updated_at
		FROM workout.sessions
		WHERE user_id = $1`
	args := []any{userID}
	i := 2
	if from != "" {
		query += fmt.Sprintf(" AND session_date >= $%d", i)
		args = append(args, from)
		i++
	}
	if to != "" {
		query += fmt.Sprintf(" AND session_date <= $%d", i)
		args = append(args, to)
		i++
	}
	query += fmt.Sprintf(" ORDER BY session_date DESC LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.Session
	for rows.Next() {
		var s model.Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.ProgramID, &s.SessionDate,
			&s.DurationMinutes, &s.CaloriesBurned, &s.Notes, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	if sessions == nil {
		sessions = []model.Session{}
	}
	return sessions, nil
}

func (r *WorkoutRepository) GetSession(ctx context.Context, id, userID uuid.UUID) (*model.Session, error) {
	s := &model.Session{}
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, program_id, session_date::text, duration_minutes, calories_burned, notes, created_at, updated_at
		FROM workout.sessions WHERE id = $1 AND user_id = $2`, id, userID,
	).Scan(&s.ID, &s.UserID, &s.ProgramID, &s.SessionDate,
		&s.DurationMinutes, &s.CaloriesBurned, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}

	exercises, err := r.getSessionExercisesWithSets(ctx, id)
	if err != nil {
		return nil, err
	}
	s.Exercises = exercises
	return s, nil
}

func (r *WorkoutRepository) UpdateSession(ctx context.Context, id, userID uuid.UUID, req *model.UpdateSessionRequest) (*model.Session, error) {
	s := &model.Session{}
	err := r.db.QueryRow(ctx, `
		UPDATE workout.sessions
		SET duration_minutes=$3, calories_burned=$4, notes=$5, updated_at=NOW()
		WHERE id=$1 AND user_id=$2
		RETURNING id, user_id, program_id, session_date::text, duration_minutes, calories_burned, notes, created_at, updated_at`,
		id, userID, req.DurationMinutes, req.CaloriesBurned, req.Notes,
	).Scan(&s.ID, &s.UserID, &s.ProgramID, &s.SessionDate,
		&s.DurationMinutes, &s.CaloriesBurned, &s.Notes, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *WorkoutRepository) DeleteSession(ctx context.Context, id, userID uuid.UUID) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM workout.sessions WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

// Session Exercises

func (r *WorkoutRepository) AddExerciseToSession(ctx context.Context, sessionID uuid.UUID, req *model.AddExerciseRequest) (*model.SessionExercise, error) {
	orderIndex := int16(0)
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	} else {
		_ = r.db.QueryRow(ctx,
			`SELECT COALESCE(MAX(order_index)+1, 0) FROM workout.session_exercises WHERE session_id=$1`, sessionID,
		).Scan(&orderIndex)
	}

	se := &model.SessionExercise{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO workout.session_exercises (session_id, exercise_id, order_index)
		VALUES ($1, $2, $3)
		RETURNING id, session_id, exercise_id, order_index, created_at`,
		sessionID, req.ExerciseID, orderIndex,
	).Scan(&se.ID, &se.SessionID, &se.ExerciseID, &se.OrderIndex, &se.CreatedAt)
	return se, err
}

func (r *WorkoutRepository) RemoveExerciseFromSession(ctx context.Context, sessionID, exerciseID uuid.UUID) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM workout.session_exercises WHERE session_id=$1 AND exercise_id=$2`, sessionID, exerciseID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("exercise not found in session")
	}
	return nil
}

func (r *WorkoutRepository) GetSessionExerciseByExerciseID(ctx context.Context, sessionID, exerciseID uuid.UUID) (*model.SessionExercise, error) {
	se := &model.SessionExercise{}
	err := r.db.QueryRow(ctx,
		`SELECT id, session_id, exercise_id, order_index, created_at
         FROM workout.session_exercises WHERE session_id=$1 AND exercise_id=$2`,
		sessionID, exerciseID,
	).Scan(&se.ID, &se.SessionID, &se.ExerciseID, &se.OrderIndex, &se.CreatedAt)
	return se, err
}

// Sets

func (r *WorkoutRepository) CreateSet(ctx context.Context, sessionExerciseID uuid.UUID, req *model.CreateSetRequest) (*model.Set, error) {
	s := &model.Set{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO workout.sets (session_exercise_id, set_number, weight_kg, reps, rpe, rest_seconds)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, session_exercise_id, set_number, weight_kg, reps, rpe, rest_seconds, created_at`,
		sessionExerciseID, req.SetNumber, req.WeightKg, req.Reps, req.RPE, req.RestSeconds,
	).Scan(&s.ID, &s.SessionExerciseID, &s.SetNumber, &s.WeightKg, &s.Reps, &s.RPE, &s.RestSeconds, &s.CreatedAt)
	return s, err
}

func (r *WorkoutRepository) UpdateSet(ctx context.Context, setID uuid.UUID, req *model.UpdateSetRequest) (*model.Set, error) {
	s := &model.Set{}
	err := r.db.QueryRow(ctx, `
		UPDATE workout.sets SET weight_kg=$2, reps=$3, rpe=$4, rest_seconds=$5
		WHERE id=$1
		RETURNING id, session_exercise_id, set_number, weight_kg, reps, rpe, rest_seconds, created_at`,
		setID, req.WeightKg, req.Reps, req.RPE, req.RestSeconds,
	).Scan(&s.ID, &s.SessionExerciseID, &s.SetNumber, &s.WeightKg, &s.Reps, &s.RPE, &s.RestSeconds, &s.CreatedAt)
	return s, err
}

func (r *WorkoutRepository) DeleteSet(ctx context.Context, setID uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM workout.sets WHERE id=$1`, setID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("set not found")
	}
	return nil
}

// Programs

func (r *WorkoutRepository) CreateProgram(ctx context.Context, userID uuid.UUID, req *model.CreateProgramRequest) (*model.Program, error) {
	p := &model.Program{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO workout.programs (user_id, name, program_type_id, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, name, program_type_id, notes, is_active, created_at, updated_at`,
		userID, req.Name, req.ProgramTypeID, req.Notes,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.ProgramTypeID, &p.Notes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *WorkoutRepository) ListPrograms(ctx context.Context, userID uuid.UUID) ([]model.Program, error) {
	rows, err := r.db.Query(ctx, `
		SELECT p.id, p.user_id, p.name, p.program_type_id, pt.name, p.notes, p.is_active, p.created_at, p.updated_at
		FROM workout.programs p
		JOIN workout.program_types pt ON pt.id = p.program_type_id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var programs []model.Program
	for rows.Next() {
		var p model.Program
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.ProgramTypeID, &p.ProgramType,
			&p.Notes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		programs = append(programs, p)
	}
	if programs == nil {
		programs = []model.Program{}
	}
	return programs, nil
}

func (r *WorkoutRepository) GetProgram(ctx context.Context, id, userID uuid.UUID) (*model.Program, error) {
	p := &model.Program{}
	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.user_id, p.name, p.program_type_id, pt.name, p.notes, p.is_active, p.created_at, p.updated_at
		FROM workout.programs p
		JOIN workout.program_types pt ON pt.id = p.program_type_id
		WHERE p.id=$1 AND p.user_id=$2`, id, userID,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.ProgramTypeID, &p.ProgramType,
		&p.Notes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *WorkoutRepository) UpdateProgram(ctx context.Context, id, userID uuid.UUID, req *model.UpdateProgramRequest) (*model.Program, error) {
	p := &model.Program{}
	err := r.db.QueryRow(ctx, `
		UPDATE workout.programs SET name=$3, program_type_id=$4, notes=$5, is_active=COALESCE($6, is_active), updated_at=NOW()
		WHERE id=$1 AND user_id=$2
		RETURNING id, user_id, name, program_type_id, notes, is_active, created_at, updated_at`,
		id, userID, req.Name, req.ProgramTypeID, req.Notes, req.IsActive,
	).Scan(&p.ID, &p.UserID, &p.Name, &p.ProgramTypeID, &p.Notes, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (r *WorkoutRepository) DeleteProgram(ctx context.Context, id, userID uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM workout.programs WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("program not found")
	}
	return nil
}

func (r *WorkoutRepository) ListProgramTypes(ctx context.Context) ([]model.ProgramType, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM workout.program_types ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var types []model.ProgramType
	for rows.Next() {
		var t model.ProgramType
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	if types == nil {
		types = []model.ProgramType{}
	}
	return types, nil
}

// Internal API

func (r *WorkoutRepository) GetSessionsInRange(ctx context.Context, userID uuid.UUID, from, to string) ([]model.SessionRangeItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, session_date::text, duration_minutes, calories_burned
		FROM workout.sessions
		WHERE user_id=$1 AND session_date BETWEEN $2 AND $3
		ORDER BY session_date ASC`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.SessionRangeItem
	for rows.Next() {
		var item model.SessionRangeItem
		if err := rows.Scan(&item.ID, &item.SessionDate, &item.DurationMinutes, &item.CaloriesBurned); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if items == nil {
		items = []model.SessionRangeItem{}
	}
	return items, nil
}

func (r *WorkoutRepository) GetExerciseHistory(ctx context.Context, userID, exerciseID uuid.UUID) ([]model.ExerciseSetHistory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT s2.session_date::text, s.set_number, s.weight_kg, s.reps, s.rpe
		FROM workout.sets s
		JOIN workout.session_exercises se ON se.id = s.session_exercise_id
		JOIN workout.sessions s2 ON s2.id = se.session_id
		WHERE s2.user_id=$1 AND se.exercise_id=$2
		ORDER BY s2.session_date ASC, s.set_number ASC`, userID, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []model.ExerciseSetHistory
	for rows.Next() {
		var h model.ExerciseSetHistory
		if err := rows.Scan(&h.SessionDate, &h.SetNumber, &h.WeightKg, &h.Reps, &h.RPE); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	if history == nil {
		history = []model.ExerciseSetHistory{}
	}
	return history, nil
}

// helpers

func (r *WorkoutRepository) getSessionExercisesWithSets(ctx context.Context, sessionID uuid.UUID) ([]model.SessionExercise, error) {
	seRows, err := r.db.Query(ctx, `
		SELECT id, session_id, exercise_id, order_index, created_at
		FROM workout.session_exercises WHERE session_id=$1 ORDER BY order_index ASC`, sessionID)
	if err != nil {
		return nil, err
	}
	defer seRows.Close()

	var exercises []model.SessionExercise
	for seRows.Next() {
		var se model.SessionExercise
		if err := seRows.Scan(&se.ID, &se.SessionID, &se.ExerciseID, &se.OrderIndex, &se.CreatedAt); err != nil {
			return nil, err
		}
		exercises = append(exercises, se)
	}

	for i, se := range exercises {
		setRows, err := r.db.Query(ctx, `
			SELECT id, session_exercise_id, set_number, weight_kg, reps, rpe, rest_seconds, created_at
			FROM workout.sets WHERE session_exercise_id=$1 ORDER BY set_number ASC`, se.ID)
		if err != nil {
			return nil, err
		}
		var sets []model.Set
		for setRows.Next() {
			var s model.Set
			if err := setRows.Scan(&s.ID, &s.SessionExerciseID, &s.SetNumber,
				&s.WeightKg, &s.Reps, &s.RPE, &s.RestSeconds, &s.CreatedAt); err != nil {
				setRows.Close()
				return nil, err
			}
			sets = append(sets, s)
		}
		setRows.Close()
		if sets == nil {
			sets = []model.Set{}
		}
		exercises[i].Sets = sets
	}
	if exercises == nil {
		exercises = []model.SessionExercise{}
	}
	return exercises, nil
}
