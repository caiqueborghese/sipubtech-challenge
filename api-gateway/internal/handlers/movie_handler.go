package handlers

import (
	"net/http"
	"strconv"

	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/domain"
	"github.com/caiqueborghese/sipubtech-challenge/api-gateway/internal/usecase"
	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	svc usecase.MovieService
}

func RegisterMovieRoutes(r *gin.Engine, svc usecase.MovieService) {
	h := &MovieHandler{svc: svc}
	g := r.Group("/movies")
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.DELETE("/:id", h.Delete)
}

// List godoc
// @Summary Lista todos os filmes
// @Tags movies
// @Produce json
// @Param limit query int false "Máximo de itens retornados (default 50, max 200)"
// @Success 200 {array} domain.Movie
// @Router /movies [get]
func (h *MovieHandler) List(c *gin.Context) {
	movies, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	// limit (default=50, max=200)
	limit := 50
	if s := c.Query("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			if v < 1 {
				v = 1
			}
			if v > 200 {
				v = 200
			}
			limit = v
		}
	}
	if len(movies) > limit {
		movies = movies[:limit]
	}

	c.JSON(http.StatusOK, movies)
}

// Get godoc
// @Summary Busca um filme por ID
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Success 200 {object} domain.Movie
// @Failure 404 {string} string "movie not found"
// @Failure 400 {string} string "invalid id"
// @Router /movies/{id} [get]
func (h *MovieHandler) Get(c *gin.Context) {
	id := c.Param("id")
	m, err := h.svc.Get(id)
	if err != nil {
		status := http.StatusBadGateway
		if err == domain.ErrNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// Create godoc
// @Summary Cria um novo filme
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body domain.Movie true "Movie"
// @Success 201 {object} domain.Movie
// @Failure 400 {string} string "invalid body"
// @Router /movies [post]
func (h *MovieHandler) Create(c *gin.Context) {
	var in domain.Movie
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	m, err := h.svc.Create(&in)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

// Delete godoc
// @Summary Remove um filme
// @Tags movies
// @Param id path string true "Movie ID"
// @Success 204
// @Failure 404 {string} string "movie not found"
// @Router /movies/{id} [delete]
func (h *MovieHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
