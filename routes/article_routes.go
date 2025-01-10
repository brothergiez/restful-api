package routes

import (
	"github.com/brothergiez/restful-api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterArticleRoutes(router *gin.Engine, handler *handlers.ArticleHandler) {
	articleRoutes := router.Group("/articles")
	{
		articleRoutes.POST("/create", handler.CreateArticleHandler)
		articleRoutes.PUT("/update/:id", handler.UpdateArticleHandler)
		articleRoutes.GET("/search", handler.SearchArticlesHandler)
		articleRoutes.GET("/get-all", handler.GetAllArticlesHandler)
	}
}
