package configure

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"metrics/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ReadConfig(t *testing.T) {
	fileTmp := "/tmp/server.json"
	defer os.Remove(fileTmp)

	os.Args = append(os.Args, "--a=127.0.0.1:9080")
	os.Args = append(os.Args, fmt.Sprintf("-c=%s", fileTmp))

	type want struct {
		cfg Config
		ok  bool
	}
	tests := []struct {
		name      string
		cfgJSON   Config
		cfg       Config
		envKey    string
		fieldName string
		envVal    string
		want      want
	}{
		{
			name: "positive test",
			cfgJSON: Config{
				Address:       "127.0.0.1:8888",
				StoreInterval: 300,
				Restore:       true,
			},
			want: want{
				ok: true,
				cfg: Config{
					Address:         "127.0.0.1:9080",
					StoreInterval:   300,
					FileStoragePath: "/tmp/metrics-db.json",
					Restore:         true,
					DatabaseDsn:     "",
					Key:             "",
					ConfigPath:      "/tmp/server.json",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.cfgJSON)
			if err != nil {
				logger.Log.Warn("не удалось преобразовать структуру")
				return
			}
			file, err := os.OpenFile(fileTmp, os.O_WRONLY|os.O_CREATE, 0o666)
			if err != nil {
				return
			}

			_, err = file.Write(data)
			if err != nil {
				return
			}

			ok := tt.cfg.ReadConfig()
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.cfg, tt.cfg)
		})
	}
}
