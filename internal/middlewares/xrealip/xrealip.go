package xrealip

import (
	"fmt"
	"net"
	"net/http"

	"metrics/internal/logger"
)

var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func RealIP(trustedSubnet string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if trustedSubnet == "" {
				next.ServeHTTP(res, req)
				return
			}

			ip := getIP(req)

			_, ipNet, err := net.ParseCIDR(trustedSubnet)
			if err != nil {
				logger.Log.Warn("Не верный фомат подсети")
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			if ipNet.Contains(ip) {
				logger.Log.Info(fmt.Sprintf("%s входит в подсеть %s", ip, ipNet))
			} else {
				logger.Log.Warn(fmt.Sprintf("%s не входит в подсеть %s", ip, ipNet))
				res.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(res, req)
		})
	}
}

func getIP(r *http.Request) net.IP {
	var ip string

	if tmpIP := r.Header.Get(xRealIP); tmpIP != "" {
		ip = tmpIP
	}

	return net.ParseIP(ip)
}
