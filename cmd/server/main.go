package main

import (
	"context"
	"errors"
	"go-mentor-gpt/internal/httpapi"
	"go-mentor-gpt/internal/storage/sqlite"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	storage, err := sqlite.NewStorageSQLite("/data/mentor.db")
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()
	handlers := httpapi.NewHandlers(storage)

	token := os.Getenv("API_TOKEN")
	if token == "" {
		log.Fatal("API_TOKEN environment variable not set")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /state", handlers.GetState)
	mux.HandleFunc("POST /topics", handlers.CreateTopic)
	mux.HandleFunc("PATCH /topics/{id}", handlers.UpdateTopic)
	mux.HandleFunc("POST /tasks", handlers.CreateTask)
	mux.HandleFunc("PATCH /tasks/{id}", handlers.UpdateTask)
	mux.HandleFunc("PATCH /goal", handlers.SetGoal)

	server := http.Server{
		Addr:    ":8080",
		Handler: httpapi.Auth(token, mux),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGKILL, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("shutdown server: %v", err)
	}
}
