package utils

import (
	"encoding/base64"
	"io"
	"crypto/rand"
	"crypto/md5"
	"encoding/hex"
	"net"
)

//生成Guid字串
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return getMd5String(base64.URLEncoding.EncodeToString(b))
}


func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetLocalAddress() string {
	conn, err := net.Dial("udp", "www.google.com.hk:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	return conn.LocalAddr().String()
}
