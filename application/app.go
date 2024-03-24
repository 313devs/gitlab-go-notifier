package application

import (
	"fmt"
	"context"
	"net/http"
)

// App is the main application struct
type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: loadRoutes(),
	}
	return app
}
func (a *App) Start(ctx context.Context) error {

	server := &http.Server{
		Addr:    ":8080",
		Handler: a.router,
	}
	fmt.Println("Starting server on port", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
