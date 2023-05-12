package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/bookCRUD/model"
)

// @Summary CreateBook
// @Description	Create a book
// @Security ApiKeyAuth
// @Tags books
// @Accept json
// @Produce	json
// @Param input body model.Book true "book details"
// @Success	200	{string} string "book created"
// @Failure	400 {object} errorResponse
// @Failure	404	{object}	errorResponse
// @Failure	500	{object} errorResponse
// @Router /api/books/ [post]
func (h *Handler) createBook(c *gin.Context) {
	_, ok := c.Get(userCtx)

	if !ok {
		logError("createBook", errors.New("user id not found"))

		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input model.Book
	if err := c.BindJSON(&input); err != nil {
		logError("createBook", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//service
	err := h.services.Create(c, input)
	if err != nil {
		logError("createBook", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"msg": "book created",
	})
}

// @Summary GetBook
// @Description	Get a book
// @Security ApiKeyAuth
// @Tags books
// @Accept json
// @Produce	json
// @Param input body model.Book true "book details"
// @Param id	path int true "Account ID"
// @Success	200	{object} model.Book
// @Failure	400 {object} errorResponse
// @Failure	404	{object}	errorResponse
// @Failure	500	{object} errorResponse
// @Router /api/books/{id} [get]
func (h *Handler) getBook(c *gin.Context) {
	_, ok := c.Get(userCtx)

	if !ok {
		logError("getBook", errors.New("user id not found"))
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input model.Book
	if err := c.BindJSON(&input); err != nil {
		logError("createBook", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//service
	book, err := h.services.GetById(c, input.Id)
	if err != nil {
		logError("createBook", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"msg": book,
	})
}

// @Summary GetBooks
// @Description	Get books
// @Security ApiKeyAuth
// @Tags books
// @Accept json
// @Produce	json
// @Success	200	{array} model.Book
// @Failure	400 {object} errorResponse
// @Failure	404	{object}	errorResponse
// @Failure	500	{object} errorResponse
// @Router /api/books/ [get]
func (h *Handler) getBooks(c *gin.Context) {
	_, ok := c.Get(userCtx)
	if !ok {
		logError("getBooks", errors.New("user id not found"))
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	books, err := h.services.GetAll(c)
	if err != nil {
		logError("getBooks", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"msg": &books,
	})
}

// @Summary UpdateBook
// @Description	Update a book
// @Security ApiKeyAuth
// @Tags books
// @Accept json
// @Produce	json
// @Param input body model.Book true "book details"
// @Param id	path int true "book ID"
// @Success	200	{string} string
// @Failure	400 {object} errorResponse
// @Failure	404	{object}	errorResponse
// @Failure	500	{object} errorResponse
// @Router /api/books/{id} [put]
func (h *Handler) updateBook(c *gin.Context) {
	_, ok := c.Get(userCtx)

	if !ok {
		logError("updateBook", errors.New("user id not found"))
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input model.Book
	if err := c.BindJSON(&input); err != nil {
		logError("updateBook", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//service
	err := h.services.Update(c, input.Id, &input)
	if err != nil {
		logError("updateBook", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book": "Book updated",
	})
}

// @Summary DeleteBook
// @Description	delete a book
// @Security ApiKeyAuth
// @Tags books
// @Accept json
// @Produce	json
// @Param id	path int true "book ID"
// @Success	200	{object} model.Book
// @Failure	400 {object} errorResponse
// @Failure	404	{object}	errorResponse
// @Failure	500	{object} errorResponse
// @Router /api/books/{id} [delete]
func (h *Handler) deleteBook(c *gin.Context) {
	_, ok := c.Get(userCtx)

	if !ok {
		logError("deleteBook", errors.New("user id not found"))
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return
	}

	var input model.Book
	if err := c.BindJSON(&input); err != nil {
		logError("deleteBook", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	//service
	book, err := h.services.GetById(c, input.Id)
	if err != nil {
		logError("deleteBook", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"book": book,
	})
}
