package domain

import "time"

type Task struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Processed bool      `json:"processed"`
	NeedRetry bool      `json:"need_retry"`
	Error     string    `json:"error"`
	Url       string    `json:"url"`
	Status    string    `json:"status"`
}
