// Модуль configure предназначен для настройки программы.

package configure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ReadConfig(t *testing.T) {
	os.Args = append(os.Args, "--a=127.0.0.1:9080")

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
					Address:         "127.0.0.1:9080",
					StoreInterval:   300,
					FileStoragePath: "/tmp/metrics-db.json",
					Restore:         true,
					DatabaseDsn:     "",
					Key:             "",
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
