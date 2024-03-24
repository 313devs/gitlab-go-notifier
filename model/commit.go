package model

type Commit struct {
	Sha     string `json:"sha"`
	Message string `json:"message"`
	Author  string `json:"author"`
	PushedAt string `json:"pushed_at"`
}