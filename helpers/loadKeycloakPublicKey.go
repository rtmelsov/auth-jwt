package helpers

import (
	"encoding/json"
	"fmt"
	"go-keycloak-jwt/models"
	"log"
	"net/http"

	"os"
)

var keycloakCertURL string

func LoadKeycloakPublicKey() ([]models.KeyData, error) {
	// Fetch the public key from Keycloak
	keycloakCertURL = os.Getenv("KEY_CLOAK_CERT_URL")
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
	var jwks models.JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		log.Fatalf("Error decoding JWKs: %v", err)
	}

	// Map JWKs to your keyData format
	var keyData []models.KeyData
	for _, jwk := range jwks.Keys {
		keyData = append(keyData, models.KeyData{
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
