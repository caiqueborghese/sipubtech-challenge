package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gdomain "github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type fakeSvc struct {
	list []gdomain.Movie
	get  *gdomain.Movie
	err  error
}

func (f *fakeSvc) List() ([]gdomain.Movie, error) { return f.list, f.err }
func (f *fakeSvc) Get(id string) (*gdomain.Movie, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.get, nil
}
func (f *fakeSvc) Create(m *gdomain.Movie) (*gdomain.Movie, error) {
	if f.err != nil {
		return nil, f.err
	}
	m.ID = "new"
	return m, nil
}
func (f *fakeSvc) Delete(id string) error { return f.err }

var _ usecase.MovieService = (*fakeSvc)(nil)

func setupRouter(svc usecase.MovieService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterMovieRoutes(r, svc)
	return r
}

func TestListHandler_OK(t *testing.T) {
	svc := &fakeSvc{list: []gdomain.Movie{{ID: "8", Title: "X", Year: 2000}}}
	r := setupRouter(svc)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies", nil)
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var got []gdomain.Movie
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	require.Len(t, got, 1)
	require.Equal(t, "8", got[0].ID)
}

func TestCreateHandler_Valid(t *testing.T) {
	svc := &fakeSvc{}
	r := setupRouter(svc)

	body := `{"title":"Hello","year":2025}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/movies", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var m gdomain.Movie
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &m))
	require.Equal(t, "new", m.ID)
}
