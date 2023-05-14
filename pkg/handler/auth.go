package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/qazaqpyn/bookCRUD/pkg/logging"
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
		logging.LogError("signup", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//have push down our parsed data to service level
	err := h.services.CreateUser(c, input)
	if err != nil {
		logging.LogError("signup", err)
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
		logging.LogError("login", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//have push down our parsed data to service level
	accessToken, refreshToken, err := h.services.SignIn(c, input)
	if err != nil {
		logging.LogError("login", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": accessToken,
	})
	if err != nil {
		logging.LogError("login", err)
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Header("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		logging.LogError("refresh", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.RefreshTokens(c, cookie)
	if err != nil {
		logging.LogError("refresh", err)
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.Header("Set-Cokie", fmt.Sprintf("refresh-token=%s; HttpOnly", refreshToken))
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}
