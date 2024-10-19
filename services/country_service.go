package services

import (
	"go-keycloak-jwt/models"
	"go-keycloak-jwt/repositories"
)

func GetCountries() (models.Countries, error) {
	return repositories.GetAllCountries()
}

func GetCountryById(reqId string) (models.Country, error) {
	return repositories.GetCountryById(reqId)
}
