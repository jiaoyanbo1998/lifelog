package errs

import "errors"

var (
	EmailExist         = errors.New("邮箱已存在")
	UserNotExist       = errors.New("用户没有注册")
	PasswordWrong      = errors.New("旧密码错误输入错误")
	UseOldPassword     = errors.New("不能使用历史密码")
	NeedUpdatePassword = "需要修改密码"
)
