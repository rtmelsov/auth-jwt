package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-keycloak-jwt/helpers"
	"log"
	"net/http"
	"os"
)

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var keycloakCertURL string

func loadKeycloakPublicKey() ([]helpers.KeyData, error) {
	// Fetch the public key from Keycloak
	resp, err := http.Get(keycloakCertURL)

	if err != nil {
		return nil, fmt.Errorf("failed to get Keycloak cert: %v", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
			return
		}
	}()

	if err != nil {
		return nil, fmt.Errorf("failed to read Keycloak cert: %v", err)
	}

	// Parse the public key
	var jwks JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		log.Fatalf("Error decoding JWKs: %v", err)
	}

	// Map JWKs to your keyData format
	var keyData []helpers.KeyData
	for _, jwk := range jwks.Keys {
		keyData = append(keyData, helpers.KeyData{
			Key:       jwk.Kid,
			Algorithm: jwk.Kty,
			N:         jwk.N,
			E:         jwk.E,
			Sig:       jwk.E, // Assuming "sig" in your format is equivalent to "e"
		})
	}

	if len(keyData) == 0 {
		return nil, fmt.Errorf("no keys found in Keycloak cert")
	}
	return keyData, nil
}

func jwtMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString != "" && len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:] // Убираем "Bearer "
	} else {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}
	jwts, err := loadKeycloakPublicKey()
	_, err = helpers.TokenLoader(tokenString, jwts)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}

	c.Next()
}

func protectedHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Access to protected resource granted"})
}

func loginHandler(c *gin.Context) {
	token, err := helpers.GetTokenFromKeycloak(c.PostForm("username"), c.PostForm("password"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to login",
			"details": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	keycloakCertURL = os.Getenv("KEY_CLOAK_CERT_URL")

	// Настраиваем Gin
	r := gin.Default()

	r.GET("/login", loginHandler)

	// Защищённый маршрут
	r.GET("/protected-route", jwtMiddleware, protectedHandler)

	fmt.Print("Server listening on port 8082")

	// Запуск сервера
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}

}
