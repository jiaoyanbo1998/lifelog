package web

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	interactivev1 "lifelog-grpc/api/proto/gen/api/proto/interactive/v1"
	lifelogv1 "lifelog-grpc/api/proto/gen/api/proto/lifelog/v1"
	"lifelog-grpc/event/lifeLogEvent"
	"lifelog-grpc/job"
	"lifelog-grpc/lifeLog/vo"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
	"strconv"
	"time"
)

type LifeLogHandler struct {
	logger                   loggerx.Logger
	lifeLogServiceClient     lifelogv1.LifeLogServiceClient
	interactiveServiceClient interactivev1.InteractiveServiceClient
	job                      job.Job
	biz                      string
	syncProducer             lifeLogEvent.Producer
	JWTHandler
}

func NewLifeLogHandler(l loggerx.Logger,
	lifeLogServiceClient lifelogv1.LifeLogServiceClient,
	job job.Job,
	syncProducer lifeLogEvent.Producer,
	interactiveServiceClient interactivev1.InteractiveServiceClient) *LifeLogHandler {
	return &LifeLogHandler{
		logger:                   l,
		lifeLogServiceClient:     lifeLogServiceClient,
		interactiveServiceClient: interactiveServiceClient,
		biz:                      "lifeLog",
		job:                      job,
		syncProducer:             syncProducer,
	}
}

func (l *LifeLogHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/lifeLog")
	// 创作者的接口
	// 编辑LifeLog（新建LifeLog/修改LifeLog）
	rg.POST("/edit", l.Edit)
	// 删除LifeLog（根据Id删除LifeLog，可以删除线上库/制作库中的LifeLog）
	rg.POST("/delete", l.Delete)
	// 撤销LifeLog（将指定id的LifeLog状态更换为"隐藏"，对读者不可见）
	rg.PUT("/revoke/:id", l.Revoke)
	// 发布LifeLog（只有发布后的LifeLog读者才能看到）
	rg.POST("/publish", l.Publish)
	// 作者的LifeLog列表（根据作者id分页查询制作库中的LifeLog）
	rg.POST("/author_id", l.DraftList)
	// 读者的接口
	// 根据标题查找LifeLog
	rg.POST("/title", l.SearchByTitle)
	// 热榜
	rg.GET("/hot", l.Hot)
	// 查看LifeLog的详情 （根据id查找LifeLog，线上库/制作库）
	rg.POST("/detail", l.Detail)
}

// Edit 编辑LifeLog（新建LifeLog/修改LifeLog）
func (l *LifeLogHandler) Edit(ctx *gin.Context) {
	type EditLifeLogReq struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req EditLifeLogReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("参数bind失败",
			loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:Edit"))
		return
	}
	// 获取token中存储的用户信息
	userInfo, ok := l.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:Edit"))
		return
	}
	// 调用service层代码
	// 传入id，表示是修改文化在哪个
	// 不传入id，表示创建新LifeLog
	res, err := l.lifeLogServiceClient.Edit(ctx, &lifelogv1.EditRequest{
		LifeLogDomain: &lifelogv1.LifeLogDomain{
			Id:      req.Id,
			Title:   req.Title,
			Content: req.Content,
			Author: &lifelogv1.Author{
				UserId: userInfo.Id,
			},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("编辑LifeLog失败",
			loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:Edit"))
		return
	}
	ctx.JSON(http.StatusOK, Result[vo.LifeLogVo]{
		Code: 200,
		Msg:  "编辑LifeLog成功",
		Data: vo.LifeLogVo{
			Id:      res.GetLifeLogDomain().Id,
			Content: res.GetLifeLogDomain().Content,
			// 将毫秒值时间戳，转换为，time.Time类型(2024-09-23 16:00:00 +0800 CST)
			CreateTime: time.UnixMilli(res.GetLifeLogDomain().CreateTime),
			Title:      res.GetLifeLogDomain().Title,
			UpdateTime: time.UnixMilli(res.GetLifeLogDomain().UpdateTime),
			AuthorId:   res.GetLifeLogDomain().GetAuthor().GetUserId(),
			AuthorName: res.GetLifeLogDomain().GetAuthor().GetNickName(),
		},
	})
}

// Delete 删除LifeLog（根据Id删除LifeLog，可以删除线上库/制作库中的LifeLog）
func (l *LifeLogHandler) Delete(ctx *gin.Context) {
	type DeleteLifeLogReq struct {
		Ids []int64 `json:"ids"`
		// 操作的是线上库，还是草稿库
		// true 线上库，false 制作库
		Public bool `json:"public"`
	}
	var req DeleteLifeLogReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("参数bind失败",
			loggerx.String("method:", "LifeLogHandler:Delete"))
		return
	}
	// 调用service层代码
	_, err = l.lifeLogServiceClient.Delete(ctx, &lifelogv1.DeleteRequest{
		Ids:      req.Ids,
		IsPublic: req.Public,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "删除LifeLog失败",
			Data: "error",
		})
		l.logger.Error("删除LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Delete"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "删除LifeLog成功",
		Data: "success",
	})
}

