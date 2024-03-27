package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/313devs/gitlab-go-notifier/handler"
	"github.com/313devs/gitlab-go-notifier/middleware"
	"github.com/313devs/gitlab-go-notifier/repository/commit"
)

func (a *App) loadRoutes(){
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Auth)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	router.Route("/commit", a.loadCommitRoutes)

	a.router = router
}
func (a *App) loadCommitRoutes(router chi.Router) {
	commitHandler := handler.Commit{
		Repo: &commit.RedisRepo{
			Client: a.rdb,
		},
	}
	router.Get("/", commitHandler.GetCommits)
	router.Post("/", commitHandler.PostCommit)
}
