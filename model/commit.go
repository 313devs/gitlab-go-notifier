package model

import "time"

type Commit struct {
	Sha     string `json:"sha"`
	Message string `json:"message"`
	Author  string `json:"author"`
	PushedAt *time.Time `json:"pushed_at"`
}