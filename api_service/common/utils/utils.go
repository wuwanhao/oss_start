package utils

import (
	"github.com/gin-gonic/gin"
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

// todo: 根据文件头部的长度信息计算文件hash
func GetHash(c *gin.Context) {

}
