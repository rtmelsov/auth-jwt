package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-keycloak-jwt/models"
	"math/big"
)

// Преобразуем модуль и экспоненту в PEM-формат публичного ключа
func createRSAPublicKeyFromModExp(nStr string, eStr string) (*rsa.PublicKey, error) {
	// Декодируем модуль (n) из base64URL
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при декодировании модуля: %v", err)
	}

	// Декодируем экспоненту (e) из base64URL
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка при декодировании экспоненты: %v", err)
	}

	// Преобразуем экспоненту в целое число
	e := 0
	for _, b := range eBytes {
		e = e*256 + int(b)
	}

	// Преобразуем модуль в big int
	n := new(big.Int).SetBytes(nBytes)

	// Создаем публичный ключ RSA
	pubKey := &rsa.PublicKey{
		N: n,
		E: e,
	}

	return pubKey, nil
}

// Функция для конвертации RSA ключа в PEM-формат
func convertRSAPublicKeyToPEM(pubKey *rsa.PublicKey) (string, error) {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", fmt.Errorf("ошибка при маршализации публичного ключа: %v", err)
	}

	pemKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	return string(pem.EncodeToMemory(pemKey)), nil
}

func TokenLoader(tokenString string, keys []models.KeyData) (string, string, error) {
	// Парсинг массива ключей
	// JWT токен для валидации

	// Структура для хранения claim'ов (данных) токена
	claims := jwt.MapClaims{}

	// Парсинг и валидация токена
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Убедитесь, что используемый метод подписи — RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Извлекаем "kid" из заголовка токена
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("не найден 'kid' в заголовке токена")
		}

		// Ищем соответствующий ключ по "kid"
		for _, key := range keys {
			if key.Key == kid {
				// Преобразуем модуль и экспоненту в публичный ключ
				pubKey, err := createRSAPublicKeyFromModExp(key.N, key.E)
				if err != nil {
					return nil, fmt.Errorf("ошибка при создании публичного ключа: %v", err)
				}
				return pubKey, nil
			}
		}

		return nil, fmt.Errorf("ключ с 'kid' %s не найден", kid)
	})

	if err != nil {
		return "", "", err
	}

	// Проверка валидности токена
	if token.Valid {
		// Извлекаем данные из токена
		userID, ok := claims["sub"].(string)
		if !ok {
			return "", "", fmt.Errorf("не найден userID в токене")
		}

		userName, ok := claims["preferred_username"].(string)
		if !ok {
			return "", "", fmt.Errorf("не найден userName в токене")
		}

		return userID, userName, nil
	} else {
		return "", "", fmt.Errorf("токен недействителен")
	}
}
