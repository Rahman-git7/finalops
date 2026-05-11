import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { workoutsApi } from '../api/workouts'
import { exercisesApi } from '../api/exercises'
import { format } from 'date-fns'
import type { Session, Exercise, Set } from '../types'

function SetRow({
  set,
  sessionId,
  exerciseId,
  onUpdate,
  onDelete,
}: {
  set: Set
  sessionId: string
  exerciseId: string
  onUpdate: () => void
  onDelete: () => void
}) {
  const qc = useQueryClient()
  const [editing, setEditing] = useState(false)
  const [weight, setWeight] = useState(set.weight_kg?.toString() ?? '')
  const [reps, setReps] = useState(set.reps?.toString() ?? '')
  const [rpe, setRpe] = useState(set.rpe?.toString() ?? '')
  const [rest, setRest] = useState(set.rest_seconds?.toString() ?? '')

  const updateMut = useMutation({
    mutationFn: () => workoutsApi.updateSet(sessionId, exerciseId, set.id, {
      weight_kg: weight ? Number(weight) : undefined,
      reps: reps ? Number(reps) : undefined,
      rpe: rpe ? Number(rpe) : undefined,
      rest_seconds: rest ? Number(rest) : undefined,
    }),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['session'] }); setEditing(false); onUpdate() },
  })

  const deleteMut = useMutation({
    mutationFn: () => workoutsApi.deleteSet(sessionId, exerciseId, set.id),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['session'] }); onDelete() },
  })

  if (!editing) {
    return (
      <div
        className="grid grid-cols-5 gap-2 text-sm py-2 px-3 rounded-lg hover:bg-zinc-800/50 cursor-pointer group"
        onClick={() => setEditing(true)}
      >
        <span className="text-zinc-500">#{set.set_number}</span>
        <span className="text-white font-medium">{set.weight_kg ? `${set.weight_kg} kg` : '—'}</span>
        <span className="text-white font-medium">{set.reps ? `${set.reps} reps` : '—'}</span>
        <span className="text-zinc-400">{set.rpe ? `RPE ${set.rpe}` : '—'}</span>
        <span className="text-zinc-400">{set.rest_seconds ? `${set.rest_seconds}s` : '—'}</span>
      </div>
    )
  }

  return (
    <div className="grid grid-cols-5 gap-2 py-2 px-3 bg-zinc-800/50 rounded-lg">
      <span className="text-zinc-500 text-sm self-center">#{set.set_number}</span>
      <input value={weight} onChange={e => setWeight(e.target.value)} placeholder="kg"
        className="bg-zinc-700 border border-zinc-600 rounded px-2 py-1 text-sm text-white focus:outline-none focus:border-brand-500" />
      <input value={reps} onChange={e => setReps(e.target.value)} placeholder="reps"
        className="bg-zinc-700 border border-zinc-600 rounded px-2 py-1 text-sm text-white focus:outline-none focus:border-brand-500" />
      <input value={rpe} onChange={e => setRpe(e.target.value)} placeholder="RPE" step="0.5" min="1" max="10"
        className="bg-zinc-700 border border-zinc-600 rounded px-2 py-1 text-sm text-white focus:outline-none focus:border-brand-500" />
      <div className="flex gap-1">
        <input value={rest} onChange={e => setRest(e.target.value)} placeholder="sec"
          className="flex-1 bg-zinc-700 border border-zinc-600 rounded px-2 py-1 text-sm text-white focus:outline-none focus:border-brand-500" />
        <button onClick={() => updateMut.mutate()} className="text-brand-400 hover:text-brand-300 px-1 text-xs">✓</button>
        <button onClick={() => deleteMut.mutate()} className="text-red-400 hover:text-red-300 px-1 text-xs">×</button>
      </div>
    </div>
  )
}

