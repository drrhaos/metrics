// Модуль cryptodata предназначен шифрования данных.

package cryptodata

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"net/http"

	"metrics/internal/logger"
)

// Encrypt шифрование сообщения открытым ключом
func Encrypt(plaintext []byte, publicKey any) ([]byte, error) {
	ciphertext := make([]byte, 0)
	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		maxSize := publicKey.N.BitLen()/8 - 11 // Максимальный размер данных для шифрования

		for i := 0; i < len(plaintext); i += maxSize {
			end := i + maxSize
			if end > len(plaintext) {
				end = len(plaintext)
			}
			block, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plaintext[i:end])
			if err != nil {
				return nil, err
			}
			ciphertext = append(ciphertext, block...)
		}
	default:
		panic("error")
	}
	return ciphertext, nil
}

// Decrypt дешифрование сообщения закрытым ключом
func Decrypt(ciphertext []byte, privateKey any) ([]byte, error) {
	plaintext := make([]byte, 0)

	switch privateKey := privateKey.(type) {
	case *rsa.PrivateKey:
		maxSize := (privateKey.N.BitLen() + 7) / 8 // Максимальный размер данных для расшифрования

		for i := 0; i < len(ciphertext); i += maxSize {
			end := i + maxSize
			if end > len(ciphertext) {
				end = len(ciphertext)
			}
			block, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext[i:end])
			if err != nil {
				return nil, err
			}
			plaintext = append(plaintext, block...)
		}
	default:
		panic("error")
	}
	return plaintext, nil
}

// DecryptMiddleware обработчик распакавывает тело ответа.
func DecryptMiddleware(key any) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if key == nil {
				next.ServeHTTP(res, req)
				return
			}

			body, err := io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(body))
			if err != nil {
				logger.Log.Warn("Ошибка чтения")
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			cipherTextBytes, err := Decrypt(body, key)
			if err != nil {
				logger.Log.Warn("Ошибка расшифроки данных")
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			req.Body = io.NopCloser(bytes.NewBuffer((cipherTextBytes)))
			req.Body.Close()

			next.ServeHTTP(res, req)
		})
	}
}
