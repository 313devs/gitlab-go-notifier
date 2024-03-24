package handler

import (
	"net/http"
)

type Commit struct {
}

func (c *Commit) PostCommit(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("committed"))
}
func (c *Commit) GetCommits(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{all commits}"))
}
