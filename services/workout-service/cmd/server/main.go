package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/finalops/workout-service/internal/config"
	"github.com/finalops/workout-service/internal/handler"
	"github.com/finalops/workout-service/internal/repository"
	"github.com/finalops/workout-service/internal/service"
	"github.com/finalops/workout-service/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DBDSN)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("cannot ping database: %v", err)
	}
	log.Println("connected to database")

	repo := repository.NewWorkoutRepository(db)
	svc := service.NewWorkoutService(repo)
	h := handler.NewWorkoutHandler(svc)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	r.Get("/health", h.Health)

	// Internal endpoints (no JWT — called only by analytics-service within Docker network)
	r.Get("/internal/sessions/range", h.InternalSessionRange)
	r.Get("/internal/sessions/exercise/{exercise_id}", h.InternalExerciseHistory)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))

		r.Get("/workouts/program-types", h.ListProgramTypes)

		r.Post("/workouts/programs", h.CreateProgram)
		r.Get("/workouts/programs", h.ListPrograms)
		r.Get("/workouts/programs/{id}", h.GetProgram)
		r.Put("/workouts/programs/{id}", h.UpdateProgram)
		r.Delete("/workouts/programs/{id}", h.DeleteProgram)

		r.Post("/workouts/sessions", h.CreateSession)
		r.Get("/workouts/sessions", h.ListSessions)
		r.Get("/workouts/sessions/{id}", h.GetSession)
		r.Put("/workouts/sessions/{id}", h.UpdateSession)
		r.Delete("/workouts/sessions/{id}", h.DeleteSession)

		r.Post("/workouts/sessions/{id}/exercises", h.AddExercise)
		r.Delete("/workouts/sessions/{id}/exercises/{exercise_id}", h.RemoveExercise)

		r.Post("/workouts/sessions/{id}/exercises/{exercise_id}/sets", h.AddSet)
		r.Put("/workouts/sessions/{id}/exercises/{exercise_id}/sets/{set_id}", h.UpdateSet)
		r.Delete("/workouts/sessions/{id}/exercises/{exercise_id}/sets/{set_id}", h.DeleteSet)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("workout-service listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("workout-service shutdown")
}
