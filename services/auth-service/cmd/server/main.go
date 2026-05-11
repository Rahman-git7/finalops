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

	"github.com/finalops/auth-service/internal/config"
	"github.com/finalops/auth-service/internal/handler"
	"github.com/finalops/auth-service/internal/middleware"
	"github.com/finalops/auth-service/internal/repository"
	"github.com/finalops/auth-service/internal/service"
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

	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo, cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	h := handler.NewAuthHandler(svc)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", h.Health)

	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	r.Post("/auth/refresh", h.Refresh)
	r.Post("/auth/logout", h.Logout)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Get("/auth/me", h.Me)
		r.Get("/auth/profile", h.GetProfile)
		r.Put("/auth/profile", h.UpdateProfile)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("auth-service listening on :%s", cfg.Port)
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
	log.Println("auth-service shutdown")
}