function ExerciseBlock({
  sessionExercise,
  session,
  exerciseMap,
}: {
  sessionExercise: Session['exercises'][0]
  session: Session
  exerciseMap: Map<string, Exercise>
}) {
  const qc = useQueryClient()
  const exercise = exerciseMap.get(sessionExercise.exercise_id)
  const [addingSet, setAddingSet] = useState(false)
  const [weight, setWeight] = useState('')
  const [reps, setReps] = useState('')
  const [rpe, setRpe] = useState('')
  const [rest, setRest] = useState('')

  const nextSetNumber = (sessionExercise.sets?.length ?? 0) + 1

  const addSetMut = useMutation({
    mutationFn: () => workoutsApi.addSet(session.id, sessionExercise.exercise_id, {
      set_number: nextSetNumber,
      weight_kg: weight ? Number(weight) : undefined,
      reps: reps ? Number(reps) : undefined,
      rpe: rpe ? Number(rpe) : undefined,
      rest_seconds: rest ? Number(rest) : undefined,
    }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['session', session.id] })
      setWeight(''); setReps(''); setRpe(''); setRest(''); setAddingSet(false)
    },
  })

  const removeExerciseMut = useMutation({
    mutationFn: () => workoutsApi.removeExercise(session.id, sessionExercise.exercise_id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['session', session.id] }),
  })

  return (
    <div className="bg-zinc-900 border border-zinc-800 rounded-xl p-5">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h4 className="font-semibold text-white">{exercise?.name ?? 'Exercice inconnu'}</h4>
          {exercise && (
            <p className="text-xs text-zinc-500 mt-0.5">
              {exercise.primary_muscle_name} · {exercise.category_name}
              {exercise.equipment && ` · ${exercise.equipment}`}
            </p>
          )}
        </div>
        <button
          onClick={() => removeExerciseMut.mutate()}
          className="text-zinc-600 hover:text-red-400 transition-colors text-lg leading-none"
        >
          ×
        </button>
      </div>

      {sessionExercise.sets?.length > 0 && (
        <div className="mb-3">
          <div className="grid grid-cols-5 gap-2 text-xs text-zinc-600 px-3 mb-1">
            <span>Série</span><span>Poids</span><span>Reps</span><span>RPE</span><span>Repos</span>
          </div>
          {sessionExercise.sets.map(set => (
            <SetRow
              key={set.id}
              set={set}
              sessionId={session.id}
              exerciseId={sessionExercise.exercise_id}
              onUpdate={() => {}}
              onDelete={() => {}}
            />
          ))}
        </div>
      )}

      {addingSet ? (
        <div className="grid grid-cols-5 gap-2 mt-2">
          <span className="text-zinc-500 text-sm self-center">#{nextSetNumber}</span>
          <input value={weight} onChange={e => setWeight(e.target.value)} placeholder="kg"
            className="bg-zinc-800 border border-zinc-700 rounded px-2 py-1.5 text-sm text-white focus:outline-none focus:border-brand-500" />
          <input value={reps} onChange={e => setReps(e.target.value)} placeholder="reps"
            className="bg-zinc-800 border border-zinc-700 rounded px-2 py-1.5 text-sm text-white focus:outline-none focus:border-brand-500" />
          <input value={rpe} onChange={e => setRpe(e.target.value)} placeholder="RPE"
            className="bg-zinc-800 border border-zinc-700 rounded px-2 py-1.5 text-sm text-white focus:outline-none focus:border-brand-500" />
          <div className="flex gap-1">
            <input value={rest} onChange={e => setRest(e.target.value)} placeholder="sec"
              className="flex-1 bg-zinc-800 border border-zinc-700 rounded px-2 py-1.5 text-sm text-white focus:outline-none focus:border-brand-500" />
            <button onClick={() => addSetMut.mutate()} className="text-brand-400 hover:text-brand-300 px-1 text-sm font-bold">+</button>
            <button onClick={() => setAddingSet(false)} className="text-zinc-500 hover:text-zinc-300 px-1 text-sm">×</button>
          </div>
        </div>
      ) : (
        <button
          onClick={() => setAddingSet(true)}
          className="mt-2 text-xs font-medium text-brand-500 hover:text-brand-400 transition-colors flex items-center gap-1"
        >
          + Ajouter une série
        </button>
      )}
    </div>
  )
}

