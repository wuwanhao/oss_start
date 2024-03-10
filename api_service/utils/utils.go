package utils

import (
	"net/http"
	"strconv"
)

// 从请求头中获取文件大小
func GetSizeFromHeader(h http.Header) int64 {
	size, _ := strconv.ParseInt(h.Get("Content-Length"), 0, 64)
	return size
}

// 从请求头中获取文件hash
func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}
