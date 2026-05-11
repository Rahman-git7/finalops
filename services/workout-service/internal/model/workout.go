package model

import (
	"time"

	"github.com/google/uuid"
)

type ProgramType struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type Program struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"user_id"`
	Name          string     `json:"name"`
	ProgramTypeID int16      `json:"program_type_id"`
	ProgramType   string     `json:"program_type,omitempty"`
	Notes         *string    `json:"notes"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type Set struct {
	ID                uuid.UUID  `json:"id"`
	SessionExerciseID uuid.UUID  `json:"session_exercise_id"`
	SetNumber         int16      `json:"set_number"`
	WeightKg          *float64   `json:"weight_kg"`
	Reps              *int16     `json:"reps"`
	RPE               *float64   `json:"rpe"`
	RestSeconds       *int16     `json:"rest_seconds"`
	CreatedAt         time.Time  `json:"created_at"`
}

type SessionExercise struct {
	ID         uuid.UUID `json:"id"`
	SessionID  uuid.UUID `json:"session_id"`
	ExerciseID uuid.UUID `json:"exercise_id"`
	OrderIndex int16     `json:"order_index"`
	CreatedAt  time.Time `json:"created_at"`
	Sets       []Set     `json:"sets,omitempty"`
}

type Session struct {
	ID              uuid.UUID         `json:"id"`
	UserID          uuid.UUID         `json:"user_id"`
	ProgramID       *uuid.UUID        `json:"program_id"`
	SessionDate     string            `json:"session_date"`
	DurationMinutes *int16            `json:"duration_minutes"`
	CaloriesBurned  *int16            `json:"calories_burned"`
	Notes           *string           `json:"notes"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Exercises       []SessionExercise `json:"exercises,omitempty"`
}

// Requests

type CreateSessionRequest struct {
	ProgramID       *uuid.UUID `json:"program_id"`
	SessionDate     string     `json:"session_date"`
	DurationMinutes *int16     `json:"duration_minutes"`
	CaloriesBurned  *int16     `json:"calories_burned"`
	Notes           *string    `json:"notes"`
}

type UpdateSessionRequest struct {
	DurationMinutes *int16  `json:"duration_minutes"`
	CaloriesBurned  *int16  `json:"calories_burned"`
	Notes           *string `json:"notes"`
}

type AddExerciseRequest struct {
	ExerciseID uuid.UUID `json:"exercise_id"`
	OrderIndex *int16    `json:"order_index"`
}

type CreateSetRequest struct {
	SetNumber   int16    `json:"set_number"`
	WeightKg    *float64 `json:"weight_kg"`
	Reps        *int16   `json:"reps"`
	RPE         *float64 `json:"rpe"`
	RestSeconds *int16   `json:"rest_seconds"`
}

type UpdateSetRequest struct {
	WeightKg    *float64 `json:"weight_kg"`
	Reps        *int16   `json:"reps"`
	RPE         *float64 `json:"rpe"`
	RestSeconds *int16   `json:"rest_seconds"`
}

type CreateProgramRequest struct {
	Name          string  `json:"name"`
	ProgramTypeID int16   `json:"program_type_id"`
	Notes         *string `json:"notes"`
}

type UpdateProgramRequest struct {
	Name          string  `json:"name"`
	ProgramTypeID int16   `json:"program_type_id"`
	Notes         *string `json:"notes"`
	IsActive      *bool   `json:"is_active"`
}

// Internal API models (used by analytics-service)

type SessionRangeItem struct {
	ID              uuid.UUID `json:"id"`
	SessionDate     string    `json:"session_date"`
	DurationMinutes *int16    `json:"duration_minutes"`
	CaloriesBurned  *int16    `json:"calories_burned"`
}

type ExerciseSetHistory struct {
	SessionDate string   `json:"session_date"`
	SetNumber   int16    `json:"set_number"`
	WeightKg    *float64 `json:"weight_kg"`
	Reps        *int16   `json:"reps"`
	RPE         *float64 `json:"rpe"`
}
