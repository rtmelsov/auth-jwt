package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"math/big"
)

// Структура для хранения ключей из JSON
type KeyData struct {
	Key              string `json:"key"`
	Algorithm        string `json:"algorithm"`
	EncryptionScheme string `json:"encryption_scheme,omitempty"`
	SignatureScheme  string `json:"signature_scheme,omitempty"`
	N                string `json:"n"`
	E                string `json:"e"`
	Sig              string `json:"sig"`
}

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

func tokenLoader(tokenString string, keys []KeyData) (*jwt.Token, error) {
	// Парсинг массива ключей
	// JWT токен для валидации

	// Парсинг и валидация токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
				// Преобразуем в PEM и возвращаем ключ
				pemKey, err := convertRSAPublicKeyToPEM(pubKey)
				if err != nil {
					return nil, err
				}
				fmt.Printf("PEM ключ:\n%s\n", pemKey)
				return pubKey, nil
			}
		}

		return nil, fmt.Errorf("ключ с 'kid' %s не найден %v", kid, keys)
	})

	if err != nil {
		return nil, err
	}

	// Проверка, является ли токен действительным
	if token.Valid {
		return token, err
	} else {
		return nil, fmt.Errorf("Ошибка при валидности токена")
	}
}
