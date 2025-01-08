package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	userv1 "lifelog-grpc/api/proto/gen/user/v1"
	"lifelog-grpc/errs"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/user/vo"
	"net/http"
	"strconv"
	"strings"
)

var ErrVerifyToMany = errors.New("验证码验证太频繁")

// UserHandler 用户处理器
type UserHandler struct {
	userServiceClient userv1.UserServiceClient
	logger            loggerx.Logger
	emailRegexp       *regexp.Regexp // 用于校验邮箱格式的正则表达式
	passwordRegexp    *regexp.Regexp // 用于校验密码的正则表达式
	phoneRegexp       *regexp.Regexp
	// 组合jwtHandler
	jwtHandler *JWTHandler
}

// NewUserHandler 构造函数创建用户处理器
func NewUserHandler(
	userServiceClient userv1.UserServiceClient,
	l loggerx.Logger,
	jwtHandler *JWTHandler) *UserHandler {
	// 定义正则表达式常量
	const (
		EmailRegexp    = "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
		PasswordRegexp = "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[&$#])[A-Za-z\\d&$#]{8,}$"
		PhoneRegexp    = "^1[3-9]\\d{9}$"
	)
	return &UserHandler{
		userServiceClient: userServiceClient,
		logger:            l,
		// 预编译正则表达式（性能优化：只需要编译一次，后续就可以重复使用）
		emailRegexp:    regexp.MustCompile(EmailRegexp, regexp.None),
		passwordRegexp: regexp.MustCompile(PasswordRegexp, regexp.None),
		phoneRegexp:    regexp.MustCompile(PhoneRegexp, regexp.None),
		jwtHandler:     jwtHandler,
	}
}

// RegisterRoutes 注册路由
func (userHandler *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 用户路由组
	rg := server.Group("/user")
	// 邮箱和密码注册
	rg.POST("/register_email_password", userHandler.RegisterByEmailAndPassword)
	// 邮箱和密码登录
	rg.POST("/login_email_password", userHandler.LoginByEmailAndPassword)
	// 获取用户信息
	rg.GET("/info/:id", userHandler.GetUserInfoById)
	// 修改用户信息
	rg.PUT("/info", userHandler.UpdateUserInfoById)
	// 删除用户信息
	rg.DELETE("/info/:ids", userHandler.DeleteUSerInfoByIds)
	// token续约
	rg.GET("/refresh_token", userHandler.RefreshToken)
	// 退出登录
	rg.GET("/logout", userHandler.Logout)
	// 手机号登录
	rg.POST("/login_phone_code", userHandler.LoginByPhoneCode)
}

