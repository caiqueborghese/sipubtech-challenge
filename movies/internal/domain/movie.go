package domain

import (
	"errors"
	"strings"
)

type Movie struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	LegacyID string `bson:"legacy_id,omitempty" json:"-"`
}

var (
	ErrInvalidID  = errors.New("invalid id")
	ErrNotFound   = errors.New("movie not found")
	ErrValidation = errors.New("validation error")
)

func (m *Movie) Normalize() {
	m.Title = strings.TrimSpace(m.Title)
}

// Validação mínima de domínio usada no usecase.
func (m *Movie) Validate() error {
	if strings.TrimSpace(m.Title) == "" {
		return errors.New("title required")
	}

	if m.Year < 1800 || m.Year > 3000 {
		return errors.New("year out of range")
	}
	return nil
}
