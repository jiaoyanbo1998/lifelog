package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	codev1 "lifelog-grpc/api/proto/gen/code/v1"
	"lifelog-grpc/code/service"
	"lifelog-grpc/errs"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
)

// CodeHandler 短信处理器
type CodeHandler struct {
	codeServiceClient codev1.CodeServiceClient
	logger            loggerx.Logger
	phoneRegexp       *regexp.Regexp
}

// NewCodeHandler 构造短信处理器
func NewCodeHandler(l loggerx.Logger, codeServiceClient codev1.CodeServiceClient) *CodeHandler {
	// 定义正则表达式常量
	const PhoneRegexp = "^1[3-9]\\d{9}$"
	return &CodeHandler{
		logger: l,
		// 预编译正则表达式（性能优化：只需要编译一次，后续就可以重复使用）
		phoneRegexp:       regexp.MustCompile(PhoneRegexp, regexp.None),
		codeServiceClient: codeServiceClient,
	}
}

// RegisterRoutes 注册路由
func (codeHandler *CodeHandler) RegisterRoutes(server *gin.Engine) {
	// 用户路由组
	rg := server.Group("/code")
	// 发送验证码
	rg.POST("/send_code", codeHandler.SendPhoneCode)
}

func (codeHandler *CodeHandler) SendPhoneCode(ctx *gin.Context) {
	type RegisterReq struct {
		Phone string `json:"phone"`
		Biz   string `json:"biz"`
	}
	var req RegisterReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "系统错误",
			Data: "error",
		})
		codeHandler.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"))
		return
	}
	// 检查手机号是否为黑名单
	isBackPhone, err := codeHandler.codeServiceClient.IsBackPhone(ctx.Request.Context(), &codev1.IsBackPhoneRequest{
		Phone: req.Phone,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		codeHandler.logger.Error("查询是否为黑名单失败", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"))
		return
	}
	if isBackPhone.GetIsBack() {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrPhoneIsBlack,
			Msg:  "系统错误",
			Data: "error",
		})
		codeHandler.logger.Error("黑名单手机访问，怀疑是非法攻击", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"),
			loggerx.String("phone", req.Phone))
		return
	}
	// 校验手机号
	ok, err := codeHandler.phoneRegexp.MatchString(req.Phone)
	if !ok {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "手机号格式错误",
			Data: "error",
		})
		codeHandler.logger.Error("手机号格式错误", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		codeHandler.logger.Error("手机号正则表达式校验失败", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"))
		return
	}
	_, err = codeHandler.codeServiceClient.SendPhoneCode(ctx.Request.Context(), &codev1.SendPhoneCodeRequest{
		Phone: req.Phone,
		Biz:   req.Biz,
	})
	if errors.Is(err, service.ErrCodeSendFrequent) {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrCodeInputFrequently,
			Msg:  "系统错误",
			Data: "error",
		})
		codeHandler.logger.Error("发送验证码太频繁", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"))
		return
	}
	if errors.Is(err, service.ErrCodeSendMany) {
		_, err := codeHandler.codeServiceClient.SetBlackPhone(ctx.Request.Context(), &codev1.SetBlackPhoneRequest{
			Phone: req.Phone,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, Result[string]{
				Code: errs.ErrSystemError,
				Msg:  "系统错误",
				Data: "error",
			})
			codeHandler.logger.Error("设置手机号黑名单失败", loggerx.Error(err),
				loggerx.String("method", "CodeHandler:SendPhoneCode"))
			return
		}
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrCodeInputMany,
			Msg:  "系统错误",
			Data: "error",
		})
		codeHandler.logger.Error("发送验证码次数已达上限，明日再试", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"))
		return
	}
	if errors.Is(err, service.ErrLimited) {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrLimitHttp,
			Msg:  "系统错误",
			Data: "error",
		})
		codeHandler.logger.Error("触发限流", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		codeHandler.logger.Error("发送验证码失败", loggerx.Error(err),
			loggerx.String("method", "CodeHandler:SendPhoneCode"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "发送验证码成功",
		Data: "success",
	})
}