func (userHandler *UserHandler) RegisterByEmailAndPassword(ctx *gin.Context) {
	// 定义结构体
	type RegisterReq struct {
		// json字段email和结构体字段Email相对应
		Email string `json:"email"`
		// json字段password和结构体字段Password相对应
		Password string `json:"password"`
		// json字段confirm_password和结构体字段confirmPassword相对应
		ConfirmPassword string `json:"confirm_password"`
	}
	// 将http请求携带的参数，绑定到对应的结构体字段上
	var req RegisterReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "请求参数错误",
			Data: "error",
		})
		userHandler.logger.Error("参数bind失败，RegisterByEmailAndPassword方法",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	if req.Email == "" {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "邮箱不能为空",
			Data: "error",
		})
		userHandler.logger.Error("邮箱为空",
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	if req.Password == "" || req.ConfirmPassword == "" {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "密码不能为空",
			Data: "error",
		})
		userHandler.logger.Error("密码为空",
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	// 校验邮箱格式
	ok, err := userHandler.emailRegexp.MatchString(req.Email)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrUserInputError,
			Msg:  "邮箱格式不正确",
			Data: "error",
		})
		userHandler.logger.Error("邮箱格式不正确",
			loggerx.Error(err),
			loggerx.String("email", req.Email),
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("正则表达式校验超时，RegisterByEmailAndPassword方法",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	// 检查两次输入的密码是否相同
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrUserInputError,
			Msg:  "两次输入的密码不一致",
			Data: "error",
		})
		userHandler.logger.Error("两次输入的密码不一致，RegisterByEmailAndPassword方法",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	// 校验密码格式
	ok, err = userHandler.passwordRegexp.MatchString(req.Password)
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrUserInputError,
			Msg:  "密码格式错误",
			Data: "error",
		})
		userHandler.logger.Error("密码格式错误，RegisterByEmailAndPassword方法",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("正则校验超时，RegisterByEmailAndPassword方法",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	// 调用service层
	_, err = userHandler.userServiceClient.RegisterByEmailAndPassword(ctx.Request.Context(),
		&userv1.RegisterByEmailAndPasswordRequest{
			UserDomain: &userv1.UserDomain{
				Email:    req.Email,
				Password: req.Password,
			},
		})
	if err != nil {
		if strings.Contains(err.Error(), errs.EmailExist.Error()) {
			ctx.JSON(http.StatusOK, Result[string]{
				Code: errs.ErrEmailAlreadyRegistered,
				Msg:  "邮箱已存在",
				Data: "error",
			})
			userHandler.logger.Error("邮箱重复，可能有人在搞你的系统",
				loggerx.Error(err),
				loggerx.String("email", req.Email),
				loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
			return
		}
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("注册失败，RegisterByEmailAndPassword方法",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:RegisterByEmailAndPassword"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "注册成功",
		Data: "success",
	})
}

func (userHandler *UserHandler) LoginByEmailAndPassword(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("参数bind失败，LoginByEmailAndPassword方法",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:Login"))
		return
	}
	res, err := userHandler.userServiceClient.LoginByEmailAndPassword(ctx.Request.Context(),
		&userv1.LoginByEmailAndPasswordRequest{
			UserDomain: &userv1.UserDomain{
				Email:    req.Email,
				Password: req.Password,
			},
		})
	if err != nil {
		if strings.Contains(err.Error(), errs.UserNotExist.Error()) {
			ctx.JSON(http.StatusOK, Result[string]{
				Code: errs.ErrUserNotRegistered,
				Msg:  "用户没有被注册",
				Data: "error",
			})
			userHandler.logger.Error("用户不存在，可能有人在搞你的系统",
				loggerx.Error(err),
				loggerx.String("email", req.Email),
				loggerx.String("method", "UserHandler:Login"),
			)
			return
		}
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrUsernameOrPasswordError,
			Msg:  "用户名或密码错误",
			Data: "error",
		})
		userHandler.logger.Error("用户名或密码错误，可能有人在搞你的系统",
			loggerx.Error(err),
			loggerx.String("email", req.Email),
			loggerx.String("method", "UserHandler:Login"))
		return
	}
	// 设置token的负载，也就是存储在token中的数据
	userClaims := UserClaims{
		Id:        res.GetUserDomain().Id,
		NickName:  res.GetUserDomain().NickName,
		Authority: res.GetUserDomain().Authority,
	}
	// 登陆成功，生成长短token
	userHandler.jwtHandler.SetJwt(ctx, userClaims, true)
	// 提醒用户需要修改密码
	if res.Info == errs.NeedUpdatePassword {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 200,
			Msg:  "登录成功",
			Data: "你的密码太久没有修改了，你需要去修改密码",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "登录成功",
		Data: "success",
	})
}

func (userHandler *UserHandler) GetUserInfoById(ctx *gin.Context) {
	// 获取路径参数
	idStr := ctx.Param("id")
	// 将string转换为int64
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("参数转换失败", loggerx.Error(err),
			loggerx.String("method", "UserHandler:GetUserInfoById"))
		return
	}
	// 获取当前用户信息
	userInfo, ok := userHandler.jwtHandler.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "UserHandler:GetUserInfoById"))
		return
	}
	if userInfo.Id != id {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("不能非法查找其他用户信息", loggerx.Error(err),
			loggerx.String("method", "UserHandler:GetUserInfoById"))
		return
	}
	// 调用service层
	user, err := userHandler.userServiceClient.GetUserInfoById(ctx, &userv1.GetUserInfoByIdRequest{
		UserDomain: &userv1.UserDomain{
			Id: id,
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), errs.UserNotExist.Error()) {
			ctx.JSON(http.StatusOK, Result[string]{
				Code: errs.ErrUserNotRegistered,
				Msg:  "用户不存在",
				Data: "error",
			})
			userHandler.logger.Error("用户不存在", loggerx.Error(err),
				loggerx.String("method", "UserHandler:GetUserInfoById"))
			return
		}
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("查询用户信息失败", loggerx.Error(err),
			loggerx.String("method", "UserHandler:GetUserInfoById"))
		return
	}
	ctx.JSON(http.StatusOK, Result[vo.UserVo]{
		Code: 200,
		Msg:  "查询成功",
		Data: vo.UserVo{
			Id:       user.GetUserDomain().Id,
			Email:    user.GetUserDomain().Email,
			Phone:    user.GetUserDomain().Phone,
			NickName: user.GetUserDomain().NickName,
		},
	})
}

