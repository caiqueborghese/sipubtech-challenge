package seed

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/caiqueborghese/sipubtech-challenge/movies/internal/domain"
)

type rawMovie struct {
	ID    any    `json:"id"`
	Title string `json:"title"`
	Year  any    `json:"year"`
}

// LoadSeed lê o arquivo JSON e converte para []domain.Movie
func LoadSeed(path string) ([]domain.Movie, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read seed: %w", err)
	}

	// usar decoder com UseNumber para preservar números quando possível
	dec := json.NewDecoder(strings.NewReader(string(data)))
	dec.UseNumber()

	var raws []rawMovie
	if err := dec.Decode(&raws); err != nil {
		return nil, fmt.Errorf("unmarshal seed: %w", err)
	}

	movies := make([]domain.Movie, 0, len(raws))
	for _, r := range raws {
		legacyID := anyToString(r.ID)
		yearInt := anyToInt(r.Year)

		movies = append(movies, domain.Movie{
			Title:    strings.TrimSpace(r.Title),
			Year:     yearInt,
			LegacyID: legacyID,
		})
	}
	return movies, nil
}

func anyToString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return strings.TrimSpace(t)
	case json.Number:
		return t.String()
	case float64:
		return strconv.Itoa(int(t))
	case int:
		return strconv.Itoa(t)
	case int32:
		return strconv.Itoa(int(t))
	case int64:
		return strconv.FormatInt(t, 10)
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

func anyToInt(v any) int {
	switch t := v.(type) {
	case nil:
		return 0
	case string:
		t = strings.TrimSpace(t)
		if t == "" {
			return 0
		}
		if n, err := strconv.Atoi(t); err == nil {
			return n
		}
		return 0
	case json.Number:
		if n, err := t.Int64(); err == nil {
			return int(n)
		}
		return 0
	case float64:
		return int(t)
	case int:
		return t
	case int32:
		return int(t)
	case int64:
		return int(t)
	default:
		if n, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(t))); err == nil {
			return n
		}
		return 0
	}
}
