import { api } from './client'
import type { TokenPair, User, Profile } from '../types'

export const authApi = {
  register: (email: string, username: string, password: string) =>
    api.post<TokenPair>('/auth/register', { email, username, password }).then(r => r.data),

  login: (email: string, password: string) =>
    api.post<TokenPair>('/auth/login', { email, password }).then(r => r.data),

  logout: (refreshToken: string) =>
    api.post('/auth/logout', { refresh_token: refreshToken }),

  me: () => api.get<User>('/auth/me').then(r => r.data),

  getProfile: () => api.get<Profile>('/auth/profile').then(r => r.data),

  updateProfile: (data: { weight_kg?: number; age?: number; sex?: string }) =>
    api.put<Profile>('/auth/profile', data).then(r => r.data),
}
