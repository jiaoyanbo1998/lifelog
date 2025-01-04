package web

import (
	"github.com/gin-gonic/gin"
	"lifelog-grpc/collect/domain"
	"lifelog-grpc/collect/service"
	"lifelog-grpc/collect/vo"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
	"strconv"
	"strings"
)

type CollectHandler struct {
	collectService service.CollectService
	biz            string
	JWTHandler
}

func NewCollectHandler(collectService service.CollectService) *CollectHandler {
	return &CollectHandler{
		collectService: collectService,
		biz:            "lifeLog",
	}
}

func (c *CollectHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/collect")
	// 编辑收藏夹（创建/update）
	rg.POST("/edit", c.EditCollect)
	// 删除收藏夹
	rg.DELETE("/:ids", c.DeleteCollect)
	// 获取收藏夹列表
	rg.POST("/list", c.CollectList)
	// 将LifeLog插入收藏夹
	rg.POST("/insert", c.InsertCollectDetail)
	// 查看收藏夹详情
	rg.POST("/detail", c.CollectDetail)
}

// EditCollect 编辑收藏夹（创建或更新）
func (c *CollectHandler) EditCollect(ctx *gin.Context) {
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
			loggerx.String("method:", "ArticleHandler:EditCollect"))
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
			loggerx.String("method：", "ArticleHandler:EditCollect"))
		return
	}
	err = c.collectService.EditCollect(ctx.Request.Context(),
		domain.CollectDomain{
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
			loggerx.String("method:", "ArticleHandler:EditCollect"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "编辑收藏夹成功",
		Data: "success",
	})
}

// DeleteCollect 删除收藏夹
func (c *CollectHandler) DeleteCollect(ctx *gin.Context) {
	idsString, ok := ctx.Params.Get("ids")
	if !ok {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数获取失败",
			loggerx.String("method:", "ArticleHandler:DeleteCollect"))
		return
	}
	// 按照,分割
	idsStringSplit := strings.Split(idsString, ",")
	var ids []int64
	for _, idString := range idsStringSplit {
		// string转换为int64
		collectId, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, Result[string]{
				Code: 400,
				Msg:  "参数错误",
				Data: "error",
			})
			c.logger.Error("string转为int64失败", loggerx.Error(err),
				loggerx.String("method:", "ArticleHandler:DeleteCollect"))
			return
		}
		ids = append(ids, collectId)
	}
	err := c.collectService.DeleteCollect(ctx.Request.Context(), ids)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "删除收藏夹失败",
			Data: "error",
		})
		c.logger.Error("删除收藏夹失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:DeleteCollect"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "删除收藏夹成功",
		Data: "success",
	})
}

// CollectList 获取收藏夹列表
func (c *CollectHandler) CollectList(ctx *gin.Context) {
	userInfo, ok := c.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "ArticleHandler:CollectList"))
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
			loggerx.String("method:", "ArticleHandler:CollectList"))
		return
	}
	cds, err := c.collectService.CollectList(ctx.Request.Context(), userInfo.Id,
		req.Limit, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "获取收藏夹列表失败",
			Data: "error",
		})
		c.logger.Error("获取收藏夹列表失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:CollectList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.CollectVo]{
		Code: 200,
		Msg:  "获取收藏夹列表成功",
		Data: c.collectsToCollectVo(cds),
	})
}

// collectsToCollectVo 将domain.CollectDomain转换为vo.CollectVo
func (c *CollectHandler) collectsToCollectVo(cds []domain.CollectDomain) []vo.CollectVo {
	ccvs := make([]vo.CollectVo, 0, len(cds))
	for _, cd := range cds {
		ccvs = append(ccvs, vo.CollectVo{
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

// InsertCollectDetail 将LifeLog插入收藏夹
func (c *CollectHandler) InsertCollectDetail(ctx *gin.Context) {
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
			loggerx.String("method:", "ArticleHandler:InsertCollect"))
		return
	}
	err = c.collectService.InsertCollectDetail(ctx.Request.Context(),
		domain.CollectDetailDomain{
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
			loggerx.String("method:", "ArticleHandler:InsertCollect"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "插入收藏夹成功",
		Data: "success",
	})
}

// CollectDetail 获取收藏夹详情
func (c *CollectHandler) CollectDetail(ctx *gin.Context) {
	type CollectDetailReq struct {
		CollectId int64 `json:"collect_id"`
		Limit     int   `json:"limit"`
		Offset    int   `json:"offset"`
	}
	var req CollectDetailReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "参数错误",
			Data: "error",
		})
		c.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:CollectDetail"))
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
			loggerx.String("method：", "ArticleHandler:CollectDetail"))
		return
	}
	collectDetails, err := c.collectService.CollectDetail(ctx.Request.Context(), req.CollectId,
		req.Limit, req.Offset, userInfo.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("查询收藏夹详情失败", loggerx.Error(err),
			loggerx.String("method:", "ArticleHandler:CollectDetail"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.CollectDetailVo]{
		Code: 200,
		Msg:  "查询成功",
		Data: c.collectDetailToCollectVo(collectDetails),
	})
}

// collectDetailToCollectVo 将domain.CollectDetailDomain转换为vo.CollectVo
func (c *CollectHandler) collectDetailToCollectVo(
	cds []domain.CollectDetailDomain) []vo.CollectDetailVo {
	ccvs := make([]vo.CollectDetailVo, 0, len(cds))
	for _, cd := range cds {
		ccvs = append(ccvs, vo.CollectDetailVo{
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
