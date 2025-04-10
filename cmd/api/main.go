package main

import (
	"context"
	"fmt"
	"log"
	"lostnfound-api/internal/config"
	"lostnfound-api/internal/handler"
	"lostnfound-api/internal/repository"
	"lostnfound-api/internal/router"
	"lostnfound-api/internal/service"
	"lostnfound-api/internal/util/storage"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up database
	db, err := repository.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Google Cloud Storage
	gcs, err := storage.NewGoogleCloudStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Google Cloud Storage: %v", err)
	}
	defer gcs.Close()

	// Initialize repositories
	itemRepo := repository.NewItemRepository(db)

	// Initialize services
	itemService := service.NewItemService(itemRepo)

	// Initialize handlers
	itemHandler := handler.NewItemHandler(itemService)

	// Setup router
	r := router.SetupRouter(cfg, itemHandler)

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	// Server run context
	go func() {
		log.Printf("Server starting on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
