import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { authApi } from '../api/auth'
import { workoutsApi } from '../api/workouts'
import { useAuthStore } from '../store/authStore'

export default function ProfilePage() {
  const qc = useQueryClient()
  const { user } = useAuthStore()
  const [weightKg, setWeightKg] = useState('')
  const [age, setAge] = useState('')
  const [sex, setSex] = useState('')
  const [saved, setSaved] = useState(false)

  const { data: profile } = useQuery({
    queryKey: ['profile'],
    queryFn: authApi.getProfile,
  })

  const { data: programTypes } = useQuery({
    queryKey: ['program-types'],
    queryFn: workoutsApi.programTypes,
  })

  const { data: programs } = useQuery({
    queryKey: ['programs'],
    queryFn: workoutsApi.listPrograms,
  })

  useEffect(() => {
    if (profile) {
      setWeightKg(profile.weight_kg?.toString() ?? '')
      setAge(profile.age?.toString() ?? '')
      setSex(profile.sex ?? '')
    }
  }, [profile])

  const updateMut = useMutation({
    mutationFn: () => authApi.updateProfile({
      weight_kg: weightKg ? Number(weightKg) : undefined,
      age: age ? Number(age) : undefined,
      sex: sex || undefined,
    }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['profile'] })
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    },
  })

  const [newProgramName, setNewProgramName] = useState('')
  const [newProgramType, setNewProgramType] = useState<number | ''>('')
  const [newProgramNotes, setNewProgramNotes] = useState('')

  const createProgramMut = useMutation({
    mutationFn: () => workoutsApi.createProgram({
      name: newProgramName,
      program_type_id: Number(newProgramType),
      notes: newProgramNotes || undefined,
    }),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['programs'] })
      setNewProgramName(''); setNewProgramType(''); setNewProgramNotes('')
    },
  })

  const deleteProgramMut = useMutation({
    mutationFn: (id: string) => workoutsApi.deleteProgram(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['programs'] }),
  })

  return (
    <div className="p-8 max-w-2xl">
      <h2 className="text-3xl font-bold text-white mb-8">Profil</h2>

      <section className="bg-slate-900 border border-slate-800 rounded-2xl p-6 mb-6">
        <h3 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">Informations du compte</h3>
        <div className="space-y-1">
          <div className="flex justify-between py-2 border-b border-slate-800">
            <span className="text-sm text-slate-500">Nom d'utilisateur</span>
            <span className="text-sm font-medium text-white">{user?.username}</span>
          </div>
          <div className="flex justify-between py-2">
            <span className="text-sm text-slate-500">Email</span>
            <span className="text-sm font-medium text-white">{user?.email}</span>
          </div>
        </div>
      </section>

      <section className="bg-slate-900 border border-slate-800 rounded-2xl p-6 mb-6">
        <h3 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">Données physiques</h3>
        <div className="grid grid-cols-3 gap-4 mb-4">
          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1.5">Poids (kg)</label>
            <input
              type="number" value={weightKg} onChange={e => setWeightKg(e.target.value)} placeholder="75"
              className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-brand-500"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1.5">Âge</label>
            <input
              type="number" value={age} onChange={e => setAge(e.target.value)} placeholder="25"
              className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-brand-500"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-slate-400 mb-1.5">Sexe</label>
            <select
              value={sex} onChange={e => setSex(e.target.value)}
              className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2.5 text-sm text-white focus:outline-none focus:border-brand-500"
            >
              <option value="">—</option>
              <option value="male">Homme</option>
              <option value="female">Femme</option>
              <option value="other">Autre</option>
            </select>
          </div>
        </div>
        <button
          onClick={() => updateMut.mutate()}
          disabled={updateMut.isPending}
          className="bg-brand-500 hover:bg-brand-600 disabled:opacity-50 text-white font-semibold px-5 py-2 rounded-lg text-sm transition-colors"
        >
          {saved ? '✓ Sauvegardé' : 'Sauvegarder'}
        </button>
      </section>

      <section className="bg-slate-900 border border-slate-800 rounded-2xl p-6">
        <h3 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-4">Programmes d'entraînement</h3>

        <div className="space-y-2 mb-4">
          {programs?.length === 0 && (
            <p className="text-slate-600 text-sm">Aucun programme créé.</p>
          )}
          {programs?.map(p => (
            <div key={p.id} className="flex items-center justify-between py-3 border-b border-slate-800 last:border-0">
              <div>
                <p className="text-sm font-medium text-white">{p.name}</p>
                <p className="text-xs text-slate-500">{p.program_type ?? `Type ${p.program_type_id}`}</p>
              </div>
              <button
                onClick={() => deleteProgramMut.mutate(p.id)}
                className="text-xs text-slate-600 hover:text-red-400 transition-colors"
              >
                Supprimer
              </button>
            </div>
          ))}
        </div>

        <div className="space-y-3 border-t border-slate-800 pt-4">
          <p className="text-xs font-medium text-slate-400">Nouveau programme</p>
          <div className="grid grid-cols-2 gap-3">
            <input
              value={newProgramName} onChange={e => setNewProgramName(e.target.value)} placeholder="Nom du programme"
              className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-brand-500"
            />
            <select
              value={newProgramType} onChange={e => setNewProgramType(Number(e.target.value))}
              className="bg-slate-800 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-brand-500"
            >
              <option value="">Type de programme</option>
              {programTypes?.map(t => (
                <option key={t.id} value={t.id}>{t.name}</option>
              ))}
            </select>
          </div>
          <input
            value={newProgramNotes} onChange={e => setNewProgramNotes(e.target.value)} placeholder="Notes (optionnel)"
            className="w-full bg-slate-800 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:outline-none focus:border-brand-500"
          />
          <button
            onClick={() => createProgramMut.mutate()}
            disabled={!newProgramName || !newProgramType || createProgramMut.isPending}
            className="bg-brand-500 hover:bg-brand-600 disabled:opacity-50 text-white font-semibold px-5 py-2 rounded-lg text-sm transition-colors"
          >
            Créer le programme
          </button>
        </div>
      </section>
    </div>
  )
}