// SearchByTitle 查找LifeLog（根据title）
func (l *LifeLogHandler) SearchByTitle(ctx *gin.Context) {
	type SearchReq struct {
		Limit  int64  `json:"limit"`  // 每页显示的条数
		Offset int64  `json:"offset"` // 偏移量
		Title  string `json:"title"`
	}
	var req SearchReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("参数bind失败",
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	res, err := l.lifeLogServiceClient.SearchByTitle(ctx, &lifelogv1.SearchByTitleRequest{
		LifeLogDomain: &lifelogv1.LifeLogDomain{
			Title:  req.Title,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "查找LifeLog失败",
			Data: "error",
		})
		l.logger.Error("查找LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	// 将[]*lifelogv1.LifeLogDomain转换为[]vo.LifeLogVo
	llvs := make([]vo.LifeLogVo, 0, len(res.GetLifeLogDomain()))
	err = copier.Copy(&llvs, res.GetLifeLogDomain())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("copier失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:SearchByTitle"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.LifeLogVo]{
		Code: 200,
		Msg:  "查找LifeLog成功",
		Data: llvs,
	})
}

// DraftList 作者的LifeLog列表（根据作者id分页查询制作库中的LifeLog）
func (l *LifeLogHandler) DraftList(ctx *gin.Context) {
	type SearchReq struct {
		Limit  int64 `json:"limit"`  // 每页显示的条数
		Offset int64 `json:"offset"` // 偏移量
	}
	var req SearchReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("参数bind失败",
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	// 获取token中存储的用户信息
	userInfo, ok := l.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:Search"))
		return
	}
	res, err := l.lifeLogServiceClient.DraftList(ctx, &lifelogv1.DraftListRequest{
		LifeLogDomain: &lifelogv1.LifeLogDomain{
			Author: &lifelogv1.Author{
				UserId: userInfo.Id,
			},
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "查找LifeLog列表失败",
			Data: "error",
		})
		l.logger.Error("查找LifeLog列表失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	// 将[]*lifelogv1.LifeLogDomain转换为[]vo.LifeLogVo
	llvs := make([]vo.LifeLogVo, 0, len(res.GetLifeLogDomain()))
	err = copier.Copy(&llvs, res.GetLifeLogDomain())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("copier失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:DraftList"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.LifeLogVo]{
		Code: 200,
		Msg:  "查找LifeLog列表成功",
		Data: llvs,
	})
}

// Revoke 撤销LifeLog（将指定id的LifeLog状态更换为"隐藏"，对读者不可见）
func (l *LifeLogHandler) Revoke(ctx *gin.Context) {
	// 获取token中存储的用户信息
	userInfo, ok := l.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:Search"))
		return
	}
	// 获取请求参数
	idString := ctx.Param("id")
	// 将获取到的id字符串，转为int64
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("id string转int64失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Revoke"))
		return
	}
	_, err = l.lifeLogServiceClient.Revoke(ctx, &lifelogv1.RevokeRequest{
		LifeLogDomain: &lifelogv1.LifeLogDomain{
			Id: id,
			Author: &lifelogv1.Author{
				UserId: userInfo.Id,
			},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "撤销LifeLog失败",
			Data: "error",
		})
		l.logger.Error("撤销LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Revoke"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "撤销LifeLog成功",
		Data: "success",
	})
}

