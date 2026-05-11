package model

type CalendarDay struct {
	Date            string `json:"date"`
	HasSession      bool   `json:"has_session"`
	SessionID       string `json:"session_id,omitempty"`
	DurationMinutes *int16 `json:"duration_minutes,omitempty"`
	CaloriesBurned  *int16 `json:"calories_burned,omitempty"`
}

type ProgressionPoint struct {
	Date     string   `json:"date"`
	MaxWeight *float64 `json:"max_weight"`
	MaxReps   *int16   `json:"max_reps"`
	TotalVolume float64 `json:"total_volume"`
	AvgRPE   *float64 `json:"avg_rpe"`
}

type WeeklyStats struct {
	WeekStart       string  `json:"week_start"`
	WeekEnd         string  `json:"week_end"`
	TotalSessions   int     `json:"total_sessions"`
	TotalDuration   int     `json:"total_duration_minutes"`
	TotalCalories   int     `json:"total_calories"`
}

type MonthlyStats struct {
	Year            int     `json:"year"`
	Month           int     `json:"month"`
	TotalSessions   int     `json:"total_sessions"`
	TotalDuration   int     `json:"total_duration_minutes"`
	TotalCalories   int     `json:"total_calories"`
	TrainingDays    int     `json:"training_days"`
}

// From workout-service internal API

type SessionRangeItem struct {
	ID              string `json:"id"`
	SessionDate     string `json:"session_date"`
	DurationMinutes *int16 `json:"duration_minutes"`
	CaloriesBurned  *int16 `json:"calories_burned"`
}

type ExerciseSetHistory struct {
	SessionDate string   `json:"session_date"`
	SetNumber   int16    `json:"set_number"`
	WeightKg    *float64 `json:"weight_kg"`
	Reps        *int16   `json:"reps"`
	RPE         *float64 `json:"rpe"`
}
