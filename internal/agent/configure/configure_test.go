// Модуль configure предназначен для настройки программы.

package configure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ReadConfig(t *testing.T) {
	os.Args = append(os.Args, "--a=127.0.0.1:9080")
	os.Args = append(os.Args, "--p=33")

	type want struct {
		cfg Config
		ok  bool
	}
	tests := []struct {
		name  string
		isEnv bool
		want  want
	}{
		{
			name:  "positive test #1",
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
			var cfg Config
			ok := cfg.ReadConfig()
			assert.Equal(t, tt.want.ok, ok)
			assert.Equal(t, tt.want.cfg.Address, cfg.Address)
			assert.Equal(t, tt.want.cfg.PollInterval, cfg.PollInterval)
		})
	}
}
