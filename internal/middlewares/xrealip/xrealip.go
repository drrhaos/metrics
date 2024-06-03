// Package xrealip предназначен для проверки IP клиента.
package xrealip

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"metrics/internal/logger"
)

var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

// RealIP проверяет входи ли адрес клиета в список разрешенных подсетей
func RealIP(trustedSubnet string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if trustedSubnet == "" {
				next.ServeHTTP(res, req)
				return
			}

			ips := getIP(req)

			_, ipNet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				logger.Log.Warn("Не верный фомат подсети")
				res.WriteHeader(http.StatusBadRequest)
				return
			}
			inContain := false
			for _, ip := range strings.Split(ips, ", ") {
				if ipNet.Contains(net.ParseIP(ip)) {
					logger.Log.Info(fmt.Sprintf("%s входит в подсеть %s", ip, ipNet))
					inContain = true
				} else {
					logger.Log.Warn(fmt.Sprintf("%s не входит в подсеть %s", ip, ipNet))
				}
			}

			if !inContain {
				res.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}

func getIP(r *http.Request) string {
	var ipAddresses string

	if tmpIP := r.Header.Get(xRealIP); tmpIP != "" {
		ipAddresses = tmpIP
	}

	return ipAddresses
}
