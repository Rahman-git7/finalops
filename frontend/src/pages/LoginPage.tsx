import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth'
import { useAuthStore } from '../store/authStore'

export default function LoginPage() {
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [email, setEmail] = useState('')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { setTokens, setUser } = useAuthStore()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const tokens = mode === 'login'
        ? await authApi.login(email, password)
        : await authApi.register(email, username, password)
      setTokens(tokens.access_token, tokens.refresh_token)
      const user = await authApi.me()
      setUser(user)
      navigate('/dashboard')
    } catch (err: any) {
      setError(err.response?.data?.error ?? 'Une erreur est survenue')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-zinc-950 flex items-center justify-center p-4">
      <div className="w-full max-w-sm">
        <div className="text-center mb-10">
          <h1 className="text-5xl font-black tracking-tight mb-2">
            <span className="text-brand-500">FINAL</span>
            <span className="text-white">OPS</span>
          </h1>
          <p className="text-zinc-500 text-sm">Tracker de performance fitness</p>
        </div>

        <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-8">
          <div className="flex gap-2 mb-6 bg-zinc-800 rounded-lg p-1">
            <button
              onClick={() => setMode('login')}
              className={`flex-1 py-1.5 text-sm font-medium rounded-md transition-all ${
                mode === 'login' ? 'bg-brand-500 text-white' : 'text-zinc-400 hover:text-white'
              }`}
            >
              Connexion
            </button>
            <button
              onClick={() => setMode('register')}
              className={`flex-1 py-1.5 text-sm font-medium rounded-md transition-all ${
                mode === 'register' ? 'bg-brand-500 text-white' : 'text-zinc-400 hover:text-white'
              }`}
            >
              Inscription
            </button>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            {mode === 'register' && (
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1.5">Nom d'utilisateur</label>
                <input
                  type="text"
                  value={username}
                  onChange={e => setUsername(e.target.value)}
                  className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-white placeholder-zinc-600 focus:outline-none focus:border-brand-500 transition-colors"
                  placeholder="johndoe"
                  required
                />
              </div>
            )}
            <div>
              <label className="block text-xs font-medium text-zinc-400 mb-1.5">Email</label>
              <input
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-white placeholder-zinc-600 focus:outline-none focus:border-brand-500 transition-colors"
                placeholder="john@example.com"
                required
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-zinc-400 mb-1.5">Mot de passe</label>
              <input
                type="password"
                value={password}
                onChange={e => setPassword(e.target.value)}
                className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-white placeholder-zinc-600 focus:outline-none focus:border-brand-500 transition-colors"
                placeholder="••••••••"
                required
              />
            </div>

            {error && (
              <p className="text-red-400 text-xs bg-red-400/10 border border-red-400/20 rounded-lg px-3 py-2">
                {error}
              </p>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-brand-500 hover:bg-brand-600 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg transition-colors text-sm mt-2"
            >
              {loading ? '...' : mode === 'login' ? 'Se connecter' : "S'inscrire"}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}
