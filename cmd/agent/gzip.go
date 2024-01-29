package main

import (
	"bytes"
	"compress/gzip"
	"log"
)

func compressReqData(reqData []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	zipF := gzip.NewWriter(&buf)
	_, err := zipF.Write(reqData)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = zipF.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &buf, nil
}
