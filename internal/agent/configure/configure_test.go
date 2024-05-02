// Модуль configure предназначен для настройки программы.

package configure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_ReadConfig(t *testing.T) {
	t.Setenv("ADDRESS", "127.0.0.1:9090")
	t.Setenv("REPORT_INTERVAL", "11")
	t.Setenv("POLL_INTERVAL", "5")
	t.Setenv("RATE_LIMIT", "2")
	t.Setenv("KEY", "test")

	type want struct {
		ok  bool
		cfg Config
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
					Address:        "127.0.0.1:9090",
					ReportInterval: 11,
					PollInterval:   5,
					RateLimit:      2,
					Key:            "test",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := tt.cfg.ReadConfig()
			assert.Equal(t, tt.want.ok, ok)
		})
	}
}
