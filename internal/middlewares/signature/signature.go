// Package signature предназначен для проверки целостности запроса.
package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"

	"metrics/internal/logger"
)

type hashResponseWriter struct {
	http.ResponseWriter
	key string
}

// AddSignatureMiddleware добавляет ключ проверки целостности.
func AddSignatureMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if key == "" {
				next.ServeHTTP(res, req)
				return
			}
			responseWriter := &hashResponseWriter{ResponseWriter: res, key: key}

			next.ServeHTTP(responseWriter, req)
		})
	}
}

// Write записывает в заголовок hash для проверки целостности.
func (hrw *hashResponseWriter) Write(b []byte) (int, error) {
	h := hmac.New(sha256.New, []byte(hrw.key))
	h.Write(b)
	hashReq := h.Sum(nil)
	hrw.Header().Set("HashSHA256", base64.URLEncoding.EncodeToString(hashReq))

	return hrw.ResponseWriter.Write(b)
}

// CheckSignaturMiddleware проверяет целостность пакета.
func CheckSignaturMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if key == "" {
				next.ServeHTTP(res, req)
				return
			}

			hashReq := req.Header.Get("HashSHA256")
			if hashReq == "" {
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
			req.Body.Close()

			hashReqByte, err := base64.URLEncoding.DecodeString(hashReq)
			if err != nil {
				logger.Log.Warn("Ошибка декодирования hash")
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			h := hmac.New(sha256.New, []byte(key))
			h.Write(body)
			hashReqCalc := h.Sum(nil)
			if !bytes.Equal(hashReqByte, hashReqCalc) {
				logger.Log.Warn("Не пройдена проверка целостности данных")
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}
