// Модуль cryptodata предназначен шифрования данных.

package cryptodata

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestDecryptMiddleware(t *testing.T) {
	privateKeyGood, _ := rsa.GenerateKey(rand.Reader, 2048)
	privateKeyBad, _ := rsa.GenerateKey(rand.Reader, 2048)

	r := chi.NewRouter()
	r.Use(DecryptMiddleware(privateKeyGood))
	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
		w.WriteHeader(http.StatusOK)
	})

	type want struct {
		code int
		body string
	}
	tests := []struct {
		name       string
		privateKey *rsa.PrivateKey
		want       want
	}{
		{
			name:       "positive positive check Decrypt #1",
			privateKey: privateKeyGood,
			want: want{
				body: "test",
				code: 200,
			},
		},
		{
			name:       "negative positive check Decrypt #2",
			privateKey: privateKeyBad,
			want: want{
				code: 400,
			},
		},
		{
			name:       "negative positive check Decrypt #3",
			privateKey: privateKeyBad,
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ciphertext, _ := rsa.EncryptPKCS1v15(rand.Reader, &test.privateKey.PublicKey, []byte("test"))

			req := httptest.NewRequest(http.MethodPost, "/update", bytes.NewReader(ciphertext))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, test.want.body, w.Body.String())
			assert.Equal(t, test.want.code, w.Code)
		})
	}
}

func TestEncrypt(t *testing.T) {
	privateKeyGood, _ := rsa.GenerateKey(rand.Reader, 2048)
	type args struct {
		plaintext []byte
		publicKey *rsa.PublicKey
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positive positive check Decrypt #1",
			args: args{
				plaintext: []byte("test"),
				publicKey: &privateKeyGood.PublicKey,
			},
			wantErr: false,
			want:    "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt([]byte(tt.args.plaintext), tt.args.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			text, _ := Decrypt(got, privateKeyGood)
			if !reflect.DeepEqual(string(text), tt.want) {
				t.Errorf("Encrypt() = %v, want %v", text, tt.want)
			}
		})
	}
}
