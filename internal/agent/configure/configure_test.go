// Модуль configure предназначен для настройки программы.

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
	fileTmp := "/tmp/agent.json"
	defer os.Remove(fileTmp)

	os.Args = append(os.Args, "--a=127.0.0.1:9080")
	os.Args = append(os.Args, "--p=33")
	os.Args = append(os.Args, fmt.Sprintf("-c=%s", fileTmp))

	type want struct {
		cfg Config
		ok  bool
	}
	tests := []struct {
		name    string
		cfgJSON Config
		want    want
		isEnv   bool
	}{
		{
			name: "positive test #1",
			cfgJSON: Config{
				Address:        "127.0.0.1:8888",
				ReportInterval: 300,
				PollInterval:   300,
				RateLimit:      3,
				Key:            "test",
			},
			isEnv: false,
			want: want{
				ok: true,
				cfg: Config{
					Address:        "127.0.0.1:9080",
					ReportInterval: 10,
					PollInterval:   33,
					RateLimit:      1,
					Key:            "test",
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

			var cfg Config
			ok := cfg.ReadConfig()
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.cfg.Address, cfg.Address)
			assert.Equal(t, tt.want.cfg.PollInterval, cfg.PollInterval)
		})
	}
}
