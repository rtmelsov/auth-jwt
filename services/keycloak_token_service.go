package services

import (
	"encoding/json"
	"go-keycloak-jwt/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetTokenFromKeycloak(username string, password string) (models.TokenData, error) {
	tokenUrl := os.Getenv("TOKEN_URL")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)

	resp, err := http.Post(tokenUrl, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			return
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dat models.TokenData
	if err := json.Unmarshal(body, &dat); err != nil {
		return nil, err
	}

	return dat, nil
}
