package models

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

type TokenData map[string]interface{}

// Структура для получения данных из POST-запроса
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
