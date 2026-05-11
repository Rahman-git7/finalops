import { api } from './client'
import type { Exercise, Category, MuscleGroup } from '../types'

export const exercisesApi = {
  list: (params?: { category?: number; muscle?: number; q?: string }) =>
    api.get<Exercise[]>('/exercises', { params }).then(r => r.data),

  get: (id: string) => api.get<Exercise>(`/exercises/${id}`).then(r => r.data),

  create: (data: Partial<Exercise>) =>
    api.post<Exercise>('/exercises', data).then(r => r.data),

  update: (id: string, data: Partial<Exercise>) =>
    api.put<Exercise>(`/exercises/${id}`, data).then(r => r.data),

  delete: (id: string) => api.delete(`/exercises/${id}`),

  categories: () => api.get<Category[]>('/exercises/categories').then(r => r.data),

  muscleGroups: () => api.get<MuscleGroup[]>('/exercises/muscle-groups').then(r => r.data),
}
