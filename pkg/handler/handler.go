package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/bookCRUD/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/login", h.login)
		auth.POST("/logout", h.logout)
		auth.POST("/signup", h.signup)
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
