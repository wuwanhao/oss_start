// 状态码相关
package result

// Codes 定义的状态
type Codes struct {
	SUCCESS                                 uint
	FAILED                                  uint
	Message                                 map[uint]string
	NOAUTH                                  uint
	AUTHFORMATERROR                         uint
	INVALIDTOKEN                            uint
	MissingLoginParameter                   uint
	VerificationCodeHasExpired              uint
	CAPTCHANOTTRUE                          uint
	PASSWORDNOTTRUE                         uint
	STATUSISENABLE                          uint
	ROLENAMEALREADYEXISTS                   uint
	MENUISEXIST                             uint
	DELSYSMENUFAILED                        uint
	DEPTISEXIST                             uint
	DEPTISDISTRIBUTE                        uint
	POSTALREADYEXISTS                       uint
	MissingNewAdminParameter                uint
	USERNAMEALREADYEXISTS                   uint
	MissingModificationOfPersonalParameters uint
	MissingChangePasswordParameter          uint
	RESETPASSWORD                           uint
	FILEUPLOADERROR                         uint
}

// ApiCode 状态码
var ApiCode = &Codes{
	SUCCESS:                                 200,
	FAILED:                                  501,
	NOAUTH:                                  403,
	AUTHFORMATERROR:                         405,
	INVALIDTOKEN:                            406,
	MissingLoginParameter:                   407,
	VerificationCodeHasExpired:              408,
	CAPTCHANOTTRUE:                          409,
	PASSWORDNOTTRUE:                         410,
	STATUSISENABLE:                          411,
	ROLENAMEALREADYEXISTS:                   412,
	MENUISEXIST:                             413,
	DELSYSMENUFAILED:                        414,
	DEPTISEXIST:                             415,
	DEPTISDISTRIBUTE:                        416,
	POSTALREADYEXISTS:                       417,
	MissingNewAdminParameter:                418,
	USERNAMEALREADYEXISTS:                   419,
	MissingModificationOfPersonalParameters: 420,
	MissingChangePasswordParameter:          421,
	RESETPASSWORD:                           422,
	FILEUPLOADERROR:                         427,
}

// 状态信息
func init() {
	ApiCode.Message = map[uint]string{
		ApiCode.SUCCESS:                                 "成功",
		ApiCode.FAILED:                                  "失败",
		ApiCode.NOAUTH:                                  "请求头中auth为空",
		ApiCode.AUTHFORMATERROR:                         "请求头中auth格式有误",
		ApiCode.INVALIDTOKEN:                            "token无效或者登录过期，请重新登录",
		ApiCode.MissingLoginParameter:                   "缺少登录参数",
		ApiCode.VerificationCodeHasExpired:              "验证码过期",
		ApiCode.CAPTCHANOTTRUE:                          "验证码错误，请重新输入",
		ApiCode.PASSWORDNOTTRUE:                         "密码错误",
		ApiCode.STATUSISENABLE:                          "您的账号已被停用，请联系管理员",
		ApiCode.ROLENAMEALREADYEXISTS:                   "角色名称或者权限字符已经存在，请重新输入",
		ApiCode.MENUISEXIST:                             "菜单名称已经存在，请重新输入",
		ApiCode.DELSYSMENUFAILED:                        "菜单已分配，不能删除",
		ApiCode.DEPTISEXIST:                             "部门名称已存在，请重新输入",
		ApiCode.DEPTISDISTRIBUTE:                        "部门已分配，不能删除",
		ApiCode.POSTALREADYEXISTS:                       "岗位名称或岗位编码已存在，请重新输入",
		ApiCode.MissingNewAdminParameter:                "缺少新增管理员参数",
		ApiCode.USERNAMEALREADYEXISTS:                   "用户名已存在",
		ApiCode.MissingModificationOfPersonalParameters: "缺少修改个人信息的参数",
		ApiCode.MissingChangePasswordParameter:          "缺少修改密码参数",
		ApiCode.RESETPASSWORD:                           "两次输入的密码不一致，请重新输入",
		ApiCode.FILEUPLOADERROR:                         "文件上传失败",
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
