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

	"github.com/finalops/exercise-service/internal/config"
	"github.com/finalops/exercise-service/internal/handler"
	"github.com/finalops/exercise-service/internal/repository"
	"github.com/finalops/exercise-service/internal/service"
	"github.com/finalops/exercise-service/pkg/middleware"
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

	repo := repository.NewExerciseRepository(db)
	svc := service.NewExerciseService(repo)
	h := handler.NewExerciseHandler(svc)

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

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Get("/exercises", h.List)
		r.Post("/exercises", h.Create)
		r.Get("/exercises/categories", h.ListCategories)
		r.Get("/exercises/muscle-groups", h.ListMuscleGroups)
		r.Get("/exercises/{id}", h.Get)
		r.Put("/exercises/{id}", h.Update)
		r.Delete("/exercises/{id}", h.Delete)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("exercise-service listening on :%s", cfg.Port)
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
	log.Println("exercise-service shutdown")
}