function ExerciseSelector({
  sessionId,
  onClose,
}: {
  sessionId: string
  onClose: () => void
}) {
  const qc = useQueryClient()
  const [query, setQuery] = useState('')
  const { data: exercises } = useQuery({
    queryKey: ['exercises', query],
    queryFn: () => exercisesApi.list({ q: query }),
  })

  const addMut = useMutation({
    mutationFn: (exerciseId: string) => workoutsApi.addExercise(sessionId, exerciseId),
    onSuccess: () => { qc.invalidateQueries({ queryKey: ['session'] }); onClose() },
  })

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
      <div className="bg-zinc-900 border border-zinc-800 rounded-2xl w-full max-w-md max-h-[70vh] flex flex-col">
        <div className="p-5 border-b border-zinc-800">
          <div className="flex items-center justify-between mb-3">
            <h3 className="font-semibold text-white">Ajouter un exercice</h3>
            <button onClick={onClose} className="text-zinc-500 hover:text-white text-xl leading-none">×</button>
          </div>
          <input
            value={query}
            onChange={e => setQuery(e.target.value)}
            placeholder="Rechercher un exercice..."
            className="w-full bg-zinc-800 border border-zinc-700 rounded-lg px-3 py-2 text-sm text-white placeholder-zinc-600 focus:outline-none focus:border-brand-500"
            autoFocus
          />
        </div>
        <div className="overflow-y-auto flex-1 p-3 space-y-1">
          {exercises?.map(ex => (
            <button
              key={ex.id}
              onClick={() => addMut.mutate(ex.id)}
              className="w-full text-left px-3 py-3 rounded-lg hover:bg-zinc-800 transition-colors"
            >
              <p className="text-sm font-medium text-white">{ex.name}</p>
              <p className="text-xs text-zinc-500 mt-0.5">
                {ex.primary_muscle_name} · {ex.category_name}
              </p>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}

export default function WorkoutLogPage() {
  const navigate = useNavigate()
  const qc = useQueryClient()
  const [showNewSession, setShowNewSession] = useState(false)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)
  const [showExerciseSelector, setShowExerciseSelector] = useState(false)
  const [newDate, setNewDate] = useState(format(new Date(), 'yyyy-MM-dd'))
  const [newDuration, setNewDuration] = useState('')
  const [newCalories, setNewCalories] = useState('')
  const [newNotes, setNewNotes] = useState('')

  const { data: sessions, isLoading } = useQuery({
    queryKey: ['sessions'],
    queryFn: () => workoutsApi.listSessions({ limit: 20 }),
  })

  const { data: fullSession, refetch: refetchSession } = useQuery({
    queryKey: ['session', selectedSession?.id],
    queryFn: () => workoutsApi.getSession(selectedSession!.id),
    enabled: !!selectedSession?.id,
  })

  const { data: exercises } = useQuery({
    queryKey: ['exercises-all'],
    queryFn: () => exercisesApi.list(),
  })

  const exerciseMap = new Map((exercises ?? []).map(e => [e.id, e]))

  const createMut = useMutation({
    mutationFn: () => workoutsApi.createSession({
      session_date: newDate,
      duration_minutes: newDuration ? Number(newDuration) : undefined,
      calories_burned: newCalories ? Number(newCalories) : undefined,
      notes: newNotes || undefined,
    }),
    onSuccess: (session) => {
      qc.invalidateQueries({ queryKey: ['sessions'] })
      qc.invalidateQueries({ queryKey: ['weekly-stats'] })
      qc.invalidateQueries({ queryKey: ['monthly-stats'] })
      setShowNewSession(false)
      setSelectedSession(session)
      setNewDate(format(new Date(), 'yyyy-MM-dd'))
      setNewDuration(''); setNewCalories(''); setNewNotes('')
    },
  })

  const deleteMut = useMutation({
    mutationFn: (id: string) => workoutsApi.deleteSession(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['sessions'] })
      setSelectedSession(null)
    },
  })

  const activeSession = fullSession ?? selectedSession

  return (
    <div className="flex h-full">
      <div className="w-80 border-r border-zinc-800 flex flex-col">
        <div className="p-5 border-b border-zinc-800">
          <h2 className="font-bold text-white text-lg">Séances</h2>
        </div>
        <div className="p-3">
          <button
            onClick={() => setShowNewSession(true)}
            className="w-full bg-brand-500 hover:bg-brand-600 text-white font-semibold py-2.5 rounded-lg text-sm transition-colors"
          >
            + Nouvelle séance
          </button>
        </div>
        <div className="flex-1 overflow-y-auto px-3 pb-3 space-y-2">
          {isLoading && <p className="text-zinc-500 text-sm text-center py-8">Chargement...</p>}
          {sessions?.map(session => (
            <button
              key={session.id}
              onClick={() => setSelectedSession(session)}
              className={`w-full text-left p-4 rounded-xl border transition-all ${
                selectedSession?.id === session.id
                  ? 'bg-brand-500/10 border-brand-500/30 text-brand-400'
                  : 'bg-zinc-900 border-zinc-800 hover:border-zinc-700 text-white'
              }`}
            >
              <p className="text-sm font-semibold">{session.session_date}</p>
              <p className="text-xs text-zinc-500 mt-1">
                {session.duration_minutes ? `${session.duration_minutes} min` : 'Durée non définie'}
                {session.calories_burned ? ` · ${session.calories_burned} kcal` : ''}
              </p>
            </button>
          ))}
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-8">
        {!selectedSession && !showNewSession && (
          <div className="flex flex-col items-center justify-center h-full text-center">
            <div className="text-5xl mb-4">💪</div>
            <p className="text-zinc-500">Sélectionne une séance ou crée-en une nouvelle</p>
          </div>
        )}

        {showNewSession && (
          <div className="max-w-lg">
            <h3 className="text-xl font-bold text-white mb-6">Nouvelle séance</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1.5">Date</label>
                <input type="date" value={newDate} onChange={e => setNewDate(e.target.value)}
                  className="w-full bg-zinc-900 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-brand-500" />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1.5">Durée (min)</label>
                  <input type="number" value={newDuration} onChange={e => setNewDuration(e.target.value)} placeholder="60"
                    className="w-full bg-zinc-900 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-brand-500" />
                </div>
                <div>
                  <label className="block text-xs font-medium text-zinc-400 mb-1.5">Calories</label>
                  <input type="number" value={newCalories} onChange={e => setNewCalories(e.target.value)} placeholder="400"
                    className="w-full bg-zinc-900 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-brand-500" />
                </div>
              </div>
              <div>
                <label className="block text-xs font-medium text-zinc-400 mb-1.5">Notes</label>
                <textarea value={newNotes} onChange={e => setNewNotes(e.target.value)} rows={3} placeholder="Ressenti, PR, conditions..."
                  className="w-full bg-zinc-900 border border-zinc-700 rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-brand-500 resize-none" />
              </div>
              <div className="flex gap-3">
                <button onClick={() => createMut.mutate()} disabled={createMut.isPending}
                  className="flex-1 bg-brand-500 hover:bg-brand-600 disabled:opacity-50 text-white font-semibold py-2.5 rounded-lg text-sm transition-colors">
                  Créer la séance
                </button>
                <button onClick={() => setShowNewSession(false)}
                  className="px-4 bg-zinc-800 hover:bg-zinc-700 text-white rounded-lg text-sm transition-colors">
                  Annuler
                </button>
              </div>
            </div>
          </div>
        )}

        {selectedSession && !showNewSession && activeSession && (
          <div className="max-w-2xl">
            <div className="flex items-start justify-between mb-6">
              <div>
                <h3 className="text-2xl font-bold text-white">{activeSession.session_date}</h3>
                <div className="flex gap-4 mt-1">
                  {activeSession.duration_minutes && (
                    <span className="text-sm text-zinc-400">{activeSession.duration_minutes} min</span>
                  )}
                  {activeSession.calories_burned && (
                    <span className="text-sm text-brand-400">{activeSession.calories_burned} kcal</span>
                  )}
                </div>
                {activeSession.notes && (
                  <p className="text-sm text-zinc-500 mt-2 italic">"{activeSession.notes}"</p>
                )}
              </div>
              <button
                onClick={() => deleteMut.mutate(selectedSession.id)}
                className="text-xs text-zinc-600 hover:text-red-400 transition-colors border border-zinc-800 hover:border-red-400/30 rounded-lg px-3 py-1.5"
              >
                Supprimer
              </button>
            </div>

            <div className="space-y-3 mb-4">
              {(activeSession.exercises ?? []).map(se => (
                <ExerciseBlock
                  key={se.id}
                  sessionExercise={se}
                  session={activeSession}
                  exerciseMap={exerciseMap}
                />
              ))}
            </div>

            <button
              onClick={() => setShowExerciseSelector(true)}
              className="w-full border-2 border-dashed border-zinc-700 hover:border-brand-500/50 text-zinc-500 hover:text-brand-500 rounded-xl py-4 text-sm font-medium transition-all"
            >
              + Ajouter un exercice
            </button>
          </div>
        )}
      </div>

      {showExerciseSelector && selectedSession && (
        <ExerciseSelector
          sessionId={selectedSession.id}
          onClose={() => { setShowExerciseSelector(false); refetchSession() }}
        />
      )}
    </div>
  )
}
