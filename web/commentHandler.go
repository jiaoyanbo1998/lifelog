package web

import (
	"github.com/gin-gonic/gin"
	commentv1 "lifelog-grpc/api/proto/gen/comment/v1"
	"lifelog-grpc/comment/vo"
	"lifelog-grpc/errs"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
	"strconv"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	Biz                  string
	logger               loggerx.Logger
	commentServiceClient commentv1.CommentServiceClient
}

// NewCommentHandler 创建评论处理器
func NewCommentHandler(logger loggerx.Logger,
	commentServiceClient commentv1.CommentServiceClient) *CommentHandler {
	return &CommentHandler{
		Biz:                  "lifeLog",
		logger:               logger,
		commentServiceClient: commentServiceClient,
	}
}

// RegisterRoutes 注册路由
func (c *CommentHandler) RegisterRoutes(server *gin.Engine) {
	// 评论路由组
	rg := server.Group("/comment")
	// 创建评论
	rg.POST("/add", c.CreateComment)
	// 删除评论（级联删除，创建一个外键parent_id，引用本表的id，然后开启级联删除）
	rg.DELETE("/:id", c.DeleteComment)
	// 修改评论
	rg.PUT("/edit", c.EditComment)
	// 查找一级评论（parent_id==null）
	rg.POST("/FirstList", c.FirstList)
	// 查找根评论下的所有子孙评论，根据id降序排序（id小的评论肯定是更早发表的评论）
	rg.POST("/EveryRootChildSonList", c.EveryRootChildSonList)
	// 查询某个评论的，一级子评论
	rg.POST("/SonList", c.SonList)
}

func (c *CommentHandler) CreateComment(ctx *gin.Context) {
	type CommentReq struct {
		UserId   int64  `json:"user_id"`
		BizId    int64  `json:"biz_id"` // 文章id
		Content  string `json:"content"`
		ParentId int64  `json:"parent_id"`
		RootId   int64  `json:"root_id"`
	}
	var cq CommentReq
	err := ctx.Bind(&cq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "CommentHandler:CreateComment"))
		return
	}
	// 创建评论，发送评论到kafka
	_, err = c.commentServiceClient.ProducerCommentEvent(ctx.Request.Context(),
		&commentv1.ProducerCommentEventRequest{
			Comment: &commentv1.Comment{
				UserId:   cq.UserId,
				Biz:      c.Biz,
				BizId:    cq.BizId,
				Content:  cq.Content,
				ParentId: cq.ParentId,
				RootId:   cq.RootId,
			},
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("发送评论到kafka失败", loggerx.Error(err),
			loggerx.String("method:", "CommentHandler:CreateComment"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "创建成功",
		Data: "success",
	})
}

func (c *CommentHandler) DeleteComment(ctx *gin.Context) {
	getIdStr, ok := ctx.Params.Get("id")
	if getIdStr == " " || !ok {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("非法参数",
			loggerx.String("method", "CommentHandler:DeleteComment"))
		return
	}
	id, err := strconv.ParseInt(getIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "string转为int64失败",
			Data: "error",
		})
		c.logger.Error("string转int64失败", loggerx.Error(err),
			loggerx.String("method", "CommentHandler:DeleteComment"))
		return
	}
	_, err = c.commentServiceClient.DeleteComment(ctx.Request.Context(),
		&commentv1.DeleteCommentRequest{
			Id: id,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "删除评论失败",
			Data: "error",
		})
		c.logger.Error("删除评论失败", loggerx.Error(err),
			loggerx.String("method", "CommentHandler:DeleteComment"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "删除评论成功",
		Data: "success",
	})
}

func (c *CommentHandler) EditComment(ctx *gin.Context) {
	type CommentReq struct {
		Id      int64  `json:"id"`
		Content string `json:"content"`
	}
	var cq CommentReq
	err := ctx.Bind(&cq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err))
	}
	_, err = c.commentServiceClient.EditComment(ctx.Request.Context(),
		&commentv1.EditCommentRequest{
			Comment: &commentv1.Comment{
				Id:      cq.Id,
				Content: cq.Content,
			},
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "修改评论失败",
			Data: "error",
		})
		c.logger.Error("修改评论失败", loggerx.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "修改评论成功",
		Data: "success",
	})
}

func (c *CommentHandler) FirstList(ctx *gin.Context) {
	type FirstListReq struct {
		BizId int64 `json:"biz_id"`
		Min   int64 `json:"min"` // 只展示前几条数据
	}
	var req FirstListReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method", "CommentHandler:FirstList"))
		return
	}
	comments, err := c.commentServiceClient.FirstList(ctx.Request.Context(), &commentv1.FirstListRequest{
		Comment: &commentv1.Comment{
			Biz:   c.Biz,
			BizId: req.BizId,
		},
		Min: req.Min,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "获取评论失败",
			Data: "error",
		})
		c.logger.Error("获取评论失败", loggerx.Error(err),
			loggerx.String("method", "CommentHandler:FirstList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.CommentVo]{
		Code: 200,
		Msg:  "获取评论成功",
		Data: c.toVo(comments.GetComments()),
	})
}

func (c *CommentHandler) EveryRootChildSonList(ctx *gin.Context) {
	type FirstListSonReq struct {
		Id     int64 `json:"id"`
		RootId int64 `json:"root_id"`
		Limit  int64 `json:"limit"`
	}
	var req FirstListSonReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method", "CommentHandler:FirstListSon"))
		return
	}
	comments, err := c.commentServiceClient.EveryRootChildSonList(ctx.Request.Context(),
		&commentv1.EveryRootChildSonListRequest{
			Comment: &commentv1.Comment{
				Id:     req.Id,
				RootId: req.RootId,
			},
			Limit: req.Limit,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "获取评论失败",
			Data: "error",
		})
		c.logger.Error("获取评论失败", loggerx.Error(err),
			loggerx.String("method", "CommentHandler:FirstListSon"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.CommentVo]{
		Code: 200,
		Msg:  "获取评论成功",
		Data: c.toVo(comments.GetComments()),
	})
}

func (c *CommentHandler) SonList(ctx *gin.Context) {
	type SonListReq struct {
		ParentId int64 `json:"parent_id"`
		Offset   int64 `json:"offset"`
		Limit    int64 `json:"limit"`
	}
	var req SonListReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: errs.ErrInvalidParams,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method", "CommentHandler:SonList"))
		return
	}
	comments, err := c.commentServiceClient.SonList(ctx.Request.Context(), &commentv1.SonListRequest{
		Comment: &commentv1.Comment{
			ParentId: req.ParentId,
		},
		Limit:  req.Limit,
		Offset: req.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: errs.ErrSystemError,
			Msg:  "获取评论失败",
			Data: "error",
		})
		c.logger.Error("获取评论失败", loggerx.Error(err),
			loggerx.String("method", "CommentHandler:SonList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.CommentVo]{
		Code: 200,
		Msg:  "获取评论成功",
		Data: c.toVo(comments.GetComments()),
	})
}

func (c *CommentHandler) toVo(comments []*commentv1.Comment) []vo.CommentVo {
	var res []vo.CommentVo
	for _, comment := range comments {
		c := vo.CommentVo{
			Content: comment.GetContent(),
			Id:      comment.GetId(),
			UserId:  comment.GetUserId(),
		}
		res = append(res, c)
	}
	return res
}
