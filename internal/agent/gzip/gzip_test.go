package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressReqData(t *testing.T) {
	tests := []struct {
		name    string
		reqData []byte
	}{
		{
			name:    "positive test",
			reqData: []byte("test"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compData, _ := CompressReqData(tt.reqData)

			var r io.Reader
			r, _ = gzip.NewReader(compData)

			var resB bytes.Buffer
			resB.ReadFrom(r)

			result := resB.Bytes()

			assert.Equal(t, tt.reqData, result)
		})
	}
}
