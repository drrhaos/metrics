package main

import (
	"net/http"
	"testing"
)

func Test_updateMetricHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateMetricHandler(tt.args.res, tt.args.req)
		})
	}
}
