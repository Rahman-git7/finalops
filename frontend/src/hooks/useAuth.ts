import { useAuthStore } from '../store/authStore'

export function useAuth() {
  const { accessToken, user } = useAuthStore()
  return {
    isAuthenticated: !!accessToken,
    user,
  }
}
