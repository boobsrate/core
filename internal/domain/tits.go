package domain

import "time"

type Tits struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Rating    int64     `json:"rating"`
	URL       string    `json:"url"`
	FullURL   string    `json:"full_url"`
	Abyss     bool      `json:"abyss"`
}

type Vote struct {
	TitsID    string    `json:"tits_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Report struct {
	TitsID    string    `json:"tits_id"`
	CreatedAt time.Time `json:"created_at"`
}
