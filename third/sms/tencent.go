package sms

import (
	"context"
	"fmt"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// TencentSmsService 封装与短信相关的信息
type TencentSmsService struct {
	client   *sms.Client // 短信服务的客户端
	appId    *string     // 应用ID
	signName *string     // 签名
}

// NewTencentService 使用构造函数创建一个Service实例
func NewTencentService(c *sms.Client, appId string, signName string) *TencentSmsService {
	return &TencentSmsService{
		client:   c,
		appId:    &appId,
		signName: &signName,
	}
}

// Send 发送短信
//    ctx	上下文
//    tplId  模板id
//    args	 模板参数，使用参数替换模板中的占位符
//    numbers	接收短信的手机号，可以一次性给多个用户发送短信
func (s *TencentSmsService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 创建一个发送短信实例
	req := sms.NewSendSmsRequest()
	// 要给哪些手机号发送短信
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	// 短信服务，应用id
	req.SmsSdkAppId = s.appId
	// 设置请求上下文
	req.SetContext(ctx)
	// 设置模板参数，使用模板参数替换模板中的占位符
	req.TemplateParamSet = s.toStringPtrSlice(args)
	// 设置模板ID，要操作哪一个模板
	req.TemplateId = &tplId
	// 签名
	req.SignName = s.signName
	// 发送短信
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	// 检查是否有发送失败的情况
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败，code: %s, 原因：%s", *status.Code, *status.Message)
		}
	}
	// 全部成功发送
	return nil
}

// 将[]string，转换为[]*string
func (s *TencentSmsService) toStringPtrSlice(ss []string) []*string {
	// 使用len(ss)预先分配容量，避免频繁扩容
	ptrS := make([]*string, 0, len(ss))
	for _, s := range ss {
		ptrS = append(ptrS, &s)
	}
	return ptrS
}
