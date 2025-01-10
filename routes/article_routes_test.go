package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brothergiez/restful-api/handlers"
	"github.com/brothergiez/restful-api/repositories"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	// Buat repository dan handler untuk digunakan dalam tes
	repo := repositories.NewArticleRepository()
	handler := handlers.NewArticleHandler(repo)

	// Inisialisasi router dan daftar route
	router := gin.Default()
	RegisterArticleRoutes(router, handler)
	return router
}

func TestRegisterArticleRoutes_CreateArticle(t *testing.T) {
	router := setupRouter()

	// Simulasi request untuk /articles/create
	payload := `{"title":"Test Title","content":"Test Content"}`
	req := httptest.NewRequest(http.MethodPost, "/articles/create", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Validasi response
	assert.Equal(t, http.StatusCreated, resp.Code)

	var article map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &article)
	assert.NoError(t, err)
	assert.Equal(t, "Test Title", article["title"])
	assert.Equal(t, "Test Content", article["content"])
	assert.Equal(t, float64(1), article["id"])
}

func TestRegisterArticleRoutes_UpdateArticle(t *testing.T) {
	router := setupRouter()

	// Buat artikel terlebih dahulu
	payload := `{"title":"Original Title","content":"Original Content"}`
	req := httptest.NewRequest(http.MethodPost, "/articles/create", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var article map[string]interface{}
	_ = json.Unmarshal(resp.Body.Bytes(), &article)
	id := int(article["id"].(float64))

	// Simulasi request untuk /articles/update/:id
	updatePayload := `{"title":"Updated Title","content":"Updated Content"}`
	req = httptest.NewRequest(http.MethodPut, "/articles/update/"+strconv.Itoa(id), bytes.NewBufferString(updatePayload))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Validasi response
	assert.Equal(t, http.StatusOK, resp.Code)

	var updatedArticle map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &updatedArticle)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedArticle["title"])
	assert.Equal(t, "Updated Content", updatedArticle["content"])
}

func TestRegisterArticleRoutes_SearchArticles(t *testing.T) {
	router := setupRouter()

	// Buat beberapa artikel
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/articles/create", bytes.NewBufferString(`{"title":"First Article","content":"Content 1"}`)))
	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/articles/create", bytes.NewBufferString(`{"title":"Second Article","content":"Content 2"}`)))

	// Simulasi request untuk /articles/search
	req := httptest.NewRequest(http.MethodGet, "/articles/search?keyword=Article", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Validasi response
	assert.Equal(t, http.StatusOK, resp.Code)

	var articles []map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &articles)
	assert.NoError(t, err)
	assert.Len(t, articles, 2)
	assert.Equal(t, "First Article", articles[0]["title"])
	assert.Equal(t, "Second Article", articles[1]["title"])
}

func TestRegisterArticleRoutes_GetAllArticles(t *testing.T) {
	router := setupRouter()

	// Buat beberapa artikel
	for i := 1; i <= 15; i++ {
		payload := `{"title":"Title ` + strconv.Itoa(i) + `","content":"Content ` + strconv.Itoa(i) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/articles/create", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}

	// Simulasi request untuk /articles/get-all
	req := httptest.NewRequest(http.MethodGet, "/articles/get-all?page=1&limit=5", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	// Validasi response
	assert.Equal(t, http.StatusOK, resp.Code)

	var result struct {
		Page       int                      `json:"page"`
		Limit      int                      `json:"limit"`
		Total      int                      `json:"total"`
		TotalPages int                      `json:"totalPages"`
		Articles   []map[string]interface{} `json:"articles"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 5, result.Limit)
	assert.Equal(t, 15, result.Total)
	assert.Equal(t, 3, result.TotalPages)
	assert.Len(t, result.Articles, 5)
	assert.Equal(t, "Title 1", result.Articles[0]["title"])
	assert.Equal(t, "Title 5", result.Articles[4]["title"])
}
