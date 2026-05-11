import { NavLink, useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../store/authStore'
import { authApi } from '../../api/auth'

const navItems = [
  { to: '/dashboard', label: 'Dashboard', icon: '⚡' },
  { to: '/workout', label: 'Log Séance', icon: '💪' },
  { to: '/calendar', label: 'Calendrier', icon: '📅' },
  { to: '/progress', label: 'Progression', icon: '📈' },
  { to: '/profile', label: 'Profil', icon: '👤' },
]

export default function Sidebar() {
  const { user, refreshToken, logout } = useAuthStore()
  const navigate = useNavigate()

  const handleLogout = async () => {
    if (refreshToken) {
      try { await authApi.logout(refreshToken) } catch {}
    }
    logout()
    navigate('/login')
  }

  return (
    <aside className="w-60 bg-zinc-900 border-r border-zinc-800 flex flex-col">
      <div className="p-6 border-b border-zinc-800">
        <h1 className="text-2xl font-black tracking-tight">
          <span className="text-brand-500">FINAL</span>
          <span className="text-white">OPS</span>
        </h1>
        <p className="text-xs text-zinc-500 mt-1 font-medium">FITNESS TRACKER</p>
      </div>

      <nav className="flex-1 p-4 space-y-1">
        {navItems.map(({ to, label, icon }) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 ${
                isActive
                  ? 'bg-brand-500/10 text-brand-400 border border-brand-500/20'
                  : 'text-zinc-400 hover:text-white hover:bg-zinc-800'
              }`
            }
          >
            <span className="text-base">{icon}</span>
            {label}
          </NavLink>
        ))}
      </nav>

      <div className="p-4 border-t border-zinc-800">
        <div className="flex items-center gap-3 mb-3">
          <div className="w-8 h-8 rounded-full bg-brand-500/20 flex items-center justify-center text-brand-400 text-sm font-bold">
            {user?.username?.charAt(0).toUpperCase() ?? 'U'}
          </div>
          <div className="min-w-0">
            <p className="text-sm font-medium text-white truncate">{user?.username ?? 'User'}</p>
            <p className="text-xs text-zinc-500 truncate">{user?.email ?? ''}</p>
          </div>
        </div>
        <button
          onClick={handleLogout}
          className="w-full text-xs text-zinc-500 hover:text-red-400 transition-colors py-1"
        >
          Déconnexion
        </button>
      </div>
    </aside>
  )
}
