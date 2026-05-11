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

	"github.com/finalops/analytics-service/internal/client"
	"github.com/finalops/analytics-service/internal/config"
	"github.com/finalops/analytics-service/internal/handler"
	"github.com/finalops/analytics-service/internal/service"
	"github.com/finalops/analytics-service/pkg/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatalf("invalid redis URL: %v", err)
	}
	rdb := redis.NewClient(opt)
	defer rdb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("cannot connect to redis: %v", err)
	}
	log.Println("connected to redis")

	wc := client.NewWorkoutClient(cfg.WorkoutServiceURL)
	svc := service.NewAnalyticsService(wc, rdb, cfg.CacheTTL)
	h := handler.NewAnalyticsHandler(svc)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	r.Get("/health", h.Health)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Get("/analytics/calendar", h.Calendar)
		r.Get("/analytics/progression/{exercise_id}", h.Progression)
		r.Get("/analytics/stats/weekly", h.WeeklyStats)
		r.Get("/analytics/stats/monthly", h.MonthlyStats)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("analytics-service listening on :%s", cfg.Port)
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
	log.Println("analytics-service shutdown")
}
