package errs

const (
	ErrSystemError                 = 400001 // 系统错误
	ErrUserInputError              = 400002 // 用户输入错误
	ErrInvalidParams               = 400003 // 请求参数错误
	ErrEmailAlreadyRegistered      = 400004 // 邮箱已经被注册
	ErrUserNotRegistered           = 400005 // 用户没有被注册
	ErrUsernameOrPasswordError     = 400006 // 用户名或密码错误
	ErrCodeInputFrequently         = 400007 // 验证码输入频繁
	ErrLimitHttp                   = 400008 // 触发限流
	ErrCodeInputTooManyTimes       = 400009 // 验证码输入错误次数太多
	ErrCodeInputMany               = 400010 // 验证码输入次数太多
	ErrUseOldPassword              = 400011 // 修改的密码是历史密码
	ErrOldPasswordEqualNewPassword = 400012 // 修改的新密码是旧密码
	ErrPhoneIsBlack                = 400013 // 黑名单手机号访问系统
)
