package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	collectv1 "lifelog-grpc/api/proto/gen/api/proto/collect/v1"
	"lifelog-grpc/collect/vo"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
	"strconv"
	"strings"
)

type CollectHandler struct {
	collectServiceClient collectv1.CollectServiceClient
	biz                  string
	JWTHandler
}

func NewCollectHandler(collectServiceClient collectv1.CollectServiceClient) *CollectHandler {
	return &CollectHandler{
		collectServiceClient: collectServiceClient,
		biz:                  "lifeLog",
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
			loggerx.String("method:", "CollectHandler:EditCollect"))
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
			loggerx.String("method：", "CollectHandler:EditCollect"))
		return
	}
	_, err = c.collectServiceClient.EditCollect(ctx.Request.Context(),
		&collectv1.EditCollectRequest{
			Collect: &collectv1.Collect{
				CollectId: req.Id,
				AuthorId:  userInfo.Id,
				Name:      req.Name,
			},
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "编辑收藏夹失败",
			Data: "error",
		})
		c.logger.Error("编辑收藏夹失败", loggerx.Error(err),
			loggerx.String("method:", "CollectHandler:EditCollect"))
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
			loggerx.String("method:", "CollectHandler:DeleteCollect"))
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
				loggerx.String("method:", "CollectHandler:DeleteCollect"))
			return
		}
		ids = append(ids, collectId)
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
			loggerx.String("method：", "CollectHandler:DeleteCollect"))
		return
	}
	_, err := c.collectServiceClient.DeleteCollect(ctx.Request.Context(),
		&collectv1.DeleteCollectRequest{
			Ids:      ids,
			AuthorId: userInfo.Id,
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "删除收藏夹失败",
			Data: "error",
		})
		c.logger.Error("删除收藏夹失败", loggerx.Error(err),
			loggerx.String("method:", "CollectHandler:DeleteCollect"))
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
			loggerx.String("method：", "CollectHandler:CollectList"))
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
			loggerx.String("method:", "CollectHandler:CollectList"))
		return
	}
	cds, err := c.collectServiceClient.CollectList(ctx.Request.Context(),
		&collectv1.CollectListRequest{
			AuthorId: userInfo.Id,
			Limit:    int64(req.Limit),
			Offset:   int64(req.Offset),
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "获取收藏夹列表失败",
			Data: "error",
		})
		c.logger.Error("获取收藏夹列表失败", loggerx.Error(err),
			loggerx.String("method:", "CollectHandler:CollectList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.CollectVo]{
		Code: 200,
		Msg:  "获取收藏夹列表成功",
		Data: c.collectsToCollectVo(cds.GetCollects()),
	})
}

// collectsToCollectVo 将domain.CollectDomain转换为vo.CollectVo
func (c *CollectHandler) collectsToCollectVo(cds []*collectv1.Collect) []vo.CollectVo {
	ccvs := make([]vo.CollectVo, 0, len(cds))
	for _, cd := range cds {
		ccvs = append(ccvs, vo.CollectVo{
			Id:         cd.CollectId,
			Name:       cd.Name,
			UserId:     cd.AuthorId,
			Status:     uint8(cd.Status),
			CreateTime: cd.CreateTime,
			UpdateTime: cd.UpdateTime,
		})
	}
	return ccvs
}

// InsertCollectDetail 插入收藏夹详情（将LifeLog插入收藏夹）
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
			loggerx.String("method:", "CollectHandler:InsertCollectDetail"))
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
			loggerx.String("method：", "CollectHandler:InsertCollectDetail"))
		return
	}
	_, err = c.collectServiceClient.InsertCollectDetail(ctx.Request.Context(),
		&collectv1.InsertCollectDetailRequest{
			Collect: &collectv1.Collect{
				CollectId: req.CollectId,
				AuthorId:  userInfo.Id,
			},
			CollectDetail: &collectv1.CollectDetail{
				LifeLogId: req.LifeLogId,
			},
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("插入收藏夹失败", loggerx.Error(err),
			loggerx.String("method:", "CollectHandler:InsertCollectDetail"))
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
			loggerx.String("method:", "CollectHandler:CollectDetail"))
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
			loggerx.String("method：", "CollectHandler:CollectDetail"))
		return
	}
	collectDetails, err := c.collectServiceClient.CollectDetail(ctx.Request.Context(),
		&collectv1.CollectDetailRequest{
			Collect: &collectv1.Collect{
				CollectId: req.CollectId,
				AuthorId:  userInfo.Id,
			},
			Limit:  int64(req.Limit),
			Offset: int64(req.Offset),
		})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		c.logger.Error("查询收藏夹详情失败", loggerx.Error(err),
			loggerx.String("method:", "CollectHandler:CollectDetail"))
		return
	}
	// 将collectv1.collectDetail转换为collectDetailVo
	plvs := make([]vo.PublicLifeLogVo, 0, len(collectDetails.GetPublicLifeLogs()))
	copier.Copy(&plvs, collectDetails.GetPublicLifeLogs())
	ctx.JSON(http.StatusOK, Result[vo.CollectDetailVo]{
		Code: 200,
		Msg:  "查询成功",
		Data: vo.CollectDetailVo{
			Id:              collectDetails.GetCollectDetail().GetId(),
			CollectId:       collectDetails.GetCollectDetail().GetCollectId(),
			LifeLogId:       collectDetails.GetCollectDetail().GetLifeLogId(),
			UpdateTime:      collectDetails.GetCollectDetail().GetUpdateTime(),
			CreateTime:      collectDetails.GetCollectDetail().GetCreateTime(),
			Status:          uint8(collectDetails.GetCollectDetail().GetStatus()),
			PublicLifeLogVo: plvs,
		},
	})
}
