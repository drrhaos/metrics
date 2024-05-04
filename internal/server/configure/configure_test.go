// Модуль configure предназначен для настройки программы.

package configure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ReadConfig(t *testing.T) {
	t.Setenv("ADDRESS", "127.0.0.1:9090")
	t.Setenv("STORE_INTERVAL", "11")
	t.Setenv("FILE_STORAGE_PATH", "/home/user/")
	t.Setenv("RESTORE", "true")
	t.Setenv("DATABASE_DSN", "postgres://test:test@test:4321/test?sslmode=disable")
	t.Setenv("KEY", "test")

	type want struct {
		cfg Config
		ok  bool
	}
	tests := []struct {
		name      string
		cfg       Config
		envKey    string
		fieldName string
		envVal    string
		want      want
	}{
		{
			name: "positive test",
			want: want{
				ok: true,
				cfg: Config{
					Address:         "127.0.0.1:9090",
					StoreInterval:   11,
					FileStoragePath: "/home/user/",
					Restore:         true,
					DatabaseDsn:     "postgres://test:test@test:4321/test?sslmode=disable",
					Key:             "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := tt.cfg.ReadConfig()
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.cfg, tt.cfg)
		})
	}
}
