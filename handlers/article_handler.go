package handlers

import (
	"net/http"
	"strconv"

	"github.com/brothergiez/restful-api/repositories"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	Repo *repositories.ArticleRepository
}

func NewArticleHandler(repo *repositories.ArticleRepository) *ArticleHandler {
	return &ArticleHandler{
		Repo: repo,
	}
}

func (h *ArticleHandler) CreateArticleHandler(c *gin.Context) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	article := h.Repo.CreateArticle(input.Title, input.Content)
	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) UpdateArticleHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: Title and Content are required"})
		return
	}

	article, err := h.Repo.UpdateArticle(id, input.Title, input.Content)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) SearchArticlesHandler(c *gin.Context) {
	keyword := c.Query("keyword")
	articles := h.Repo.SearchArticles(keyword)
	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetAllArticlesHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	articles, total := h.Repo.GetAllArticlesWithPagination(page, limit)

	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, gin.H{
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": totalPages,
		"articles":   articles,
	})
}
