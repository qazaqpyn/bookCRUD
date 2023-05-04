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

func (h *Handler) logout(c *gin.Context) {

}
