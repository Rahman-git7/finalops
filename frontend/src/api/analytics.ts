import { api } from './client'
import type { CalendarDay, ProgressionPoint, WeeklyStats, MonthlyStats } from '../types'

export const analyticsApi = {
  calendar: (year?: number, month?: number) =>
    api.get<CalendarDay[]>('/analytics/calendar', { params: { year, month } }).then(r => r.data),

  progression: (exerciseId: string) =>
    api.get<ProgressionPoint[]>(`/analytics/progression/${exerciseId}`).then(r => r.data),

  weeklyStats: () => api.get<WeeklyStats>('/analytics/stats/weekly').then(r => r.data),

  monthlyStats: (year?: number, month?: number) =>
    api.get<MonthlyStats>('/analytics/stats/monthly', { params: { year, month } }).then(r => r.data),
}
