package model

import "time"

type Work struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Exp         time.Time `json:"exp_at"`
}
