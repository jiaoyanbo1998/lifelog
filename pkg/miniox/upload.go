package miniox

import (
	"context"
	"github.com/gin-gonic/gin"
	"lifelog-grpc/errs"
	"lifelog-grpc/pkg/loggerx"
	"lifelog-grpc/web"
	"net/http"
	"strconv"
	"time"
)

type UploadHandler struct {
	minio  *MinioClient
	logger loggerx.Logger
}

func NewUploadHandler(minio *MinioClient, logger loggerx.Logger) *UploadHandler {
	return &UploadHandler{
		minio:  minio,
		logger: logger,
	}
}

func (uh *UploadHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/upload", uh.uploadFiles)
}

func (uh *UploadHandler) uploadFiles(ctx *gin.Context) {
	type uploadFileReq struct {
		Filename    string `json:"filename"`    // 文件名
		TotalChunks int    `json:"totalChunks"` // 分片总数
		bucket      string `json:"bucket"`      // 桶名称
		FileSize    int    `json:"fileSize"`    // 文件大小
	}
	var req uploadFileReq
	// 自动选择合适的解析器，将http请求参数绑定到结构体
	err := ctx.ShouldBind(&req)
	if err != nil {
		uh.logger.Error("参数绑定失败", loggerx.Error(err))
		return
	}
	// 获取文件信息
	multipartForm, err := ctx.MultipartForm()
	if err != nil {
		uh.logger.Error("获取文件信息失败", loggerx.Error(err))
		return
	}
	// 获取文件列表
	files := multipartForm.File["file"]
	// 判断文件列表是否为空
	if len(files) == 0 {
		uh.logger.Warn("没有上传文件")
		return
	}
	// 桶名称
	bucket := req.bucket
	// 文件名
	fileName := req.Filename
	// 文件类型
	contentType := files[0].Header.Get("Content-Type")
	// 文件路径
	path := ""
	// 单个文件上传
	if req.TotalChunks == 1 {
		// 超时
		ctx1, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		// 打开文件
		open, _ := files[0].Open()
		// 读取文件内容
		buf := make([]byte, req.FileSize)
		n, _ := open.Read(buf)
		// 上传文件
		info, er := uh.minio.Upload(ctx1, bucket, fileName, contentType, buf[:n])
		if er != nil {
			ctx.JSON(http.StatusOK, web.Result[string]{
				Code: errs.ErrSystemError,
				Msg:  "上传失败",
				Data: "error",
			})
			return
		}
		// 返回文件路径
		path = info.Bucket + "/" + req.Filename
		cancel()
		open.Close()
	} else {
		// 多个文件上传
		sizeSum := req.FileSize
		size := req.FileSize / req.TotalChunks
		for i := 0; i < req.TotalChunks; i++ {
			// 每一片的内存大小
			if i == req.TotalChunks-1 {
				size = sizeSum
			}
			buf := make([]byte, size)
			// 超时
			ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
			// 打开文件
			open, _ := files[i].Open()
			// 读取文件内容
			n, _ := open.Read(buf)
			// 上传文件
			// 编号
			idx := strconv.Itoa(i)
			info, er := uh.minio.Upload(ctx2, bucket, fileName+idx, contentType, buf[:n])
			if er != nil {
				ctx.JSON(http.StatusOK, web.Result[string]{
					Code: errs.ErrSystemError,
					Msg:  "上传失败",
					Data: "error",
				})
				return
			}
			// 返回文件路径
			if i == req.TotalChunks-1 {
				path = info.Bucket + "/" + req.Filename + idx
			}
			cancel2()
			// 关闭文件
			open.Close()
			sizeSum = sizeSum - size
		}
		// 合并文件
		ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel2()
		errr := uh.minio.Compose(ctx2, bucket, fileName, contentType, req.TotalChunks)
		if errr != nil {
			ctx.JSON(http.StatusOK, web.Result[string]{
				Code: errs.ErrSystemError,
				Msg:  "上传失败",
				Data: "error",
			})
		}
	}
	ctx.JSON(http.StatusOK, web.Result[string]{
		Code: 200,
		Msg:  "上传成功",
		Data: path,
	})
}