// Publish 发布LifeLog（只有发布后的LifeLog读者才能看到）
func (l *LifeLogHandler) Publish(ctx *gin.Context) {
	// 有id表示修改LifeLog，无id表示新建LifeLog
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 4,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("参数bind失败，Publish方法", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Publish"))
		return
	}
	// 获取作者Id，作者id也就是登陆的用户id
	userClaims := ctx.MustGet("userClaims")
	userInfo, ok := userClaims.(UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("未发现用户的session信息，Publish方法", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Publish"))
		return
	}
	_, err = l.lifeLogServiceClient.Publish(ctx, &lifelogv1.PublishRequest{
		LifeLogDomain: &lifelogv1.LifeLogDomain{
			Id:      req.Id,
			Title:   req.Title,
			Content: req.Content,
			Author: &lifelogv1.Author{
				UserId: userInfo.Id,
			},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("LifeLog发表失败，Publish方法", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Publish"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "LifeLog发表成功",
		Data: "success",
	})
}

// Detail 查看LifeLog详情（根据id查找LifeLog）
func (l *LifeLogHandler) Detail(ctx *gin.Context) {
	type SearchLifeLogReq struct {
		Id int64 `json:"id"`
		// 操作的是线上库，还是草稿库
		// true 线上库，false 制作库
		Public bool `json:"public"`
	}
	var req SearchLifeLogReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("参数bind失败",
			loggerx.String("method:", "LifeLogHandler:SearchById"))
		return
	}
	if req.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "参数错误",
			Data: "error",
		})
		l.logger.Error("id参数<=0",
			loggerx.String("method:", "LifeLogHandler:SearchById"))
		return
	}
	// 获取token中存储的用户信息
	userInfo, ok := l.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		l.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:Search"))
		return
	}
	res, err := l.lifeLogServiceClient.Detail(ctx, &lifelogv1.DetailRequest{
		LifeLogDomain: &lifelogv1.LifeLogDomain{
			Id: req.Id,
		},
		IsPublic: req.Public,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "查找LifeLog失败",
			Data: "error",
		})
		l.logger.Error("查找LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	// 增加阅读量
	// 异步执行，防止子协程出现错误，导致主协程阻塞
	er := l.syncProducer.ProduceReadEvent(lifeLogEvent.ReadEvent{
		LifeLogId: res.GetLifeLogDomain().GetId(),
		UserId:    userInfo.Id,
	})
	if er != nil {
		l.logger.Error("增加阅读量失败", loggerx.Error(er),
			loggerx.String("method:", "LifeLogHandler:Search"))
	}
	// 获取点赞数，收藏数，阅读数
	interactiveInfo, err := l.interactiveServiceClient.GetInteractiveInfo(ctx, &interactivev1.GetInteractiveInfoRequest{
		InteractiveDomain: &interactivev1.InteractiveDomain{
			Biz:   l.biz,
			BizId: res.GetLifeLogDomain().GetId(),
		},
	})
	if err != nil {
		l.logger.Error("获取点赞数，收藏数，阅读数失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Detail"))
	}
	ctx.JSON(http.StatusOK, Result[vo.LifeLogVo]{
		Code: 200,
		Msg:  "查找LifeLog成功",
		Data: vo.LifeLogVo{
			Id:      res.GetLifeLogDomain().GetId(),
			Title:   res.GetLifeLogDomain().GetTitle(),
			Content: res.GetLifeLogDomain().GetContent(),
			// 将毫秒值时间戳，转换为，time.Time类型(2024-09-23 16:00:00 +0800 CST)
			CreateTime:   time.UnixMilli(res.GetLifeLogDomain().GetCreateTime()),
			UpdateTime:   time.UnixMilli(res.GetLifeLogDomain().GetUpdateTime()),
			AuthorId:     userInfo.Id,
			AuthorName:   userInfo.NickName,
			Status:       uint8(res.GetLifeLogDomain().GetStatus()),
			ReadCount:    interactiveInfo.GetInteractiveDomain().GetReadCount(),
			LikeCount:    interactiveInfo.GetInteractiveDomain().GetLikeCount(),
			CollectCount: interactiveInfo.GetInteractiveDomain().GetCollectCount(),
		},
	})
}

func (l *LifeLogHandler) Hot(c *gin.Context) {
	// 获取热榜
	res, err := l.job.Run()
	// 将res转为[]string
	r, _ := res.([]string)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "获取热门LifeLog失败",
			Data: "error",
		})
		l.logger.Error("获取热门LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Hot"))
		return
	}
	c.JSON(http.StatusOK, Result[[]string]{
		Code: 200,
		Msg:  "获取热门LifeLog成功",
		Data: r,
	})
}

// Abstract LifeLog摘要，取前100字
func (l *LifeLogHandler) Abstract(content string) string {
	// 将字符串转为UTF-8的rune数组
	runes := []rune(content)
	// 判断长度，如果长度小于100，直接返回
	if len(runes) < 100 {
		return content
	}
	// 截取前100个字符
	return string(runes[:100]) + "..."
}
