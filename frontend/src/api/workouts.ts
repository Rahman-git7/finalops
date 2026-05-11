import { api } from './client'
import type { Session, Program, ProgramType } from '../types'

export const workoutsApi = {
  createSession: (data: { session_date: string; program_id?: string; duration_minutes?: number; calories_burned?: number; notes?: string }) =>
    api.post<Session>('/workouts/sessions', data).then(r => r.data),

  listSessions: (params?: { from?: string; to?: string; limit?: number; offset?: number }) =>
    api.get<Session[]>('/workouts/sessions', { params }).then(r => r.data),

  getSession: (id: string) => api.get<Session>(`/workouts/sessions/${id}`).then(r => r.data),

  updateSession: (id: string, data: { duration_minutes?: number; calories_burned?: number; notes?: string }) =>
    api.put<Session>(`/workouts/sessions/${id}`, data).then(r => r.data),

  deleteSession: (id: string) => api.delete(`/workouts/sessions/${id}`),

  addExercise: (sessionId: string, exerciseId: string) =>
    api.post(`/workouts/sessions/${sessionId}/exercises`, { exercise_id: exerciseId }).then(r => r.data),

  removeExercise: (sessionId: string, exerciseId: string) =>
    api.delete(`/workouts/sessions/${sessionId}/exercises/${exerciseId}`),

  addSet: (sessionId: string, exerciseId: string, data: { set_number: number; weight_kg?: number; reps?: number; rpe?: number; rest_seconds?: number }) =>
    api.post(`/workouts/sessions/${sessionId}/exercises/${exerciseId}/sets`, data).then(r => r.data),

  updateSet: (sessionId: string, exerciseId: string, setId: string, data: { weight_kg?: number; reps?: number; rpe?: number; rest_seconds?: number }) =>
    api.put(`/workouts/sessions/${sessionId}/exercises/${exerciseId}/sets/${setId}`, data).then(r => r.data),

  deleteSet: (sessionId: string, exerciseId: string, setId: string) =>
    api.delete(`/workouts/sessions/${sessionId}/exercises/${exerciseId}/sets/${setId}`),

  programTypes: () => api.get<ProgramType[]>('/workouts/program-types').then(r => r.data),

  listPrograms: () => api.get<Program[]>('/workouts/programs').then(r => r.data),

  createProgram: (data: { name: string; program_type_id: number; notes?: string }) =>
    api.post<Program>('/workouts/programs', data).then(r => r.data),

  deleteProgram: (id: string) => api.delete(`/workouts/programs/${id}`),
}
