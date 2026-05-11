export interface User {
  id: string
  email: string
  username: string
  created_at: string
}

export interface Profile {
  id?: string
  user_id: string
  weight_kg?: number
  age?: number
  sex?: string
  created_at?: string
  updated_at?: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
}

export interface MuscleGroup {
  id: number
  name: string
}

export interface Category {
  id: number
  name: string
}

export interface Exercise {
  id: string
  name: string
  category_id: number
  category_name?: string
  primary_muscle: number
  primary_muscle_name?: string
  secondary_muscle?: number
  secondary_muscle_name?: string
  equipment?: string
  description?: string
  is_custom: boolean
  created_at: string
}

export interface Set {
  id: string
  session_exercise_id: string
  set_number: number
  weight_kg?: number
  reps?: number
  rpe?: number
  rest_seconds?: number
  created_at: string
}

export interface SessionExercise {
  id: string
  session_id: string
  exercise_id: string
  order_index: number
  created_at: string
  sets: Set[]
}

export interface Session {
  id: string
  user_id: string
  program_id?: string
  session_date: string
  duration_minutes?: number
  calories_burned?: number
  notes?: string
  created_at: string
  updated_at: string
  exercises: SessionExercise[]
}

export interface ProgramType {
  id: number
  name: string
}

export interface Program {
  id: string
  user_id: string
  name: string
  program_type_id: number
  program_type?: string
  notes?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface CalendarDay {
  date: string
  has_session: boolean
  session_id?: string
  duration_minutes?: number
  calories_burned?: number
}

export interface ProgressionPoint {
  date: string
  max_weight?: number
  max_reps?: number
  total_volume: number
  avg_rpe?: number
}

export interface WeeklyStats {
  week_start: string
  week_end: string
  total_sessions: number
  total_duration_minutes: number
  total_calories: number
}

export interface MonthlyStats {
  year: number
  month: number
  total_sessions: number
  total_duration_minutes: number
  total_calories: number
  training_days: number
}
