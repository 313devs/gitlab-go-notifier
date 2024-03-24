package application

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", func (w http.ResponseWriter, r *http.Request)  {
		w.Write([]byte("Hello, World!"))
	})
	router.Route("/commit", loadCommitRoutes)
	

	return router
}
func loadCommitRoutes(router chi.Router) {
	router.Get("/", func (w http.ResponseWriter, r *http.Request)  {
	})
	router.Post("/", func (w http.ResponseWriter, r *http.Request)  {
	})
}