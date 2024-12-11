package api

import (
	"doc-server/internal/api/handlers/auth"
	docs "doc-server/internal/api/handlers/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine, authHandler *auth.AuthHandler, docsHandler *docs.DocsHandler) {
	authGroup := router.Group("/api")
	{
		authGroup.POST("/register", authHandler.Register) // Регистрация
		authGroup.POST("/auth", authHandler.Auth)         // Регистрация
		authGroup.POST("/docs", docsHandler.Save)         // Сохранение файла
		authGroup.GET("/docs", docsHandler.GetList)       // Список документов
		authGroup.GET("/docs/:id", docsHandler.Get)       // Получение документа
		authGroup.DELETE("/docs/:id", docsHandler.Delete) // Удаление документа

	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
