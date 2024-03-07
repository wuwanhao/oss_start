// 通用返回结构
package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Result 通用结构
type Result struct {
	Code    int `json:"code"` // 状态码
	Message string `json:message` // 提示信息
	Data    interface{} `json:"data"` // 返回的数据
}


// Success 返回成功消息体
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := Result{}
	res.Code = int(ApiCode.SUCCESS)
	res.Message = ApiCode.GetMessage(ApiCode.SUCCESS)
	res.Data = data
	c.JSON(http.StatusOK, res)
}
// Failed 返回失败消息体
func Failed(c *gin.Context, code int, message string) {
	res := Result{}
	res.Code = code
	res.Message = message
	res.Data = gin.H{}
	c.JSON(http.StatusOK, res)
}