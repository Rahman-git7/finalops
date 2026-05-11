package model

import (
	"time"

	"github.com/google/uuid"
)

type MuscleGroup struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   int16  `json:"id"`
	Name string `json:"name"`
}

type Exercise struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	CategoryID      int16      `json:"category_id"`
	CategoryName    string     `json:"category_name,omitempty"`
	PrimaryMuscle   int16      `json:"primary_muscle"`
	PrimaryMuscleName string   `json:"primary_muscle_name,omitempty"`
	SecondaryMuscle *int16     `json:"secondary_muscle"`
	SecondaryMuscleName *string `json:"secondary_muscle_name,omitempty"`
	Equipment       *string    `json:"equipment"`
	Description     *string    `json:"description"`
	IsCustom        bool       `json:"is_custom"`
	CreatedAt       time.Time  `json:"created_at"`
}

type CreateExerciseRequest struct {
	Name            string  `json:"name"`
	CategoryID      int16   `json:"category_id"`
	PrimaryMuscle   int16   `json:"primary_muscle"`
	SecondaryMuscle *int16  `json:"secondary_muscle"`
	Equipment       *string `json:"equipment"`
	Description     *string `json:"description"`
}

type UpdateExerciseRequest struct {
	Name            string  `json:"name"`
	CategoryID      int16   `json:"category_id"`
	PrimaryMuscle   int16   `json:"primary_muscle"`
	SecondaryMuscle *int16  `json:"secondary_muscle"`
	Equipment       *string `json:"equipment"`
	Description     *string `json:"description"`
}

type ExerciseListFilter struct {
	CategoryID *int16
	MuscleID   *int16
	Query      string
}
