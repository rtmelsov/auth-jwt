package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-keycloak-jwt/controllers"
	"go-keycloak-jwt/db"
	_ "go-keycloak-jwt/docs"
	"go-keycloak-jwt/middlewares"
	"log"
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

	// Настраиваем Gin
	r := gin.Default()

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
