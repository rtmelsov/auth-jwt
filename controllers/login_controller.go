package controllers

import (
	"github.com/gin-gonic/gin"
	"go-keycloak-jwt/models"
	"go-keycloak-jwt/services"
	"net/http"
)

// @Login
// @Description Логин на сайт
// @Tags main
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login credentials"
// @Success 200 {object}  models.TokenData
// @Failure 404 {object} map[string]string
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var login models.LoginRequest

	// Привязка JSON-данных к структуре
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := services.GetToken(login.Username, login.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to login",
			"details": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
