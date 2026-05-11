package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/finalops/exercise-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExerciseRepository struct {
	db *pgxpool.Pool
}

func NewExerciseRepository(db *pgxpool.Pool) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) List(ctx context.Context, f model.ExerciseListFilter) ([]model.Exercise, error) {
	query := `
		SELECT e.id, e.name, e.category_id, c.name,
		       e.primary_muscle, mg1.name,
		       e.secondary_muscle, mg2.name,
		       e.equipment, e.description, e.is_custom, e.created_at
		FROM exercise.exercises e
		JOIN exercise.categories c ON c.id = e.category_id
		JOIN exercise.muscle_groups mg1 ON mg1.id = e.primary_muscle
		LEFT JOIN exercise.muscle_groups mg2 ON mg2.id = e.secondary_muscle
		WHERE 1=1`

	args := []any{}
	i := 1

	if f.CategoryID != nil {
		query += fmt.Sprintf(" AND e.category_id = $%d", i)
		args = append(args, *f.CategoryID)
		i++
	}
	if f.MuscleID != nil {
		query += fmt.Sprintf(" AND (e.primary_muscle = $%d OR e.secondary_muscle = $%d)", i, i)
		args = append(args, *f.MuscleID)
		i++
	}
	if f.Query != "" {
		query += fmt.Sprintf(" AND LOWER(e.name) LIKE $%d", i)
		args = append(args, "%"+strings.ToLower(f.Query)+"%")
		i++
	}
	_ = i
	query += " ORDER BY e.name ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []model.Exercise
	for rows.Next() {
		var ex model.Exercise
		if err := rows.Scan(
			&ex.ID, &ex.Name, &ex.CategoryID, &ex.CategoryName,
			&ex.PrimaryMuscle, &ex.PrimaryMuscleName,
			&ex.SecondaryMuscle, &ex.SecondaryMuscleName,
			&ex.Equipment, &ex.Description, &ex.IsCustom, &ex.CreatedAt,
		); err != nil {
			return nil, err
		}
		exercises = append(exercises, ex)
	}
	if exercises == nil {
		exercises = []model.Exercise{}
	}
	return exercises, nil
}

func (r *ExerciseRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Exercise, error) {
	ex := &model.Exercise{}
	err := r.db.QueryRow(ctx, `
		SELECT e.id, e.name, e.category_id, c.name,
		       e.primary_muscle, mg1.name,
		       e.secondary_muscle, mg2.name,
		       e.equipment, e.description, e.is_custom, e.created_at
		FROM exercise.exercises e
		JOIN exercise.categories c ON c.id = e.category_id
		JOIN exercise.muscle_groups mg1 ON mg1.id = e.primary_muscle
		LEFT JOIN exercise.muscle_groups mg2 ON mg2.id = e.secondary_muscle
		WHERE e.id = $1`, id,
	).Scan(
		&ex.ID, &ex.Name, &ex.CategoryID, &ex.CategoryName,
		&ex.PrimaryMuscle, &ex.PrimaryMuscleName,
		&ex.SecondaryMuscle, &ex.SecondaryMuscleName,
		&ex.Equipment, &ex.Description, &ex.IsCustom, &ex.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ex, nil
}

func (r *ExerciseRepository) Create(ctx context.Context, req *model.CreateExerciseRequest) (*model.Exercise, error) {
	id := uuid.New()
	_, err := r.db.Exec(ctx, `
		INSERT INTO exercise.exercises (id, name, category_id, primary_muscle, secondary_muscle, equipment, description, is_custom)
		VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE)`,
		id, req.Name, req.CategoryID, req.PrimaryMuscle, req.SecondaryMuscle, req.Equipment, req.Description,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *ExerciseRepository) Update(ctx context.Context, id uuid.UUID, req *model.UpdateExerciseRequest) (*model.Exercise, error) {
	_, err := r.db.Exec(ctx, `
		UPDATE exercise.exercises
		SET name=$2, category_id=$3, primary_muscle=$4, secondary_muscle=$5, equipment=$6, description=$7
		WHERE id=$1 AND is_custom=TRUE`,
		id, req.Name, req.CategoryID, req.PrimaryMuscle, req.SecondaryMuscle, req.Equipment, req.Description,
	)
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, id)
}

func (r *ExerciseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM exercise.exercises WHERE id=$1 AND is_custom=TRUE`, id,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("exercise not found or not deletable")
	}
	return nil
}

func (r *ExerciseRepository) ListCategories(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM exercise.categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	if cats == nil {
		cats = []model.Category{}
	}
	return cats, nil
}

func (r *ExerciseRepository) ListMuscleGroups(ctx context.Context) ([]model.MuscleGroup, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM exercise.muscle_groups ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var mgs []model.MuscleGroup
	for rows.Next() {
		var mg model.MuscleGroup
		if err := rows.Scan(&mg.ID, &mg.Name); err != nil {
			return nil, err
		}
		mgs = append(mgs, mg)
	}
	if mgs == nil {
		mgs = []model.MuscleGroup{}
	}
	return mgs, nil
}
