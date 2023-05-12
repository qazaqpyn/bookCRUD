package handler

import (
	"github.com/gin-gonic/gin"
	_ "github.com/qazaqpyn/bookCRUD/docs"
	"github.com/qazaqpyn/bookCRUD/pkg/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/signup", h.signup)
		auth.POST("/login", h.login)
		auth.GET("/refresh", h.refresh)
	}

	api := router.Group("/api", h.userIdentity)
	{
		books := api.Group("/books")
		{
			books.GET("/", h.getBooks)
			books.GET("/:id", h.getBook)
			books.POST("/", h.createBook)
			books.PUT("/:id", h.updateBook)
			books.DELETE("/:id", h.deleteBook)
		}
	}
	return router
}
