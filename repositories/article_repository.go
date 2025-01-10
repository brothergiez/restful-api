package repositories

import (
	"errors"
	"strings"

	"github.com/brothergiez/restful-api/models"
)

type ArticleRepository struct {
	articles []models.Article
	nextID   int
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{
		articles: []models.Article{},
		nextID:   1,
	}
}

func (r *ArticleRepository) CreateArticle(title, content string) models.Article {
	article := models.Article{
		ID:      r.nextID,
		Title:   title,
		Content: content,
	}
	r.articles = append(r.articles, article)
	r.nextID++
	return article
}

func (r *ArticleRepository) UpdateArticle(id int, title, content string) (models.Article, error) {
	for i, article := range r.articles {
		if article.ID == id {
			r.articles[i].Title = title
			r.articles[i].Content = content
			return r.articles[i], nil
		}
	}

	return models.Article{}, errors.New("article not found")
}

func (r *ArticleRepository) SearchArticles(keyword string) []models.Article {
	keyword = strings.ToLower(keyword)
	result := []models.Article{}
	for _, article := range r.articles {
		if strings.Contains(strings.ToLower(article.Title), keyword) || strings.Contains(strings.ToLower(article.Content), keyword) {
			result = append(result, article)
		}
	}

	return result
}

func (r *ArticleRepository) GetAllArticlesWithPagination(page, limit int) ([]models.Article, int) {
	start := (page - 1) * limit
	end := start + limit

	if start >= len(r.articles) {
		return []models.Article{}, len(r.articles)
	}
	if end > len(r.articles) {
		end = len(r.articles)
	}

	return r.articles[start:end], len(r.articles)
}
