package domain

import (
	"errors"
	"strings"
)

type Movie struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

var (
	ErrNotFound   = errors.New("movie not found")
	ErrInvalidID  = errors.New("invalid id")
	ErrValidation = errors.New("validation error")
)

func (m *Movie) Normalize() {
	m.Title = strings.TrimSpace(m.Title)
}

func (m *Movie) Validate() error {
	if strings.TrimSpace(m.Title) == "" {
		return errors.Join(ErrValidation, errors.New("title is required"))
	}
	if m.Year < 1878 || m.Year > 3000 {
		return errors.Join(ErrValidation, errors.New("year is out of range"))
	}
	return nil
}
