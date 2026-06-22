# FinalOps

Application de suivi fitness personnelle — journal d'entraînement, progression des charges, et calendrier de séances.

---

## Fonctionnalités

- **Journal de séances** — log quotidien des entraînements avec date, durée, calories brûlées (saisie manuelle) et notes
- **Suivi par série** — poids (kg), répétitions, RPE (1–10), temps de repos pour chaque exercice
- **Catalogue d'exercices** — 20 exercices pré-chargés (développé couché, squat, tractions, dips…) extensibles à volonté
- **Programmes** — support de tous les types : Split, PPL, Upper/Lower, Full Body
- **Calendrier mensuel** — visualisation des jours d'entraînement
- **Graphiques de progression** — évolution du poids max et du volume total par exercice
- **Stats hebdo/mensuel** — séances, durée totale, calories
- **Authentification JWT** — login/register avec tokens access + refresh

---

## Stack technique

| Couche | Technologie |
|--------|-------------|
| Backend | Go 1.23 — 4 microservices |
| Router | chi v5 |
| Base de données | PostgreSQL 16 (1 instance, 3 schémas) |
| Cache | Redis 7 |
| Frontend | React 18 + TypeScript + Vite |
| UI | TailwindCSS (thème dark navy/bleu) |
| State | Zustand (auth) + TanStack Query (serveur) |
| Charts | Recharts |
| Conteneurs | Docker + Docker Compose |
| Orchestration | Kubernetes (k3s) |
| GitOps | ArgoCD |
| CI/CD | GitHub Actions |
| Bootstrap | Ansible |
| Registry | GitHub Container Registry (GHCR) |

---

## Architecture microservices

```
Browser
   │
   ▼
api-gateway (nginx :80)
   │
   ├── /api/auth/      → auth-service      :8001
   ├── /api/exercises  → exercise-service  :8002
   ├── /api/workouts/  → workout-service   :8003
   ├── /api/analytics/ → analytics-service :8004
   └── /              → frontend           :3000
                                │
                         PostgreSQL + Redis
```

### Les 4 services

| Service | Port | Rôle |
|---------|------|------|
| `auth-service` | 8001 | Register, login, JWT access/refresh, profil utilisateur |
| `exercise-service` | 8002 | Catalogue d'exercices — CRUD + seed 20 exercices |
| `workout-service` | 8003 | Sessions, exercices par séance, séries, programmes |
| `analytics-service` | 8004 | Calendrier, progression, stats (stateless — appelle workout-service via HTTP interne) |

### Base de données

Un seul cluster PostgreSQL, 3 schémas séparés avec des utilisateurs dédiés :

```
finalops (database)
├── auth     → auth_user     (users, profiles, refresh_tokens)
├── exercise → exercise_user (exercises, categories, muscle_groups)
└── workout  → workout_user  (sessions, session_exercises, sets, programs)
```

---

## Structure du projet

```
finalops/
├── services/
│   ├── auth-service/           # Go — cmd/ internal/ pkg/ migrations/
│   ├── exercise-service/       # Go — + seed/exercises_seed.sql
│   ├── workout-service/        # Go
│   ├── analytics-service/      # Go
│   └── api-gateway/            # nginx — nginx.conf + Dockerfile
├── frontend/                   # React + Vite + Tailwind
│   └── src/
│       ├── api/                # Clients HTTP (axios)
│       ├── pages/              # 6 pages : Login, Dashboard, Log Séance, Calendrier, Progression, Profil
│       ├── components/         # Layout, Sidebar
│       ├── store/              # Zustand (auth)
│       └── types/              # Types TypeScript partagés
├── infrastructure/
│   ├── kubernetes/             # Manifests k8s (Deployments, Services, Ingress, Postgres, Redis)
│   └── ansible/
│       ├── inventory/          # Inventaire homelab
│       └── playbooks/          # 3 playbooks : k3s, ArgoCD, deploy
├── argocd/                     # App-of-apps GitOps
├── .github/workflows/          # CI/CD — 1 workflow par service
├── scripts/                    # dev-up.sh, dev-down.sh
├── docker-compose.yml          # Environnement de développement local
├── docker-compose.prod.yml     # Overrides production
├── DEPLOY.md                   # Guide de déploiement homelab complet
└── .env.example                # Variables d'environnement requises
```

---

## Lancer en local (développement)

### Prérequis
- Docker Desktop (ou Docker + Docker Compose)
- Git

### Démarrage

```bash
git clone https://github.com/Rahman-git7/finalops.git
cd finalops

# (Optionnel) personnaliser le JWT secret
cp .env.example .env
# édite .env et change JWT_SECRET

docker compose up --build
```

L'app est accessible sur **http://localhost**

> Premier lancement : Docker télécharge les images de base (~2 min). Les migrations SQL et le seed des exercices sont exécutés automatiquement.

### Ports locaux

| Service | Port |
|---------|------|
| App (gateway) | http://localhost |
| Auth API | http://localhost:8001 |
| Exercise API | http://localhost:8002 |
| Workout API | http://localhost:8003 |
| Analytics API | http://localhost:8004 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

### Arrêter

```bash
docker compose down          # conserve les données
docker compose down -v       # supprime aussi les volumes (reset DB)
```

---

## API Reference

### Auth — `/api/auth/`

