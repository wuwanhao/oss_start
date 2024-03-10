// 状态码相关
package result

// Codes 定义的状态
type Codes struct {
	SUCCESS              uint
	FAILED               uint
	Message              map[uint]string
	ERROR_GET_METADATA   uint
	ERROR_OBJECT_DELETED uint
	ERROR_UPLOAD_FILE    uint
	ERROR_GET_FILE       uint
}

// ApiCode 状态码
var ApiCode = &Codes{
	SUCCESS:              200,
	FAILED:               501,
	ERROR_GET_METADATA:   50101,
	ERROR_OBJECT_DELETED: 50102,
	ERROR_UPLOAD_FILE:    50103,
	ERROR_GET_FILE:       50104,
}

// 状态信息
func init() {
	ApiCode.Message = map[uint]string{
		ApiCode.SUCCESS:              "成功",
		ApiCode.FAILED:               "失败",
		ApiCode.ERROR_GET_METADATA:   "获取对象元数据失败",
		ApiCode.ERROR_OBJECT_DELETED: "对象不存在",
		ApiCode.ERROR_UPLOAD_FILE:    "文件上传失败",
		ApiCode.ERROR_GET_FILE:       "文件获取失败",
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
