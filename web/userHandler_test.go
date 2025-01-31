package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	userv1 "lifelog-grpc/api/proto/gen/user/v1"
	userGRPCMock "lifelog-grpc/api/proto/gen/user/v1/mock"
	"lifelog-grpc/errs"
	"lifelog-grpc/event/lifeLogEvent"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func InitNoLogger() loggerx.Logger {
	return loggerx.NewZapNoLogger()
}

func SyncProducer() *lifeLogEvent.SyncProducer {
	return nil
}

func Test_UserHandler_RegisterByEmailAndPassword(t *testing.T) {
	// 定义测试用例
	testCase := []struct {
		name                  string
		mockUserServiceClient func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient
		reqBody               string
		wantCode              int
		wantBody              Result[string]
	}{
		{
			name: "注册成功",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 定义模拟方法
				mockUserServiceClient.EXPECT().RegisterByEmailAndPassword(
					gomock.Any(), // 第一个参数：上下文
					&userv1.RegisterByEmailAndPasswordRequest{
						UserDomain: &userv1.UserDomain{
							Email:    "1242644834@qq.com",
							Password: "123456As#",
						},
					},
				).Return(&userv1.RegisterByEmailAndPasswordResponse{}, nil)
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"email": "1242644834@qq.com",
				"password": "123456As#",
				"confirm_password": "123456As#"
			}`,
			wantCode: http.StatusOK,
			wantBody: Result[string]{
				Code: 200,
				Msg:  "注册成功",
				Data: "success",
			},
		},
		{
			name: "参数绑定错误",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"email": "1242644834@qq.com",
				"password": "123456As#"
				"confirm_password": "123456As#"
			`,
			wantCode: http.StatusBadRequest,
			wantBody: Result[string]{
				Code: errs.ErrInvalidParams,
				Msg:  "请求参数错误",
				Data: "error",
			},
		},
		{
			name: "邮箱为空",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"password": "123456As#",
				"confirm_password": "123456As#"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: Result[string]{
				Code: errs.ErrInvalidParams,
				Msg:  "邮箱不能为空",
				Data: "error",
			},
		},
		{
			name: "密码为空",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"email": "1242644834@qq.com",
				"confirm_password": "123456As#"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: Result[string]{
				Code: errs.ErrInvalidParams,
				Msg:  "密码不能为空",
				Data: "error",
			},
		},
		{
			name: "邮箱格式错误",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"email": "1242644834qq.com",
				"password": "123456As#",
				"confirm_password": "123456As#"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: Result[string]{
				Code: errs.ErrUserInputError,
				Msg:  "邮箱格式不正确",
				Data: "error",
			},
		},
		{
			name: "两次输入密码不同",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"email": "1242644834@qq.com",
				"password": "123456As#",
				"confirm_password": "123456As"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: Result[string]{
				Code: errs.ErrUserInputError,
				Msg:  "两次输入的密码不一致",
				Data: "error",
			},
		},
		{
			name: "密码格式错误",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"email": "1242644834@qq.com",
				"password": "123456As",
				"confirm_password": "123456As"
			}`,
			wantCode: http.StatusBadRequest,
			wantBody: Result[string]{
				Code: errs.ErrUserInputError,
				Msg:  "密码格式错误",
				Data: "error",
			},
		},
		{
			name: "邮箱已存在",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 定义模拟方法
				mockUserServiceClient.EXPECT().RegisterByEmailAndPassword(
					gomock.Any(), // 第一个参数：上下文
					&userv1.RegisterByEmailAndPasswordRequest{
						UserDomain: &userv1.UserDomain{
							Email:    "1242644834@qq.com",
							Password: "123456As#",
						},
					},
				).Return(&userv1.RegisterByEmailAndPasswordResponse{},
					fmt.Errorf("%s", errs.EmailExist.Error()))
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"email": "1242644834@qq.com",
				"password": "123456As#",
				"confirm_password": "123456As#"
			}`,
			wantCode: http.StatusOK,
			wantBody: Result[string]{
				Code: errs.ErrEmailAlreadyRegistered,
				Msg:  "邮箱已存在",
				Data: "error",
			},
		},
		{
			name: "注册失败",
			mockUserServiceClient: func(controller *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 创建模拟对象
				mockUserServiceClient := userGRPCMock.NewMockUserServiceClient(controller)
				// 定义模拟方法
				mockUserServiceClient.EXPECT().RegisterByEmailAndPassword(
					gomock.Any(), // 第一个参数：上下文
					&userv1.RegisterByEmailAndPasswordRequest{
						UserDomain: &userv1.UserDomain{
							Email:    "1242644834@qq.com",
							Password: "123456As#",
						},
					},
				).Return(&userv1.RegisterByEmailAndPasswordResponse{},
					fmt.Errorf("%s", "注册错误"))
				// 返回模拟对象
				return mockUserServiceClient
			},
			reqBody: `{
				"email": "1242644834@qq.com",
				"password": "123456As#",
				"confirm_password": "123456As#"
			}`,
			wantCode: http.StatusOK,
			wantBody: Result[string]{
				Code: errs.ErrSystemError,
				Msg:  "系统错误",
				Data: "error",
			},
		},
	}
	// 循环使用测试用例
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			// 创建一个web服务器
			server := gin.Default()
			// 创建mock控制器
			controller := gomock.NewController(t)
			// 关闭mock控制器
			defer controller.Finish()
			// 创建UserServiceClient的mock对象
			mockUserServiceClient := tc.mockUserServiceClient(controller)
			// 创建UserHandler的实例对象
			userHandler := NewUserHandler(mockUserServiceClient,
				InitNoLogger(), nil, nil, nil)
			// 路由注册
			userHandler.RegisterRoutes(server)
			// 创建http请求
			request, err := http.NewRequest(
				http.MethodPost,
				"/user/register_email_password",
				bytes.NewBuffer([]byte(tc.reqBody)),
			)
			// 此处不会有错误，如果发生错误，直接panic
			require.NoError(t, err)
			// 指定请求体的数据类型为json
			request.Header.Set("Content-Type", "application/json")
			// 创建http响应
			response := httptest.NewRecorder()
			// 执行http请求(request)，并将响应结果写入到response中
			server.ServeHTTP(response, request)
			// 反序列化响应结果
			var responseBody Result[string]
			err = json.NewDecoder(response.Body).Decode(&responseBody)
			require.NoError(t, err)
			// 断言响应结果，想要的错误和得到的错误是否相等
			assert.Equal(t, tc.wantCode, response.Code)
			// 断言响应结果，想要的响应结果和得到的响应结果是否相等
			assert.Equal(t, tc.wantBody, responseBody)
		})
	}
}

