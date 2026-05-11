package service

import (
	"context"

	"github.com/finalops/exercise-service/internal/model"
	"github.com/finalops/exercise-service/internal/repository"
	"github.com/google/uuid"
)

type ExerciseService struct {
	repo *repository.ExerciseRepository
}

func NewExerciseService(repo *repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

func (s *ExerciseService) List(ctx context.Context, f model.ExerciseListFilter) ([]model.Exercise, error) {
	return s.repo.List(ctx, f)
}

func (s *ExerciseService) Get(ctx context.Context, id uuid.UUID) (*model.Exercise, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ExerciseService) Create(ctx context.Context, req *model.CreateExerciseRequest) (*model.Exercise, error) {
	return s.repo.Create(ctx, req)
}

func (s *ExerciseService) Update(ctx context.Context, id uuid.UUID, req *model.UpdateExerciseRequest) (*model.Exercise, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *ExerciseService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ExerciseService) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *ExerciseService) ListMuscleGroups(ctx context.Context) ([]model.MuscleGroup, error) {
	return s.repo.ListMuscleGroups(ctx)
}
