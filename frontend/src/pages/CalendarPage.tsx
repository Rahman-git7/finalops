import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { analyticsApi } from '../api/analytics'
import { format, startOfMonth, getDaysInMonth, getDay } from 'date-fns'
import { fr } from 'date-fns/locale'

const DAYS = ['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim']
const MONTHS = ['Janvier','Février','Mars','Avril','Mai','Juin','Juillet','Août','Septembre','Octobre','Novembre','Décembre']

export default function CalendarPage() {
  const navigate = useNavigate()
  const [year, setYear] = useState(new Date().getFullYear())
  const [month, setMonth] = useState(new Date().getMonth() + 1)

  const { data: days, isLoading } = useQuery({
    queryKey: ['calendar', year, month],
    queryFn: () => analyticsApi.calendar(year, month),
  })

  const sessionDays = new Map(days?.filter(d => d.has_session).map(d => [d.date, d]) ?? [])
  const today = format(new Date(), 'yyyy-MM-dd')

  const firstDay = startOfMonth(new Date(year, month - 1))
  const startPad = (getDay(firstDay) + 6) % 7
  const daysInMonth = getDaysInMonth(new Date(year, month - 1))

  const prev = () => {
    if (month === 1) { setYear(y => y - 1); setMonth(12) }
    else setMonth(m => m - 1)
  }
  const next = () => {
    if (month === 12) { setYear(y => y + 1); setMonth(1) }
    else setMonth(m => m + 1)
  }

  const cells: (number | null)[] = [
    ...Array(startPad).fill(null),
    ...Array.from({ length: daysInMonth }, (_, i) => i + 1),
  ]

  while (cells.length % 7 !== 0) cells.push(null)

  return (
    <div className="p-8 max-w-3xl">
      <div className="flex items-center justify-between mb-8">
        <h2 className="text-3xl font-bold text-white">Calendrier</h2>
        <div className="flex items-center gap-4">
          <button onClick={prev} className="text-zinc-400 hover:text-white transition-colors p-2">◀</button>
          <span className="text-lg font-semibold text-white min-w-[200px] text-center">
            {MONTHS[month - 1]} {year}
          </span>
          <button onClick={next} className="text-zinc-400 hover:text-white transition-colors p-2">▶</button>
        </div>
      </div>

      <div className="bg-zinc-900 border border-zinc-800 rounded-2xl p-6">
        <div className="grid grid-cols-7 mb-4">
          {DAYS.map(d => (
            <div key={d} className="text-center text-xs font-semibold text-zinc-500 uppercase tracking-wider py-2">
              {d}
            </div>
          ))}
        </div>

        {isLoading ? (
          <div className="text-center py-12 text-zinc-500 text-sm">Chargement...</div>
        ) : (
          <div className="grid grid-cols-7 gap-1">
            {cells.map((day, i) => {
              if (!day) return <div key={i} />
              const dateStr = `${year}-${String(month).padStart(2,'0')}-${String(day).padStart(2,'0')}`
              const session = sessionDays.get(dateStr)
              const isToday = dateStr === today
              const hasSession = !!session

              return (
                <button
                  key={i}
                  onClick={() => session?.session_id && navigate(`/workout`)}
                  className={`
                    relative aspect-square flex flex-col items-center justify-center rounded-xl text-sm font-medium transition-all
                    ${hasSession
                      ? 'bg-brand-500 text-white hover:bg-brand-600'
                      : isToday
                      ? 'bg-zinc-800 text-brand-400 ring-1 ring-brand-500/50'
                      : 'text-zinc-400 hover:bg-zinc-800 hover:text-white'
                    }
                  `}
                >
                  <span>{day}</span>
                  {hasSession && session?.duration_minutes && (
                    <span className="text-[9px] opacity-75 leading-none">{session.duration_minutes}m</span>
                  )}
                </button>
              )
            })}
          </div>
        )}
      </div>

      <div className="mt-6 flex items-center gap-6">
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-brand-500" />
          <span className="text-sm text-zinc-400">Jour d'entraînement</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-4 h-4 rounded bg-zinc-800 ring-1 ring-brand-500/50" />
          <span className="text-sm text-zinc-400">Aujourd'hui</span>
        </div>
      </div>

      {days && (
        <div className="mt-6 bg-zinc-900 border border-zinc-800 rounded-xl p-5">
          <h3 className="text-sm font-semibold text-zinc-400 mb-3">Résumé du mois</h3>
          <div className="flex gap-8">
            <div>
              <p className="text-2xl font-bold text-white">{sessionDays.size}</p>
              <p className="text-xs text-zinc-500">Jours d'entraînement</p>
            </div>
            <div>
              <p className="text-2xl font-bold text-white">
                {Array.from(sessionDays.values()).reduce((sum, d) => sum + (d.duration_minutes ?? 0), 0)}
              </p>
              <p className="text-xs text-zinc-500">Minutes totales</p>
            </div>
            <div>
              <p className="text-2xl font-bold text-brand-400">
                {Array.from(sessionDays.values()).reduce((sum, d) => sum + (d.calories_burned ?? 0), 0)}
              </p>
              <p className="text-xs text-zinc-500">Calories totales</p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
