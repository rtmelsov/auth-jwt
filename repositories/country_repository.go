package repositories

import (
	"context"
	"go-keycloak-jwt/db"
	"go-keycloak-jwt/models"
)

func GetAllCountries() (models.Countries, error) {
	rows, err := db.DB.Query(context.Background(), "SELECT id, name, code FROM countries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var countries models.Countries
	for rows.Next() {
		var id int
		var name, code string
		if err := rows.Scan(&id, &name, &code); err != nil {
			return nil, err
		}
		country := map[string]interface{}{
			"id":   id,
			"name": name,
			"code": code,
		}
		countries = append(countries, country)
	}
	return countries, nil
}

func GetCountryById(reqId string) (models.Country, error) {
	var id int
	var name, code string
	err := db.DB.QueryRow(context.Background(), "SELECT id, name, code FROM countries WHERE id=$1", reqId).Scan(&id, &name, &code)
	if err != nil {
		return nil, err
	}
	return models.Country{
		"id":   id, // Теперь id имеет тип int
		"name": name,
		"code": code,
	}, nil
}
