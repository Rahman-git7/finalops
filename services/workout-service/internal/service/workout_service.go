package service

import (
	"context"

	"github.com/finalops/workout-service/internal/model"
	"github.com/finalops/workout-service/internal/repository"
	"github.com/google/uuid"
)

type WorkoutService struct {
	repo *repository.WorkoutRepository
}

func NewWorkoutService(repo *repository.WorkoutRepository) *WorkoutService {
	return &WorkoutService{repo: repo}
}

func (s *WorkoutService) CreateSession(ctx context.Context, userID uuid.UUID, req *model.CreateSessionRequest) (*model.Session, error) {
	return s.repo.CreateSession(ctx, userID, req)
}

func (s *WorkoutService) ListSessions(ctx context.Context, userID uuid.UUID, from, to string, limit, offset int) ([]model.Session, error) {
	return s.repo.ListSessions(ctx, userID, from, to, limit, offset)
}

func (s *WorkoutService) GetSession(ctx context.Context, id, userID uuid.UUID) (*model.Session, error) {
	return s.repo.GetSession(ctx, id, userID)
}

func (s *WorkoutService) UpdateSession(ctx context.Context, id, userID uuid.UUID, req *model.UpdateSessionRequest) (*model.Session, error) {
	return s.repo.UpdateSession(ctx, id, userID, req)
}

func (s *WorkoutService) DeleteSession(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.DeleteSession(ctx, id, userID)
}

func (s *WorkoutService) AddExercise(ctx context.Context, sessionID uuid.UUID, req *model.AddExerciseRequest) (*model.SessionExercise, error) {
	return s.repo.AddExerciseToSession(ctx, sessionID, req)
}

func (s *WorkoutService) RemoveExercise(ctx context.Context, sessionID, exerciseID uuid.UUID) error {
	return s.repo.RemoveExerciseFromSession(ctx, sessionID, exerciseID)
}

func (s *WorkoutService) AddSet(ctx context.Context, sessionID, exerciseID uuid.UUID, req *model.CreateSetRequest) (*model.Set, error) {
	se, err := s.repo.GetSessionExerciseByExerciseID(ctx, sessionID, exerciseID)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateSet(ctx, se.ID, req)
}

func (s *WorkoutService) UpdateSet(ctx context.Context, setID uuid.UUID, req *model.UpdateSetRequest) (*model.Set, error) {
	return s.repo.UpdateSet(ctx, setID, req)
}

func (s *WorkoutService) DeleteSet(ctx context.Context, setID uuid.UUID) error {
	return s.repo.DeleteSet(ctx, setID)
}

func (s *WorkoutService) CreateProgram(ctx context.Context, userID uuid.UUID, req *model.CreateProgramRequest) (*model.Program, error) {
	return s.repo.CreateProgram(ctx, userID, req)
}

func (s *WorkoutService) ListPrograms(ctx context.Context, userID uuid.UUID) ([]model.Program, error) {
	return s.repo.ListPrograms(ctx, userID)
}

func (s *WorkoutService) GetProgram(ctx context.Context, id, userID uuid.UUID) (*model.Program, error) {
	return s.repo.GetProgram(ctx, id, userID)
}

func (s *WorkoutService) UpdateProgram(ctx context.Context, id, userID uuid.UUID, req *model.UpdateProgramRequest) (*model.Program, error) {
	return s.repo.UpdateProgram(ctx, id, userID, req)
}

func (s *WorkoutService) DeleteProgram(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.DeleteProgram(ctx, id, userID)
}

func (s *WorkoutService) ListProgramTypes(ctx context.Context) ([]model.ProgramType, error) {
	return s.repo.ListProgramTypes(ctx)
}

func (s *WorkoutService) GetSessionsInRange(ctx context.Context, userID uuid.UUID, from, to string) ([]model.SessionRangeItem, error) {
	return s.repo.GetSessionsInRange(ctx, userID, from, to)
}

func (s *WorkoutService) GetExerciseHistory(ctx context.Context, userID, exerciseID uuid.UUID) ([]model.ExerciseSetHistory, error) {
	return s.repo.GetExerciseHistory(ctx, userID, exerciseID)
}
