package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	gdomain "github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type nfSvc struct{ usecase.MovieService }

func (nf nfSvc) Get(id string) (*gdomain.Movie, error) {
	return nil, gdomain.ErrNotFound
}

func setupNFRouter(svc usecase.MovieService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterMovieRoutes(r, svc)
	return r
}

func TestGet_NotFound_MapsTo404(t *testing.T) {
	r := setupNFRouter(nfSvc{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/does-not-exist", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusNotFound, w.Code)
}
