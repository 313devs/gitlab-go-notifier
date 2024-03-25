package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// App is the main application struct
type App struct {
	router http.Handler
	rdb    *redis.Client
}

func New() *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	rdAddr := os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT")
	rdPassword := os.Getenv("REDIS_PASSWORD")
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr:    rdAddr,
			Password: rdPassword,
			DB:       0,
		}),
	}

	app.loadRoutes()

	return app
}
func (a *App) Start(ctx context.Context) error {

	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Printf("failed to close redis connection: %v\n", err)
		}
	}()

	ch := make(chan error, 1)

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	fmt.Println("Starting server on port", server.Addr, "and connected to redis")
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to listen to server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return server.Shutdown(timeoutCtx)
	}
}
