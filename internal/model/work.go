package model

import (
	"errors"
	"strings"
	"time"
)

type Work struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	User_ID     int       `json:"user_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiredAt   time.Time `json:"exp_at"`
}

func (w *Work) BeforeCreate() error {
	if len(w.Name) > 0 {
		w.CreatedAt = time.Now()
		w.ExpiredAt = w.CreatedAt.Add(time.Hour * 24)
	} else {
		return errors.New("len name equal 0, add name")
	}
	return nil
}

func (w *Work) IsEmpty() bool {
	return strings.Trim(w.Name, " ") == "" && strings.Trim(w.Description, " ") == ""
}
