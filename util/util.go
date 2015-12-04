package util

import (
	"net"
	"net/http"
	"strconv"
	"strings"
)

var localIp string

func LocalIP() string {

	if localIp != "" {
		return localIp
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ip := strings.Split(addr.String(), "/")[0]
		if ip == "127.0.0.1" || ip == "::1" {
			continue
		}
		first := strings.Split(ip, ".")[0]
		if _, err := strconv.Atoi(first); err == nil {
			localIp = ip
			return ip
		}
	}
	return ""
}

func HttpClietIP(req *http.Request) string {
	ip := ""
	if ips := req.Header.Get("X-Forwarded-For"); ips != "" {
		ip = strings.Split(ips, ",")[0]
	}

	if ip != "" {
		rip := strings.Split(ip, ":")
		return rip[0]
	}
	ips := strings.Split(req.RemoteAddr, ":")
	if len(ips) > 0 {
		if ips[0] != "[" {
			return ips[0]
		}
	}
	return "127.0.0.1"
}
