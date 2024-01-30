package main

import (
	"bytes"
	"compress/gzip"

	"github.com/drrhaos/metrics/internal/logger"
	"go.uber.org/zap"
)

func compressReqData(reqData []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	zipF := gzip.NewWriter(&buf)
	_, err := zipF.Write(reqData)
	if err != nil {
		logger.Log.Warn("Ошибка записи данных", zap.Error(err))
		return nil, err
	}
	err = zipF.Close()
	if err != nil {
		logger.Log.Warn("Ошибка сжатия данных", zap.Error(err))
		return nil, err
	}
	return &buf, nil
}
