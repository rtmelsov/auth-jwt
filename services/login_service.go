package services

import (
	"go-keycloak-jwt/models"
)

func GetToken(username string, password string) (models.TokenData, error) {
	return GetTokenFromKeycloak(username, password)
}