func (userHandler *UserHandler) UpdateUserInfoById(ctx *gin.Context) {
	type updateReq struct {
		Id          int64  `json:"id" binding:"required"`
		Password    string `json:"password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
		NickName    string `json:"nick_name" binding:"required"`
		Phone       string `json:"phone" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Authority   int64  `json:"authority" binding:"required"`
	}
	var req updateReq
	// 如果绑定失败或字段为空，返回错误
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "请求参数错误",
			Data: "error",
		})
		userHandler.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method", "UserHandler:UpdateUserInfoById"))
		return
	}
	// 获取用户信息
	info, ok := userHandler.jwtHandler.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:UpdateUserInfoById"))
		return
	}
	if info.Id != req.Id {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("非法用户修改信息",
			loggerx.Error(err),
			loggerx.String("method", "UserHandler:UpdateUserInfoById"))
		return
	}
	// 调用service层
	_, err = userHandler.userServiceClient.UpdateUserInfoById(ctx, &userv1.UpdateUserInfoByIdRequest{
		UserDomain: &userv1.UserDomain{
			Id:          req.Id,
			Password:    req.Password,
			NickName:    req.NickName,
			Phone:       req.Phone,
			Email:       req.Email,
			NewPassword: req.NewPassword,
		},
	})
	if err != nil {
		if errors.Is(err, errs.UseOldPassword) {
			ctx.JSON(http.StatusOK, Result[string]{
				Code: errs.ErrUseOldPassword,
				Msg:  "更新密码时，不要使用历史密码",
				Data: "error",
			})
			userHandler.logger.Error("更新密码时，不要使用历史密码", loggerx.Error(err),
				loggerx.String("method", "UserHandler:UpdateUserInfoById"))
			return
		}
		if strings.Contains(err.Error(), errs.PasswordWrong.Error()) {
			ctx.JSON(http.StatusOK, Result[string]{
				Code: errs.ErrUsernameOrPasswordError,
				Msg:  "旧密码输入错误",
				Data: "error",
			})
			userHandler.logger.Error("旧密码错误", loggerx.Error(err),
				loggerx.String("method", "UserHandler:UpdateUserInfoById"))
			return
		}
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("修改用户信息失败", loggerx.Error(err),
			loggerx.String("method", "UserHandler:UpdateUserInfoById"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "修改成功",
		Data: "success",
	})
}

