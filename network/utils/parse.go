package utils

import (
	"encoding/hex"
	"net"
	"strconv"
	"strings"
)

type _IP string

func (_ip _IP) Parse() []byte {
	var ip []byte
	tokens := strings.Split(string(_ip), ".")
	for _, t := range tokens {
		digit, _ := strconv.Atoi(t)
		ip = append(ip, byte(digit))
	}

	return ip
}

func ParseID(id string) []byte {
	if len(id) == 0 {
		return nil
	}

	d, _ := hex.DecodeString(id)
	if len(d) < 20 {
		return append(make([]byte, 20-len(d)), d...)
	}
	return d
}

func ParseIP(ip string) *net.TCPAddr {
	if len(ip) == 0 {
		return nil
	}

	tokens := strings.Split(ip, ":")
	port, _ := strconv.Atoi(tokens[1])

	return &net.TCPAddr{
		IP:   _IP(tokens[0]).Parse(),
		Port: port,
	}
}
