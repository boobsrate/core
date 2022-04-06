package domain

import (
	"time"
)

type Vote struct {
	TitsID    string    `json:"tits_id"`
	CreatedAt time.Time `json:"created_at"`
}
