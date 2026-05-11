import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { exercisesApi } from '../api/exercises'
import { analyticsApi } from '../api/analytics'
import {
  LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend,
} from 'recharts'
import { format } from 'date-fns'
import { fr } from 'date-fns/locale'

const CustomTooltip = ({ active, payload, label }: any) => {
  if (!active || !payload?.length) return null
  return (
    <div className="bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-xs">
      <p className="text-slate-400 mb-1">{label}</p>
      {payload.map((p: any) => (
        <p key={p.name} style={{ color: p.color }}>
          {p.name}: <span className="font-bold text-white">{p.value}</span>
        </p>
      ))}
    </div>
  )
}

export default function ProgressPage() {
  const [selectedExercise, setSelectedExercise] = useState<string>('')

  const { data: exercises } = useQuery({
    queryKey: ['exercises-all'],
    queryFn: () => exercisesApi.list(),
  })

  const { data: progression, isLoading } = useQuery({
    queryKey: ['progression', selectedExercise],
    queryFn: () => analyticsApi.progression(selectedExercise),
    enabled: !!selectedExercise,
  })

  const chartData = progression?.map(p => ({
    date: format(new Date(p.date), 'd MMM', { locale: fr }),
    'Poids max (kg)': p.max_weight ?? null,
    'Volume total (kg)': Math.round(p.total_volume),
    'Reps max': p.max_reps ?? null,
    'RPE moyen': p.avg_rpe ? Number(p.avg_rpe.toFixed(1)) : null,
  })) ?? []

  const selectedEx = exercises?.find(e => e.id === selectedExercise)

  return (
    <div className="p-8">
      <h2 className="text-3xl font-bold text-white mb-8">Progression</h2>

      <div className="mb-6">
        <label className="block text-xs font-medium text-slate-400 mb-2">Choisir un exercice</label>
        <select
          value={selectedExercise}
          onChange={e => setSelectedExercise(e.target.value)}
          className="bg-slate-900 border border-slate-700 rounded-lg px-4 py-2.5 text-sm text-white focus:outline-none focus:border-brand-500 min-w-[280px]"
        >
          <option value="">-- Sélectionner --</option>
          {exercises?.map(ex => (
            <option key={ex.id} value={ex.id}>{ex.name}</option>
          ))}
        </select>
      </div>

      {!selectedExercise && (
        <div className="flex flex-col items-center justify-center py-24 text-center">
          <div className="text-5xl mb-4">📈</div>
          <p className="text-slate-500">Sélectionne un exercice pour voir ta progression</p>
        </div>
      )}

      {selectedExercise && isLoading && (
        <p className="text-slate-500 text-sm">Chargement des données...</p>
      )}

      {selectedExercise && !isLoading && progression && (
        <>
          {progression.length === 0 ? (
            <div className="bg-slate-900 border border-slate-800 rounded-xl p-8 text-center">
              <p className="text-slate-500">Aucune donnée pour cet exercice.</p>
              <p className="text-slate-600 text-xs mt-1">Logger des séances pour voir ta progression.</p>
            </div>
          ) : (
            <div className="space-y-6">
              {selectedEx && (
                <div className="flex items-center gap-4 mb-2">
                  <h3 className="text-xl font-bold text-white">{selectedEx.name}</h3>
                  <span className="text-xs bg-brand-500/10 text-brand-400 border border-brand-500/20 rounded-full px-3 py-1">
                    {selectedEx.primary_muscle_name}
                  </span>
                </div>
              )}

              <div className="grid grid-cols-3 gap-4">
                {(() => {
                  const last = progression[progression.length - 1]
                  const prev = progression.length > 1 ? progression[progression.length - 2] : null
                  const weightDiff = prev?.max_weight && last.max_weight
                    ? last.max_weight - prev.max_weight : null
                  return (
                    <>
                      <div className="bg-slate-900 border border-slate-800 rounded-xl p-5">
                        <p className="text-xs text-slate-500 uppercase tracking-wider mb-1">Poids max actuel</p>
                        <p className="text-3xl font-bold text-white">{last.max_weight ?? '—'} <span className="text-base text-slate-500 font-normal">kg</span></p>
                        {weightDiff !== null && (
                          <p className={`text-xs mt-1 font-medium ${weightDiff >= 0 ? 'text-green-400' : 'text-red-400'}`}>
                            {weightDiff >= 0 ? '+' : ''}{weightDiff.toFixed(1)} kg vs séance précédente
                          </p>
                        )}
                      </div>
                      <div className="bg-slate-900 border border-slate-800 rounded-xl p-5">
                        <p className="text-xs text-slate-500 uppercase tracking-wider mb-1">Volume dernière séance</p>
                        <p className="text-3xl font-bold text-white">{last.total_volume} <span className="text-base text-slate-500 font-normal">kg</span></p>
                      </div>
                      <div className="bg-slate-900 border border-slate-800 rounded-xl p-5">
                        <p className="text-xs text-slate-500 uppercase tracking-wider mb-1">Nombre de séances</p>
                        <p className="text-3xl font-bold text-white">{progression.length}</p>
                      </div>
                    </>
                  )
                })()}
              </div>

              <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                <h4 className="text-sm font-semibold text-slate-400 mb-4">Évolution du poids (kg)</h4>
                <ResponsiveContainer width="100%" height={240}>
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                    <XAxis dataKey="date" tick={{ fill: '#64748b', fontSize: 11 }} />
                    <YAxis tick={{ fill: '#64748b', fontSize: 11 }} />
                    <Tooltip content={<CustomTooltip />} />
                    <Line
                      type="monotone"
                      dataKey="Poids max (kg)"
                      stroke="#3b82f6"
                      strokeWidth={2}
                      dot={{ fill: '#3b82f6', r: 4 }}
                      connectNulls
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>

              <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                <h4 className="text-sm font-semibold text-slate-400 mb-4">Volume total par séance (kg)</h4>
                <ResponsiveContainer width="100%" height={200}>
                  <LineChart data={chartData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
                    <XAxis dataKey="date" tick={{ fill: '#64748b', fontSize: 11 }} />
                    <YAxis tick={{ fill: '#64748b', fontSize: 11 }} />
                    <Tooltip content={<CustomTooltip />} />
                    <Line
                      type="monotone"
                      dataKey="Volume total (kg)"
                      stroke="#3b82f6"
                      strokeWidth={2}
                      dot={{ fill: '#3b82f6', r: 3 }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  )
}
