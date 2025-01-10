package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brothergiez/restful-api/handlers"
	"github.com/brothergiez/restful-api/middlewares"
	"github.com/brothergiez/restful-api/repositories"
	"github.com/brothergiez/restful-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	repo := repositories.NewArticleRepository()
	handler := handlers.NewArticleHandler(repo)

	router := gin.Default()
	router.Use(middlewares.LoggingMiddleware())
	routes.RegisterArticleRoutes(router, handler)

	return router
}

func TestMainServer(t *testing.T) {
	router := setupTestRouter()

	payload := `{"title":"Test Title","content":"Test Content"}`
	req := httptest.NewRequest(http.MethodPost, "/articles/create", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var article map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &article)
	assert.NoError(t, err)
	assert.Equal(t, "Test Title", article["title"])
	assert.Equal(t, "Test Content", article["content"])
	assert.Equal(t, float64(1), article["id"])
}

func TestMainEnvironmentVariable(t *testing.T) {
	os.Setenv("APP_PORT", "5000")
	defer os.Unsetenv("APP_PORT")

	port := os.Getenv("APP_PORT")
	assert.Equal(t, "5000", port)
}

func TestMainInvalidRoute(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/invalid-route", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
