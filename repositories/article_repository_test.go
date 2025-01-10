package repositories

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	repo := NewArticleRepository()

	article := repo.CreateArticle("Test Title", "Test Content")
	assert.Equal(t, 1, article.ID)
	assert.Equal(t, "Test Title", article.Title)
	assert.Equal(t, "Test Content", article.Content)

	assert.Len(t, repo.articles, 1)
	assert.Equal(t, repo.articles[0], article)
}

func TestUpdateArticle(t *testing.T) {
	repo := NewArticleRepository()
	article := repo.CreateArticle("Test Title", "Test Content")

	updatedArticle, err := repo.UpdateArticle(article.ID, "Updated Title", "Updated Content")
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updatedArticle.Title)
	assert.Equal(t, "Updated Content", updatedArticle.Content)

	_, err = repo.UpdateArticle(999, "New Title", "New Content")
	assert.Error(t, err)
	assert.Equal(t, "article not found", err.Error())
}

func TestSearchArticles(t *testing.T) {
	repo := NewArticleRepository()
	repo.CreateArticle("First Article", "Content of the first article")
	repo.CreateArticle("Second Article", "Content of the second article")
	repo.CreateArticle("Another Post", "Completely unrelated content")

	results := repo.SearchArticles("article")
	assert.Len(t, results, 2)
	assert.Equal(t, "First Article", results[0].Title)
	assert.Equal(t, "Second Article", results[1].Title)

	results = repo.SearchArticles("unrelated")
	assert.Len(t, results, 1)
	assert.Equal(t, "Another Post", results[0].Title)

	results = repo.SearchArticles("nonexistent")
	assert.Len(t, results, 0)
}

func TestGetAllArticlesWithPagination(t *testing.T) {
	repo := NewArticleRepository()
	for i := 1; i <= 15; i++ {
		repo.CreateArticle("Title "+strconv.Itoa(i), "Content "+strconv.Itoa(i))
	}

	results, total := repo.GetAllArticlesWithPagination(1, 5)
	assert.Len(t, results, 5)
	assert.Equal(t, 15, total)
	assert.Equal(t, "Title 1", results[0].Title)
	assert.Equal(t, "Title 5", results[4].Title)

	results, _ = repo.GetAllArticlesWithPagination(3, 5)
	assert.Len(t, results, 5)
	assert.Equal(t, "Title 11", results[0].Title)
	assert.Equal(t, "Title 15", results[4].Title)

	results, total = repo.GetAllArticlesWithPagination(4, 5)
	assert.Len(t, results, 0)
	assert.Equal(t, 15, total)

	results, _ = repo.GetAllArticlesWithPagination(2, 10)
	assert.Len(t, results, 5)
	assert.Equal(t, "Title 11", results[0].Title)
	assert.Equal(t, "Title 15", results[4].Title)
}
