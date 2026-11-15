// Package main adalah titik masuk aplikasi SIM Dina Catering.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dinacatering/internal/config"
	"dinacatering/internal/database"
	"dinacatering/internal/handler"
	"dinacatering/internal/middleware"
	"dinacatering/internal/seed"
	"dinacatering/internal/service"
	"dinacatering/internal/store"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)
	log.SetPrefix("[sim-dina] ")

	if err := run(); err != nil {
		log.Fatalf("aplikasi berhenti: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := database.Migrate(cfg.DatabaseURL); err != nil {
		return err
	}
	log.Println("migrasi selesai")

	if cfg.SeedOnBoot {
		if err := seed.Run(ctx, pool, cfg.AdminPassword); err != nil {
			return err
		}
		log.Println("data awal siap")
	}

	queries := store.New(pool)
	svc := service.New(queries, pool)
	sessions := middleware.NewSessionManager(cfg.SessionSecret)

	h, err := handler.New(svc, sessions)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           h.Router(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("aplikasi berjalan di http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server gagal: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("mematikan server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
