package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/bookCRUD/model"
)

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
		logError("signup", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//have push down our parsed data to service level
	err := h.services.CreateUser(c, input)
	if err != nil {
		logError("signup", err)
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
	var input model.LoginInput

	if err := c.BindJSON(&input); err != nil {
		logError("login", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//have push down our parsed data to service level
	session_id, err := h.services.SignIn(c, input)
	if err != nil {
		logError("login", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := json.Marshal(map[string]string{
		"session_id": session_id,
	})
	if err != nil {
		logError("login", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}
