package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func main() {
	todos := []Todo{}

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		todo := Todo{}
		err := json.NewDecoder(r.Body).Decode(&todo)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todos = append(todos, todo)
		w.WriteHeader(http.StatusCreated)
	})

    r.Get("/todos", func(w http.ResponseWriter, r *http.Request){
        w.Header().Set("Content-Type", "application/json")
        err := json.NewEncoder(w).Encode(todos)

        if err != nil {
            w.WriteHeader(http.)
        }
    })

	srv := &http.Server{
		Addr:    ":8090",
		Handler: r,
	}

	//Create a channel to listen for OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Println("Server is running on :8090")

		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	fmt.Println("Press Ctrl+C to stop the server")
	<-sigCh

	fmt.Println("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("Sever gracefully stopped")
}