| Méthode | Route | Description |
|---------|-------|-------------|
| POST | `/auth/register` | Créer un compte `{email, username, password}` |
| POST | `/auth/login` | Login `{email, password}` → `{access_token, refresh_token}` |
| POST | `/auth/refresh` | Renouveler les tokens `{refresh_token}` |
| POST | `/auth/logout` | Invalider le refresh token |
| GET | `/auth/me` | Infos du compte connecté |
| GET/PUT | `/auth/profile` | Profil physique (poids, âge, sexe) |

### Exercises — `/api/exercises`

| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/exercises` | Liste (filtres: `?q=`, `?category=`, `?muscle=`) |
| POST | `/exercises` | Créer un exercice custom |
| GET/PUT/DELETE | `/exercises/:id` | CRUD sur un exercice |
| GET | `/exercises/categories` | Liste des catégories |
| GET | `/exercises/muscle-groups` | Liste des groupes musculaires |

### Workouts — `/api/workouts/`

| Méthode | Route | Description |
|---------|-------|-------------|
| POST/GET | `/workouts/sessions` | Créer / lister les séances |
| GET/PUT/DELETE | `/workouts/sessions/:id` | Opérations sur une séance |
| POST/DELETE | `/workouts/sessions/:id/exercises` | Ajouter/retirer un exercice |
| POST/PUT/DELETE | `/workouts/sessions/:id/exercises/:ex_id/sets` | CRUD séries |
| GET | `/workouts/program-types` | Types de programme (Split, PPL…) |
| CRUD | `/workouts/programs` | Gestion des programmes |

### Analytics — `/api/analytics/`

| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/analytics/calendar?year=&month=` | Jours d'entraînement du mois |
| GET | `/analytics/progression/:exercise_id` | Courbe poids/volume par exercice |
| GET | `/analytics/stats/weekly` | Stats de la semaine courante |
| GET | `/analytics/stats/monthly?year=&month=` | Stats mensuelles |

---

## Déploiement homelab (Kubernetes)

### Prérequis serveur
- Ubuntu 22.04 LTS
- 2+ Go RAM, 2+ CPU, 20+ Go disque
- Accès SSH depuis ta machine locale

### Pipeline CI/CD (automatique)

Chaque push sur `main` déclenche la chaîne suivante :

```
git push main
    │
    ▼
GitHub Actions
    ├── go vet (ou npm build pour le frontend)
    ├── docker build --target prod
    ├── docker push ghcr.io/rahman-git7/finalops/<service>:sha-xxxxx
    └── patch deployment.yaml avec le nouveau SHA
            │
            ▼
        ArgoCD détecte le changement dans git
            │
            ▼
        kubectl rollout (nouveau pod déployé)
```

### Déploiement initial

Voir **[DEPLOY.md](DEPLOY.md)** pour le guide complet. En résumé :

```bash
# 1. Éditer l'IP du serveur
vim infrastructure/ansible/inventory/homelab.ini

# 2. Installer k3s
ansible-playbook -i infrastructure/ansible/inventory/homelab.ini \
  infrastructure/ansible/playbooks/01-install-k3s.yml

# 3. Installer ArgoCD + créer les secrets
ansible-playbook -i infrastructure/ansible/inventory/homelab.ini \
  infrastructure/ansible/playbooks/02-install-argocd.yml

# 4. Déployer via GitOps
kubectl apply -f argocd/app-of-apps.yaml
```

### Images Docker

Les images sont publiées sur GHCR après chaque push sur `main` :

```
ghcr.io/rahman-git7/finalops/auth-service:latest
ghcr.io/rahman-git7/finalops/exercise-service:latest
ghcr.io/rahman-git7/finalops/workout-service:latest
ghcr.io/rahman-git7/finalops/analytics-service:latest
ghcr.io/rahman-git7/finalops/frontend:latest
```

---

## Développement

### Structure d'un service Go

Chaque service suit le même pattern :

```
service/
├── cmd/server/main.go          # Point d'entrée, routing chi
├── internal/
│   ├── config/config.go        # Chargement des variables d'env
│   ├── handler/                # Handlers HTTP (decode req, encode resp)
│   ├── service/                # Logique métier
│   ├── repository/             # Requêtes SQL (pgx)
│   └── model/                  # Structs de données
├── pkg/middleware/auth.go      # Middleware JWT partagé
├── migrations/                 # Fichiers SQL
└── Dockerfile                  # Multi-stage: golang:1.23-alpine → scratch
```

### Ajouter un exercice

Édite `services/exercise-service/seed/exercises_seed.sql` et relance :

```bash
docker compose down -v && docker compose up --build
```

### Variables d'environnement

| Variable | Service | Description |
|----------|---------|-------------|
| `JWT_SECRET` | Tous | Secret partagé pour signer les JWT |
| `DB_DSN` | auth, exercise, workout | Connection string PostgreSQL |
| `EXERCISE_SERVICE_URL` | workout | URL interne de l'exercise-service |
| `WORKOUT_SERVICE_URL` | analytics | URL interne du workout-service |
| `REDIS_URL` | analytics | URL Redis pour le cache |
| `CACHE_TTL` | analytics | Durée de cache (ex: `5m`) |

---

## Roadmap

- [ ] Tests d'intégration Go avec `testcontainers-go`
- [ ] Monitoring : Prometheus + Grafana
- [ ] Backup automatique PostgreSQL
- [ ] HTTPS via cert-manager (Let's Encrypt)
- [ ] Notifications PR/deploy via Webhook
