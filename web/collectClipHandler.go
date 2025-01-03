package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"lifelog-grpc/collectClip/domain"
	"lifelog-grpc/collectClip/service"
	"lifelog-grpc/collectClip/vo"
	"lifelog-grpc/pkg/loggerx"
)

type CollectClipHandler struct {
	collectClipService service.CollectClipService
	biz                string
	JWTHandler
}

func NewCollectClipHandler(collectClipService service.CollectClipService) *CollectClipHandler {
	return &CollectClipHandler{
		collectClipService: collectClipService,
		biz:                "lifeLog",
	}
}

func (c *CollectClipHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/collectClip")
	// 编辑收藏夹（创建/update）
	rg.POST("/edit", c.EditCollectClip)
	// 删除收藏夹
	rg.DELETE("/:ids", c.DeleteCollectClip)
	// 获取收藏夹列表
	rg.POST("/list", c.CollectClipList)
	// 将LifeLog插入收藏夹
	rg.POST("/insert", c.InsertCollectClipDetail)
	// 查看收藏夹详情
	rg.POST("/detail", c.CollectClipDetail)
}

// EditCollectClip 编辑收藏夹（创建或更新）
func (c *CollectClipHandler) EditCollectClip(ctx *gin.Context) {
	type CollectReq struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	var req CollectReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:EditCollectClip"))
		return
	}
	// 获取用户信息
	userInfo, ok := c.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "ArticleHandler:EditCollectClip"))
		return
	}
	err = c.collectClipService.EditCollectClip(ctx.Request.Context(),
		domain.CollectClipDomain{
			Id:     req.Id,
			Name:   req.Name,
			UserId: userInfo.Id,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "编辑收藏夹失败",
			Data: "error",
		})
		c.logger.Error("编辑收藏夹失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:EditCollectClip"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "编辑收藏夹成功",
		Data: "success",
	})
}

// DeleteCollectClip 删除收藏夹
func (c *CollectClipHandler) DeleteCollectClip(ctx *gin.Context) {
	idsString, ok := ctx.Params.Get("ids")
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数获取失败",
			loggerx.String("method:", "ArticleHandler:DeleteCollectClip"))
		return
	}
	// 按照,分割
	idsStringSplit := strings.Split(idsString, ",")
	var ids []int64
	for _, idString := range idsStringSplit {
		// string转换为int64
		collectClipId, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Result[string]{
				Code: 400,
				Msg:  "参数错误",
				Data: "error",
			})
			c.logger.Error("string转为int64失败", loggerx.Error(err),
				loggerx.String("method:", "ArticleHandler:DeleteCollectClip"))
			return
		}
		ids = append(ids, collectClipId)
	}
	err := c.collectClipService.DeleteCollectClip(ctx.Request.Context(), ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "删除收藏夹失败",
			Data: "error",
		})
		c.logger.Error("删除收藏夹失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:DeleteCollectClip"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "删除收藏夹成功",
		Data: "success",
	})
}

// CollectClipList 获取收藏夹列表
func (c *CollectClipHandler) CollectClipList(ctx *gin.Context) {
	userInfo, ok := c.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "ArticleHandler:CollectClipList"))
		return
	}
	type ListReq struct {
		Limit  int
		Offset int
	}
	var req ListReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:CollectClipList"))
		return
	}
	cds, err := c.collectClipService.CollectClipList(ctx.Request.Context(), userInfo.Id,
		req.Limit, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "获取收藏夹列表失败",
			Data: "error",
		})
		c.logger.Error("获取收藏夹列表失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:CollectClipList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.CollectClipVo]{
		Code: 200,
		Msg:  "获取收藏夹列表成功",
		Data: c.collectClipsToCollectClipVo(cds),
	})
}

// collectClipsToCollectClipVo 将domain.CollectClipDomain转换为vo.CollectClipVo
func (c *CollectClipHandler) collectClipsToCollectClipVo(cds []domain.CollectClipDomain) []vo.CollectClipVo {
	ccvs := make([]vo.CollectClipVo, 0, len(cds))
	for _, cd := range cds {
		ccvs = append(ccvs, vo.CollectClipVo{
			Id:         cd.Id,
			Name:       cd.Name,
			UserId:     cd.UserId,
			Status:     cd.Status,
			CreateTime: cd.CreateTime,
			UpdateTime: cd.UpdateTime,
		})
	}
	return ccvs
}

// InsertCollectClipDetail 将LifeLog插入收藏夹
func (c *CollectClipHandler) InsertCollectClipDetail(ctx *gin.Context) {
	type InsertCollectReq struct {
		CollectId int64 `json:"collect_id"`
		LifeLogId int64 `json:"life_log_id"`
	}
	var req InsertCollectReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:InsertCollectClip"))
		return
	}
	err = c.collectClipService.InsertCollectClipDetail(ctx.Request.Context(),
		domain.CollectClipDetailDomain{
			CollectId: req.CollectId,
			LifeLogId: req.LifeLogId,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("插入收藏夹失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:InsertCollectClip"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "插入收藏夹成功",
		Data: "success",
	})
}

// CollectClipDetail 获取收藏夹详情
func (c *CollectClipHandler) CollectClipDetail(ctx *gin.Context) {
	type CollectClipDetailReq struct {
		CollectId int64 `json:"collect_id"`
		Limit     int   `json:"limit"`
		Offset    int   `json:"offset"`
	}
	var req CollectClipDetailReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:CollectClipDetail"))
		return
	}
	// 获取登录用户的id
	userInfo, ok := c.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "ArticleHandler:CollectClipDetail"))
		return
	}
	collectClipDetails, err := c.collectClipService.CollectClipDetail(ctx.Request.Context(), req.CollectId,
		req.Limit, req.Offset, userInfo.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("查询收藏夹详情失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:CollectClipDetail"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.CollectClipDetailVo]{
		Code: 200,
		Msg:  "查询成功",
		Data: c.collectClipDetailToCollectClipVo(collectClipDetails),
	})
}

// collectClipDetailToCollectClipVo 将domain.CollectClipDetailDomain转换为vo.CollectClipVo
func (c *CollectClipHandler) collectClipDetailToCollectClipVo(
	cds []domain.CollectClipDetailDomain) []vo.CollectClipDetailVo {
	ccvs := make([]vo.CollectClipDetailVo, 0, len(cds))
	for _, cd := range cds {
		ccvs = append(ccvs, vo.CollectClipDetailVo{
			CollectId:  cd.CollectId,
			LifeLogId:  cd.LifeLogId,
			CreateTime: cd.CreateTime,
			Id:         cd.Id,
			Status:     cd.Status,
			UpdateTime: cd.UpdateTime,
			PublicLifeLogVo: vo.PublicLifeLogVo{
				Title:    cd.PublicLifeLogDomain.Title,
				Content:  cd.PublicLifeLogDomain.Content,
				AuthorId: cd.PublicLifeLogDomain.AuthorId,
			},
		})
	}
	return ccvs
}
