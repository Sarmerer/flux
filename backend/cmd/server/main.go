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

	columnsvc "github.com/example/flow/internal/app/column"
	columnshttp "github.com/example/flow/internal/app/column/adapters/http"
	columnspersistence "github.com/example/flow/internal/app/column/adapters/persistence"
	relationshipsvc "github.com/example/flow/internal/app/relationship"
	relationshiphttp "github.com/example/flow/internal/app/relationship/adapters/http"
	relationshippersistence "github.com/example/flow/internal/app/relationship/adapters/persistence"
	rowhttp "github.com/example/flow/internal/app/row/adapters/http"
	tablesvc "github.com/example/flow/internal/app/table"
	tablehttp "github.com/example/flow/internal/app/table/adapters/http"
	tablepersistence "github.com/example/flow/internal/app/table/adapters/persistence"
	"github.com/example/flow/internal/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := getenv("PORT", "8080")

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:4173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	cfg := common.LoadConfig()
	db, err := common.NewDBPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	if err := common.RunMigrations(context.Background(), db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	wsHub := common.NewWSHub(cfg.WebOrigin)
	r.Get("/ws", wsHub.HandleWS)

	tableRepo := tablepersistence.NewRepo(db)
	tableService := tablesvc.NewService(tableRepo)
	tableHandler := tablehttp.NewHandler(tableService, wsHub)

	r.Mount("/api/tables", tableHandler.Routes())

	colRepo := columnspersistence.NewRepo(db)
	colSvc := columnsvc.NewService(colRepo)
	colHandler := columnshttp.NewHandler(colSvc, wsHub)
	r.Mount("/api/columns", colHandler.Routes())

	relRepo := relationshippersistence.NewRepo(db)
	relSvc := relationshipsvc.NewService(relRepo)
	relHandler := relationshiphttp.NewHandler(relSvc, wsHub)
	r.Mount("/api/relationships", relHandler.Routes())

	rowHandler := rowhttp.NewHandler(db, wsHub)
	r.Mount("/api/rows", rowHandler.Routes())

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	fmt.Println("server stopped")
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