func Test_UserHandler_LoginByEmailAndPassword(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name                  string
		mockUserServiceClient func(ctrl *gomock.Controller) *userGRPCMock.MockUserServiceClient
		reqBody               string
		wantCode              int
		wantBody              Result[string]
	}{
		{
			name: "参数绑定错误",
			mockUserServiceClient: func(ctrl *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				// 不期望调用 UserServiceClient
				return userGRPCMock.NewMockUserServiceClient(ctrl)
			},
			reqBody:  `{ "email": "test@example.com", "password": "pass123A#" `, // 无效的 JSON
			wantCode: http.StatusBadRequest,
			wantBody: Result[string]{
				Code: errs.ErrSystemError,
				Msg:  "系统错误",
				Data: "error",
			},
		},
		{
			name: "用户不存在",
			mockUserServiceClient: func(ctrl *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				mockClient := userGRPCMock.NewMockUserServiceClient(ctrl)
				mockClient.EXPECT().LoginByEmailAndPassword(
					gomock.Any(),
					&userv1.LoginByEmailAndPasswordRequest{
						UserDomain: &userv1.UserDomain{
							Email:    "nonexistent@example.com",
							Password: "pass123A#",
						},
					},
				).Return(nil, status.Error(codes.NotFound, errs.UserNotExist.Error()))
				return mockClient
			},
			reqBody:  `{ "email": "nonexistent@example.com", "password": "pass123A#" }`,
			wantCode: http.StatusOK,
			wantBody: Result[string]{
				Code: errs.ErrUserNotRegistered,
				Msg:  "用户没有被注册",
				Data: "error",
			},
		},
		{
			name: "用户名或密码错误",
			mockUserServiceClient: func(ctrl *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				mockClient := userGRPCMock.NewMockUserServiceClient(ctrl)
				mockClient.EXPECT().LoginByEmailAndPassword(
					gomock.Any(),
					&userv1.LoginByEmailAndPasswordRequest{
						UserDomain: &userv1.UserDomain{
							Email:    "user@example.com",
							Password: "wrongpass",
						},
					},
				).Return(nil, errors.New("invalid credentials"))
				return mockClient
			},
			reqBody:  `{ "email": "user@example.com", "password": "wrongpass" }`,
			wantCode: http.StatusOK,
			wantBody: Result[string]{
				Code: errs.ErrUsernameOrPasswordError,
				Msg:  "用户名或密码错误",
				Data: "error",
			},
		},
		{
			name: "登录成功",
			mockUserServiceClient: func(ctrl *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				mockClient := userGRPCMock.NewMockUserServiceClient(ctrl)
				mockClient.EXPECT().LoginByEmailAndPassword(
					gomock.Any(),
					&userv1.LoginByEmailAndPasswordRequest{
						UserDomain: &userv1.UserDomain{
							Email:    "test@example.com",
							Password: "correctPass#123",
						},
					},
				).Return(&userv1.LoginByEmailAndPasswordResponse{
					UserDomain: &userv1.UserDomain{
						Id:        1,
						NickName:  "test",
						Authority: 1,
					},
				}, nil)
				return mockClient
			},
			reqBody:  `{ "email": "test@example.com", "password": "correctPass#123" }`,
			wantCode: http.StatusOK,
			wantBody: Result[string]{
				Code: 200,
				Msg:  "登录成功",
				Data: "success",
			},
		},
		{
			name: "系统错误",
			mockUserServiceClient: func(ctrl *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				mockClient := userGRPCMock.NewMockUserServiceClient(ctrl)
				mockClient.EXPECT().LoginByEmailAndPassword(
					gomock.Any(),
					&userv1.LoginByEmailAndPasswordRequest{
						UserDomain: &userv1.UserDomain{
							Email:    "test@example.com",
							Password: "correctPass#123",
						},
					},
				).Return(nil, status.Error(codes.Internal, "internal error"))
				return mockClient
			},
			reqBody:  `{ "email": "test@example.com", "password": "correctPass#123" }`,
			wantCode: http.StatusOK,
			wantBody: Result[string]{
				Code: errs.ErrUsernameOrPasswordError,
				Msg:  "用户名或密码错误",
				Data: "error",
			},
		},
		{
			name: "登录成功",
			mockUserServiceClient: func(ctrl *gomock.Controller) *userGRPCMock.MockUserServiceClient {
				mockClient := userGRPCMock.NewMockUserServiceClient(ctrl)
				mockClient.EXPECT().LoginByEmailAndPassword(
					gomock.Any(),
					&userv1.LoginByEmailAndPasswordRequest{
						UserDomain: &userv1.UserDomain{
							Email:    "1242644834@qq.com",
							Password: "123456As#",
						},
					},
				).Return(&userv1.LoginByEmailAndPasswordResponse{
					UserDomain: &userv1.UserDomain{
						Id:        1,
						NickName:  "test",
						Authority: 1,
					},
					Info: "需要修改密码",
				}, nil)
				return mockClient
			},
			reqBody: `{
				"email": "1242644834@qq.com",
				"password": "123456As#"
			}`,
			wantCode: http.StatusOK,
			wantBody: Result[string]{
				Code: 200,
				Msg:  "登录成功",
				Data: "你的密码太久没有修改了，你需要去修改密码",
			},
		},
	}
	// 遍历测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 创建 Gin 引擎
			server := gin.Default()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 初始化模拟对象
			mockUserServiceClient := tc.mockUserServiceClient(ctrl)
			// 创建 UserHandler 实例
			userHandler := NewUserHandler(
				mockUserServiceClient,
				InitNoLogger(),
				nil,
				nil,
				SyncProducer(),
			)
			userHandler.RegisterRoutes(server)

			// 构造请求
			request, err := http.NewRequest(
				http.MethodPost,
				"/user/login_email_password", // 根据实际路由调整
				bytes.NewBuffer([]byte(tc.reqBody)),
			)
			// 此处不会有错误，如果发生错误，直接panic
			require.NoError(t, err)
			// 指定请求体的数据类型为json
			request.Header.Set("Content-Type", "application/json")
			// 创建http响应
			response := httptest.NewRecorder()
			// 执行http请求(request)，并将响应结果写入到response中
			server.ServeHTTP(response, request)
			// 反序列化响应结果
			var responseBody Result[string]
			err = json.NewDecoder(response.Body).Decode(&responseBody)
			require.NoError(t, err)
			// 断言响应结果，想要的错误和得到的错误是否相等
			assert.Equal(t, tc.wantCode, response.Code)
			// 断言响应结果，想要的响应结果和得到的响应结果是否相等
			assert.Equal(t, tc.wantBody, responseBody)
		})
	}
}
