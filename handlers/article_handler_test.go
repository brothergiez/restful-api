package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brothergiez/restful-api/models"
	"github.com/brothergiez/restful-api/repositories"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateArticleHandler(t *testing.T) {
	repo := repositories.NewArticleRepository()
	handler := NewArticleHandler(repo)

	router := gin.Default()
	router.POST("/articles", handler.CreateArticleHandler)

	payload := `{"title":"Test Title","content":"Test Content"}`
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var article models.Article
	err := json.Unmarshal(resp.Body.Bytes(), &article)
	assert.NoError(t, err)
	assert.Equal(t, "Test Title", article.Title)
	assert.Equal(t, "Test Content", article.Content)
	assert.Equal(t, 1, article.ID)
}

func TestUpdateArticleHandler(t *testing.T) {
	repo := repositories.NewArticleRepository()
	article := repo.CreateArticle("Original Title", "Original Content")
	handler := NewArticleHandler(repo)

	router := gin.Default()
	router.PUT("/articles/:id", handler.UpdateArticleHandler)

	payload := `{"title":"Updated Title","content":"Updated Content"}`
	url := "/articles/" + strconv.Itoa(article.ID)
	req := httptest.NewRequest(http.MethodPut, url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var updatedArticle models.Article
	err := json.Unmarshal(resp.Body.Bytes(), &updatedArticle)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedArticle.Title)
	assert.Equal(t, "Updated Content", updatedArticle.Content)

	req = httptest.NewRequest(http.MethodPut, "/articles/999", bytes.NewBufferString(payload))
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestSearchArticlesHandler(t *testing.T) {
	repo := repositories.NewArticleRepository()
	repo.CreateArticle("First Article", "Content of the first article")
	repo.CreateArticle("Second Article", "Content of the second article")
	handler := NewArticleHandler(repo)

	router := gin.Default()
	router.GET("/articles/search", handler.SearchArticlesHandler)

	req := httptest.NewRequest(http.MethodGet, "/articles/search?keyword=article", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var articles []models.Article
	err := json.Unmarshal(resp.Body.Bytes(), &articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 2)
	assert.Equal(t, "First Article", articles[0].Title)
	assert.Equal(t, "Second Article", articles[1].Title)
}

func TestGetAllArticlesHandler(t *testing.T) {
	repo := repositories.NewArticleRepository()
	for i := 1; i <= 15; i++ {
		repo.CreateArticle("Title "+strconv.Itoa(i), "Content "+strconv.Itoa(i))
	}
	handler := NewArticleHandler(repo)

	router := gin.Default()
	router.GET("/articles/get-all", handler.GetAllArticlesHandler)

	req := httptest.NewRequest(http.MethodGet, "/articles/get-all?page=1&limit=5", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result struct {
		Page       int              `json:"page"`
		Limit      int              `json:"limit"`
		Total      int              `json:"total"`
		TotalPages int              `json:"totalPages"`
		Articles   []models.Article `json:"articles"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, 15, result.Total)
	assert.Equal(t, 3, result.TotalPages)
	assert.Len(t, result.Articles, 5)
	assert.Equal(t, "Title 1", result.Articles[0].Title)
	assert.Equal(t, "Title 5", result.Articles[4].Title)

	req = httptest.NewRequest(http.MethodGet, "/articles/get-all?page=2&limit=5", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, 15, result.Total)
	assert.Equal(t, 3, result.TotalPages)
	assert.Len(t, result.Articles, 5)
	assert.Equal(t, "Title 6", result.Articles[0].Title)
	assert.Equal(t, "Title 10", result.Articles[4].Title)

	req = httptest.NewRequest(http.MethodGet, "/articles/get-all?page=4&limit=5", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	err = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, 4, result.Page)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, 15, result.Total)
	assert.Equal(t, 3, result.TotalPages)
	assert.Len(t, result.Articles, 0)
}
