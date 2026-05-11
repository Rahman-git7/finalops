import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { analyticsApi } from '../api/analytics'
import { workoutsApi } from '../api/workouts'
import { useAuthStore } from '../store/authStore'
import { format } from 'date-fns'
import { fr } from 'date-fns/locale'

function StatCard({ label, value, unit }: { label: string; value: string | number; unit?: string }) {
  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-5">
      <p className="text-xs font-medium text-zinc-500 uppercase tracking-wider mb-1">{label}</p>
      <p className="text-3xl font-bold text-white">
        {value}
        {unit && <span className="text-base font-normal text-zinc-500 ml-1">{unit}</span>}
      </p>
    </div>
  )
}

export default function DashboardPage() {
  const { user } = useAuthStore()
  const { data: weekStats } = useQuery({
    queryKey: ['weekly-stats'],
    queryFn: analyticsApi.weeklyStats,
  })
  const { data: monthStats } = useQuery({
    queryKey: ['monthly-stats'],
    queryFn: () => analyticsApi.monthlyStats(),
  })
  const { data: recentSessions } = useQuery({
    queryKey: ['recent-sessions'],
    queryFn: () => workoutsApi.listSessions({ limit: 5 }),
  })

  const today = format(new Date(), "EEEE d MMMM", { locale: fr })

  return (
    <div className="p-8">
      <div className="mb-8">
        <p className="text-zinc-500 text-sm capitalize">{today}</p>
        <h2 className="text-3xl font-bold text-white mt-1">
          Bienvenue, <span className="text-brand-500">{user?.username}</span>
        </h2>
      </div>

      <section className="mb-8">
        <h3 className="text-xs font-semibold text-zinc-500 uppercase tracking-wider mb-4">Cette semaine</h3>
        <div className="grid grid-cols-3 gap-4">
          <StatCard label="Séances" value={weekStats?.total_sessions ?? 0} />
          <StatCard label="Durée totale" value={weekStats?.total_duration_minutes ?? 0} unit="min" />
          <StatCard label="Calories" value={weekStats?.total_calories ?? 0} unit="kcal" />
        </div>
      </section>

      <section className="mb-8">
        <h3 className="text-xs font-semibold text-zinc-500 uppercase tracking-wider mb-4">Ce mois</h3>
        <div className="grid grid-cols-4 gap-4">
          <StatCard label="Jours d'entrainement" value={monthStats?.training_days ?? 0} />
          <StatCard label="Séances" value={monthStats?.total_sessions ?? 0} />
          <StatCard label="Durée totale" value={monthStats?.total_duration_minutes ?? 0} unit="min" />
          <StatCard label="Calories" value={monthStats?.total_calories ?? 0} unit="kcal" />
        </div>
      </section>

      <section>
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-xs font-semibold text-zinc-500 uppercase tracking-wider">Séances récentes</h3>
          <Link to="/workout" className="text-brand-500 text-xs font-medium hover:text-brand-400 transition-colors">
            + Nouvelle séance
          </Link>
        </div>
        <div className="space-y-2">
          {recentSessions?.length === 0 && (
            <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-8 text-center">
              <p className="text-zinc-500 text-sm">Aucune séance enregistrée.</p>
              <Link to="/workout" className="mt-3 inline-block text-brand-500 text-sm font-medium hover:text-brand-400">
                Logger ta première séance →
              </Link>
            </div>
          )}
          {recentSessions?.map(session => (
            <Link
              key={session.id}
              to={`/workout/${session.id}`}
              className="flex items-center justify-between bg-zinc-900 border border-zinc-800 hover:border-zinc-700 rounded-xl px-5 py-4 transition-colors group"
            >
              <div>
                <p className="text-sm font-semibold text-white group-hover:text-brand-400 transition-colors">
                  {format(new Date(session.session_date), 'EEEE d MMMM yyyy', { locale: fr })}
                </p>
                <p className="text-xs text-zinc-500 mt-0.5">
                  {session.exercises?.length ?? 0} exercice(s)
                </p>
              </div>
              <div className="flex items-center gap-6 text-right">
                {session.duration_minutes && (
                  <div>
                    <p className="text-xs text-zinc-500">Durée</p>
                    <p className="text-sm font-semibold text-white">{session.duration_minutes} min</p>
                  </div>
                )}
                {session.calories_burned && (
                  <div>
                    <p className="text-xs text-zinc-500">Calories</p>
                    <p className="text-sm font-semibold text-brand-400">{session.calories_burned} kcal</p>
                  </div>
                )}
              </div>
            </Link>
          ))}
        </div>
      </section>
    </div>
  )
}
