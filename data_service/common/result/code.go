// 状态码相关
package result

// Codes 定义的状态
type Codes struct {
	// status
	SUCCESS uint
	FAILED  uint
	Message map[uint]string
	// error code
	FILE_NOT_FOUND        uint
	FILE_BYTE_TRANS_ERROR uint
	FILE_NAME_CHECK_ERROR uint
	FILE_PUT_EMPTY_ERROR  uint
	FILE_UPLOAD_ERROR     uint
}

// ApiCode 状态码
var ApiCode = &Codes{
	SUCCESS:               200,
	FAILED:                501,
	FILE_NOT_FOUND:        50401,
	FILE_BYTE_TRANS_ERROR: 50402,
	FILE_PUT_EMPTY_ERROR:  50403,
	FILE_UPLOAD_ERROR:     50404,
	FILE_NAME_CHECK_ERROR: 50405,
}

// 状态信息
func init() {
	ApiCode.Message = map[uint]string{
		ApiCode.SUCCESS:               "成功",
		ApiCode.FAILED:                "失败",
		ApiCode.FILE_NOT_FOUND:        "该文件不存在",
		ApiCode.FILE_BYTE_TRANS_ERROR: "文件返回数据流失败",
		ApiCode.FILE_PUT_EMPTY_ERROR:  "上传文件为空",
		ApiCode.FILE_NAME_CHECK_ERROR: "上传文件名不能为空",
		ApiCode.FILE_UPLOAD_ERROR:     "文件上传失败",
	}
}

// GetMessage 供外部调用
func (c *Codes) GetMessage(code uint) string {
	message, ok := c.Message[code]
	if !ok {
		return ""
	}
	return message
}
