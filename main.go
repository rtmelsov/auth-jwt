package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
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

var db *pgx.Conn

func connectDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	db = conn
	fmt.Println("Successfully connected to database")
}

func closeDB() {
	err := db.Close(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

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

func getCountries(c *gin.Context) {
	rows, err := db.Query(context.Background(), "SELECT id, name, code FROM countries")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch countries"})
		return
	}
	defer rows.Close()
	var countries []map[string]interface{}
	for rows.Next() {
		var id int
		var name, code string
		if err := rows.Scan(&id, &name, &code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning countries"})
			return
		}
		country := map[string]interface{}{
			"id":   id,
			"name": name,
			"code": code,
		}
		countries = append(countries, country)
	}
	c.JSON(http.StatusOK, gin.H{"countries": countries})
}

func getCountryById(c *gin.Context) {
	reqId := c.Param("id")
	var id int
	var name, code string
	err := db.QueryRow(context.Background(), "SELECT id, name, code FROM countries WHERE id=$1", reqId).Scan(&id, &name, &code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch countries"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":   id,
		"name": name,
		"code": code,
	})
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

	connectDB()
	defer closeDB()

	// Настраиваем Gin
	r := gin.Default()

	r.GET("/login", loginHandler)

	// Защищённый маршрут
	r.GET("/protected-route", jwtMiddleware, protectedHandler)
	r.GET("/countries", jwtMiddleware, getCountries)
	r.GET("/countries/:id", jwtMiddleware, getCountryById)

	fmt.Print("Server listening on port 8082")

	// Запуск сервера
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}

}
