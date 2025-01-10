package web

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	filesv1 "lifelog-grpc/api/proto/gen/files/v1"
	"lifelog-grpc/errs"
	"lifelog-grpc/files/vo"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/pkg/miniox"
	"net/http"
)

type FilesHandler struct {
	filesv1.FilesServiceClient
	minio  *miniox.FileHandler
	logger loggerx.Logger
	JWTHandler
}

func NewFilesHandler(filesv1 filesv1.FilesServiceClient, minio *miniox.FileHandler, logger loggerx.Logger) *FilesHandler {
	return &FilesHandler{
		FilesServiceClient: filesv1,
		minio:              minio,
		logger:             logger,
	}
}

func (f *FilesHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/files")
	rg.POST("/upload", f.UploadFiles)
	rg.POST("/delete", f.FileDelete)
	rg.POST("/list", f.ListFiles)
	rg.GET("/search", f.searchFile)
}

func (f *FilesHandler) UploadFiles(ctx *gin.Context) {
	// 获取请求参数
	type uploadReq struct {
		Id      int64  `form:"id"`
		Name    string `form:"name"`
		Content string `form:"content"`
	}
	var req uploadReq
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: http.StatusBadRequest,
			Msg:  "参数绑定失败",
			Data: "error",
		})
		return
	}
	// 获取配置文件
	type config struct {
		Endpoint string `yaml:"endpoint"`
		UseSSL   bool   `yaml:"use_ssl"`
	}
	var c config
	err = viper.UnmarshalKey("minio", &c)
	if err != nil {
		f.logger.Error("读取配置文件失败", loggerx.Error(err))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "上传失败",
			Data: "error",
		})
		return
	}
	bucketName := "user"
	// 上传文件
	vals, err := f.minio.UploadFiles(ctx, c.Endpoint, bucketName, req.Name, c.UseSSL)
	if err != nil {
		f.logger.Error("上传文件失败", loggerx.Error(err))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "上传失败",
			Data: "error",
		})
		return
	}
	// 获取用户信息
	info, ok := f.JWTHandler.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	// 调用服务
	var er error
	for _, val := range vals {
		if req.Id == 0 {
			_, er = f.FilesServiceClient.CreateFile(ctx, &filesv1.CreateFileRequest{
				File: &filesv1.File{
					UserId:  info.Id,
					Url:     val.URL,
					Name:    req.Name,
					Content: req.Content,
				},
			})
		} else {
			_, er = f.FilesServiceClient.UpdateFile(ctx, &filesv1.UpdateFileRequest{
				File: &filesv1.File{
					Id:      req.Id,
					UserId:  info.Id,
					Url:     val.URL,
					Name:    req.Name,
					Content: req.Content,
				},
			})
			f.logger.Error("操作数据库失败：", loggerx.Error(er),
				loggerx.Int("userId", int(info.Id)),
				loggerx.String("name", req.Name),
				loggerx.String("url", val.URL))
		}
		if er != nil {
			ctx.JSON(http.StatusOK, Result[string]{
				Code: errs.ErrSystemError,
				Msg:  "上传失败",
				Data: "error",
			})
			return
		}
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 200,
			Msg:  "上传成功",
			Data: "success",
		})
	}
}

func (f *FilesHandler) FileDelete(ctx *gin.Context) {
	// 获取请求参数
	type uploadReq struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		// 扩展名
		Ext string `json:"ext"`
	}
	var req uploadReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: http.StatusBadRequest,
			Msg:  "参数绑定失败",
			Data: "error",
		})
		return
	}
	bucketName := "user"
	fileName := req.Name + req.Ext
	err = f.minio.DeleteFile(bucketName, fileName)
	if err != nil {
		f.logger.Error("删除文件失败", loggerx.Error(err))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "上传失败",
			Data: "error",
		})
		return
	}
	// 获取用户信息
	info, ok := f.JWTHandler.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	_, err = f.FilesServiceClient.DeleteFile(ctx, &filesv1.DeleteFileRequest{
		File: &filesv1.File{
			Id:     req.Id,
			UserId: info.Id,
		},
	})
	if err != nil {
		f.logger.Error("删除数据库失败", loggerx.Error(err))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "删除失败",
			Data: "error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "删除成功",
		Data: "success",
	})
}

func (f *FilesHandler) ListFiles(ctx *gin.Context) {
	type listReq struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	}
	var req listReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: http.StatusBadRequest,
			Msg:  "参数绑定失败",
			Data: "error",
		})
		return
	}
	// 获取用户信息
	info, ok := f.JWTHandler.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	res, err := f.FilesServiceClient.GetFileByUserId(ctx, &filesv1.GetFileByUserIdRequest{
		File: &filesv1.File{
			UserId: info.Id,
		},
		Limit:  int64(req.Limit),
		Offset: int64(req.Offset),
	})
	if err != nil {
		f.logger.Error("获取文件列表失败", loggerx.Error(err))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	// 将[]*File，转为[]FileVo
	fvs := make([]vo.FileVo, 0, len(res.GetFile()))
	for _, v := range res.GetFile() {
		fvs = append(fvs, vo.FileVo{
			Name:       v.Name,
			Content:    v.Content,
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
		})
	}
	ctx.JSON(http.StatusOK, Result[[]vo.FileVo]{
		Code: 200,
		Msg:  "获取成功",
		Data: fvs,
	})
}

func (f *FilesHandler) searchFile(ctx *gin.Context) {
	type listReq struct {
		Name string `form:"name"`
	}
	var req listReq
	// 将get请求，Url中的参数，绑定到req中
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: http.StatusBadRequest,
			Msg:  "参数绑定失败",
			Data: "error",
		})
		return
	}
	res, err := f.FilesServiceClient.GetFileByName(ctx, &filesv1.GetFileByNameRequest{
		File: &filesv1.File{
			Name: req.Name,
		},
	})
	if err != nil {
		f.logger.Error("获取文件失败", loggerx.Error(err))
		ctx.JSON(http.StatusOK, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result[vo.FileVo]{
		Code: 200,
		Msg:  "获取成功",
		Data: vo.FileVo{
			Name:       res.GetFile().Name,
			Content:    res.GetFile().Content,
			CreateTime: res.GetFile().CreateTime,
			UpdateTime: res.GetFile().UpdateTime,
		},
	})
}
