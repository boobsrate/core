package domain

import "time"

type Tits struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Rating    int64     `json:"rating"`
	URL       string    `json:"url"`
	FullURL   string    `json:"full_url"`
}
