package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/bookCRUD/model"
)

type signInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary			SignUp
// @Description		Sign up with user details
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			input body model.User true "account info"
// @Success			200	{string}	string "user created"
// @Failure			404	{object}	errorResponse
// @Failure			404	{object}	errorResponse
// @Failure			500	{object}	errorResponse
// @Router			/auth/signup [post]
func (h *Handler) signup(c *gin.Context) {
	var input model.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//have push down our parsed data to service level
	err := h.services.CreateUser(c, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"msg": "user created",
	})
}

// @Summary Login
// @Description	Login with user credentials
// @Tags Auth
// @Accept json
// @Produce	json
// @Param input body signInput true "account credentials"
// @Success	200	{string} string "token"
// @Failure	400 {object} errorResponse
// @Failure	404	{object}	errorResponse
// @Failure	500	{object} errorResponse
// @Router /auth/login [post]
func (h *Handler) login(c *gin.Context) {
	var input signInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//have push down our parsed data to service level
	token, err := h.services.GenerateToken(c, input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
