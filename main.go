package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
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

func loadKeycloakPublicKey() ([]KeyData, error) {
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
	var keyData []KeyData
	for _, jwk := range jwks.Keys {
		keyData = append(keyData, KeyData{
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

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString != "" && len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:] // Убираем "Bearer "
		} else {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}
		jwts, err := loadKeycloakPublicKey()
		_, err = tokenLoader(tokenString, jwts)
		if err != nil {
			fmt.Printf("Failed to load Keycloak public key: %s", err)
			http.Error(w, "Failed to load Keycloak public key: %s", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "Access to protected resource granted")
	if err != nil {
		fmt.Println("protectedHandler", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	keycloakCertURL = os.Getenv("KEY_CLOAK_CERT_URL")

	http.Handle("/protected-route", jwtMiddleware(http.HandlerFunc(protectedHandler)))

	fmt.Print("Server listening on port 8082")

	log.Fatal(http.ListenAndServe(":8082", nil))

}