func (userHandler *UserHandler) DeleteUSerInfoByIds(ctx *gin.Context) {
	// 获取用户信息
	info, ok := userHandler.jwtHandler.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method", "UserHandler:DeleteUSerInfoByIds"))
		return
	}
	if info.Authority != 1 {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("无权限删除用户信息",
			loggerx.String("method", "UserHandler:DeleteUSerInfoByIds"))
		return
	}
	// 获取路径参数，路径参数ids，多个id
	idsStr, ok := ctx.Params.Get("ids")
	if idsStr == " " {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "非法请求参数",
			Data: "error",
		})
		userHandler.logger.Error("没有传入要删除的id",
			loggerx.String("method", "UserHandler:DeleteUSerInfoByIds"))
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "非法请求参数",
			Data: "error",
		})
		userHandler.logger.Error("获取路径参数失败",
			loggerx.String("method", "UserHandler:DeleteUSerInfoByIds"))
		return
	}
	// 将字符串按照","分隔开，返回值为[]string
	idsSplit := strings.Split(idsStr, ",")
	var ids []int64
	// 将[]string转为[]int64
	for _, idStr := range idsSplit {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Result[string]{
				Code: errs.ErrSystemError,
				Msg:  "系统错误",
				Data: "error",
			})
			userHandler.logger.Error("string转int64失败", loggerx.Error(err),
				loggerx.String("method", "UserHandler:DeleteUSerInfoByIds"))
			return
		}
		ids = append(ids, id)
	}
	// 调用service层
	_, err := userHandler.userServiceClient.DeleteUserInfoByIds(ctx, &userv1.DeleteUserInfoByIdsRequest{
		Ids: ids,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("删除失败", loggerx.Error(err),
			loggerx.String("method", "UserHandler:DeleteUSerInfoByIds"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "删除成功",
		Data: "success",
	})
}

func (userHandler *UserHandler) RefreshToken(ctx *gin.Context) {
	userInfo, ok := userHandler.jwtHandler.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("用户信息不存在",
			loggerx.String("method", "UserHandler:RefreshToken"))
		return
	}
	userHandler.jwtHandler.SetJwt(ctx, userInfo, false)
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "token续约成功",
		Data: "success",
	})
}

func (userHandler *UserHandler) Logout(ctx *gin.Context) {
	userInfo, ok := userHandler.jwtHandler.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		userHandler.logger.Error("用户信息不存在",
			loggerx.String("method", "UserHandler:Logout"))
		return
	}
	_, err := userHandler.userServiceClient.Logout(ctx, &userv1.LogoutRequest{
		SessionId: userInfo.SessionId,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "退出失败",
			Data: "error",
		})
		userHandler.logger.Error("退出失败", loggerx.Error(err),
			loggerx.String("method", "UserHandler:Logout"))
		return
	}
	ctx.Header("jwt-short-token", " ")
	ctx.Header("jwt-long-token", " ")
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "退出成功",
		Data: "success",
	})
}

func (userHandler *UserHandler) LoginByPhoneCode(ctx *gin.Context) {
	type RegisterReq struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req RegisterReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "请求参数错误",
			Data: "error",
		})
		userHandler.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method", "UserHandler:LoginByPhoneCode"))
		return
	}
	// 校验手机号
	ok, err := userHandler.phoneRegexp.MatchString(req.Phone)
	if !ok {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "手机号格式错误",
			Data: "error",
		})
		userHandler.logger.Error("手机号格式错误", loggerx.Error(err),
			loggerx.String("method", "UserHandler:LoginByPhoneCode"))
		return
	}
	if len(req.Code) != 6 {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "验证码格式错误",
			Data: "error",
		})
		userHandler.logger.Error("验证码格式错误", loggerx.String("code", req.Code),
			loggerx.String("method", "UserHandler:LoginByPhoneCode"))
		return
	}
	ud, err := userHandler.userServiceClient.LoginByPhoneCode(ctx, &userv1.LoginByPhoneCodeRequest{
		UserDomain: &userv1.UserDomain{
			Phone: req.Phone,
			Code:  req.Code,
		},
		Biz: "login",
	})
	if errors.Is(err, ErrVerifyToMany) {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrCodeInputTooManyTimes,
			Msg:  "验证码输入错误次数过多",
			Data: "error",
		})
		userHandler.logger.Error("验证码输入错误次数过多", loggerx.Error(err),
			loggerx.String("method", "UserHandler:LoginByPhoneCode"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "验证码输入错误",
			Data: "error",
		})
		userHandler.logger.Error("登录失败", loggerx.Error(err),
			loggerx.String("method", "UserHandler:LoginByPhoneCode"))
		return
	}
	// 登陆成功
	// 设置token的负载，也就是存储在token中的数据
	userClaims := UserClaims{
		Id:        ud.GetUserDomain().Id,
		NickName:  ud.GetUserDomain().NickName,
		Authority: ud.GetUserDomain().Authority,
	}
	// 登陆成功，生成长短token
	userHandler.jwtHandler.SetJwt(ctx, userClaims, true)
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "登录成功",
		Data: "success",
	})
}
