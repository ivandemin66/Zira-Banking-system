package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash проверяет соответствие пароля хешу
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// HashCVV хеширует CVV с использованием bcrypt
func HashCVV(cvv string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckCVVHash проверяет соответствие CVV хешу
func CheckCVVHash(cvv, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cvv))
	return err == nil
}

// ComputeHMAC вычисляет HMAC для данных
func ComputeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyHMAC проверяет HMAC для данных
func VerifyHMAC(data, messageMAC string, secret []byte) bool {
	expectedMAC := ComputeHMAC(data, secret)
	return hmac.Equal([]byte(messageMAC), []byte(expectedMAC))
}
