package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/313devs/gitlab-go-notifier/model"
	"github.com/313devs/gitlab-go-notifier/repository/commit"
)

type Commit struct {
	Repo *commit.RedisRepo
}

func (c *Commit) PostCommit(w http.ResponseWriter, r *http.Request) {
	var Body struct {
		Sha     string `json:"sha"`
		Message string `json:"message"`
		Author  string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}
	now := time.Now().UTC()
	commit := model.Commit{
		Sha:      Body.Sha,
		Message:  Body.Message,
		Author:   Body.Author,
		PushedAt: &now,
	}
	err := c.Repo.Insert(r.Context(), commit)
	if err != nil {
		fmt.Printf("failed to insert commit: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to insert commit"))
		return
	}

	// Send the SHA to Telegram
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")
	message := fmt.Sprintf("New feature added: %s", Body.Message)

	if strings.Contains(Body.Message, "feat") {
		sendTelegramMessage(telegramToken, telegramChatID, message)
	}

	res, err := json.Marshal(commit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshal commit"))
		return
	}
	w.Write(res)
	w.WriteHeader(http.StatusCreated)
}

func (c *Commit) GetCommits(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{all commits}"))
}

func sendTelegramMessage(token, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
	data := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
