package utils

import (
	"fmt"
	"net"
	"net/http"
)

func SecondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
  str := fmt.Sprintf("%02d:%02d\n", minutes, seconds)
	return str
}

// https://blog.golang.org/context/userip/userip.go
func GetIP(req *http.Request) string {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Printf("userip: %q is not IP:port", req.RemoteAddr)
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		fmt.Printf("userip: %q is not IP:port", req.RemoteAddr)
	}

	return fmt.Sprint(userIP)
}
