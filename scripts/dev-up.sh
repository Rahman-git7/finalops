#!/bin/bash
set -e

echo "🚀 Démarrage de FinalOps..."
docker compose up --build -d

echo ""
echo "✅ FinalOps est en cours de démarrage !"
echo ""
echo "  App:        http://localhost"
echo "  Auth API:   http://localhost:8001"
echo "  Exercises:  http://localhost:8002"
echo "  Workouts:   http://localhost:8003"
echo "  Analytics:  http://localhost:8004"
echo ""
echo "Logs: docker compose logs -f"
