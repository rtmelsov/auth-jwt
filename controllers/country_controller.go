package controllers

import (
	"github.com/gin-gonic/gin"
	"go-keycloak-jwt/services"
	"net/http"
)

// @Summary Get all countries
// @Description Возвращает список всех стран (защищенный маршрут)
// @Tags countries
// @Produce json
// @Success 200 {array} []models.Country
// @Security BearerAuth
// @Router /countries [get]
func GetCountries(c *gin.Context) {
	data, err := services.GetCountries()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"countries": data})
}

// @Summary Get country by ID
// @Description Возвращает информацию о стране по её идентификатору
// @Tags countries
// @Produce json
// @Param id path int true "Country ID"
// @Success 200 {object} models.Country
// @Failure 404 {object} map[string]string "Country not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /countries/{id} [get]
func GetCountryById(c *gin.Context) {
	reqId := c.Param("id")
	data, err := services.GetCountryById(reqId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": data})
}
