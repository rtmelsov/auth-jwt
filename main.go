package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-keycloak-jwt/controllers"
	"go-keycloak-jwt/db"
	_ "go-keycloak-jwt/docs"
	"go-keycloak-jwt/middlewares"
	"log"
	"time"
)

// @title FCB
// @version 1.0
// @description JWT сервисы для обработки данных с авторизацией

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.ConnectDB()
	defer db.CloseDB()

	// Custom CORS configuration
	config := cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},                   // Allow requests from localhost:3000
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow necessary HTTP methods
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"sec-ch-ua",
			"sec-ch-ua-mobile",
			"sec-ch-ua-platform",
			"User-Agent",
			"Referer"}, // Allow necessary headers, including sec-ch-ua and user-agent
		ExposeHeaders:    []string{"Content-Length", "Authorization"}, // Expose headers if needed
		AllowCredentials: true,                                        // Allow cookies or authentication data
		MaxAge:           12 * time.Hour,                              // Cache preflight for 12 hours
	}

	// Initialize Gin Router with custom CORS configuration
	r := gin.Default()
	r.Use(cors.New(config)) // Apply the custom CORS configuration
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/login", controllers.LoginHandler)

	// Защищённый маршрут

	// Запрос структуры score-карты
	r.POST("/get-score-cards", middlewares.JwtMiddleware, controllers.GetScoreCards)
	r.POST("/score", middlewares.JwtMiddleware, controllers.PostScore)

	// Запрос по странам score-карты
	r.GET("/countries", middlewares.JwtMiddleware, controllers.GetCountries)
	r.GET("/countries/:id", middlewares.JwtMiddleware, controllers.GetCountryById)

	fmt.Print("Server listening on port 8082")

	// Запуск сервера
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}

}
