package middlewares

import (
	"github.com/gin-gonic/gin"
	"go-keycloak-jwt/helpers"
	"net/http"
)

func JwtMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString != "" && len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:] // Убираем "Bearer "
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}
	jwts, err := helpers.LoadKeycloakPublicKey()
	userID, userName, err := helpers.TokenLoader(tokenString, jwts)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}

	// Сохраняем данные в контексте
	c.Set("userID", userID)
	c.Set("userName", userName)

	c.Next()
}
