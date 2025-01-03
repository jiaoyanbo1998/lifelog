package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"lifelog-grpc/event/lifeLogEvent"
	iService "lifelog-grpc/interactive/service"
	"lifelog-grpc/job"
	aDomain "lifelog-grpc/lifeLog/domain"
	aService "lifelog-grpc/lifeLog/service"
	"lifelog-grpc/lifeLog/vo"
	"lifelog-grpc/pkg/loggerx"
	"net/http"
	"strconv"
	"time"
)

type LifeLogHandler struct {
	logger             loggerx.Logger
	lifeLogService     aService.LifeLogService
	interactiveService iService.InteractiveService
	job                job.Job
	biz                string
	syncProducer       lifeLogEvent.Producer
	JWTHandler
}

func NewLifeLogHandler(l loggerx.Logger, lifeLogService aService.LifeLogService,
	interactiveService iService.InteractiveService, job job.Job,
	syncProducer lifeLogEvent.Producer) *LifeLogHandler {
	return &LifeLogHandler{
		logger:             l,
		lifeLogService:     lifeLogService,
		interactiveService: interactiveService,
		biz:                "lifeLog",
		job:                job,
		syncProducer:       syncProducer,
	}
}

func (a *LifeLogHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/lifeLog")
	// 1.创作者的接口
	// 编辑LifeLog（新建LifeLog/修改LifeLog）
	rg.POST("/edit", a.Edit)
	// 删除LifeLog（根据Id删除LifeLog，可以删除线上库/制作库中的LifeLog）
	rg.POST("/delete", a.Delete)
	// 撤销LifeLog（将指定id的LifeLog状态更换为"隐藏"，对读者不可见）
	rg.PUT("/revoke/:id", a.Revoke)
	// 发布LifeLog（只有发布后的LifeLog读者才能看到）
	rg.POST("/publish", a.Publish)
	// 作者的LifeLog列表（根据作者id分页查询制作库中的LifeLog）
	rg.POST("/author_id", a.DraftList)
	// 2.所有用户的接口（创作者+读者）
	// 查看LifeLog的详情 （根据id查找LifeLog，线上库/制作库）
	rg.POST("/detail", a.Detail)
	// 3.读者的接口
	// 根据标题查找LifeLog
	rg.POST("/title", a.SearchByTitle)
	// 点赞LifeLog
	rg.PUT("/like/:id", a.Like)
	// 取消点赞LifeLog
	rg.PUT("/unlike/:id", a.UnLike)
	// 收藏LifeLog
	rg.POST("/collect", a.Collect)
	// 取消收藏LifeLog
	rg.PUT("/unCollect/:id", a.UnCollect)
	// 热榜
	rg.GET("/hot", a.Hot)
}

// Edit 编辑LifeLog（新建LifeLog/修改LifeLog）
func (a *LifeLogHandler) Edit(ctx *gin.Context) {
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
		a.logger.Error("参数bind失败",
			loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:Edit"))
		return
	}
	// 获取token中存储的用户信息
	userInfo, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:Edit"))
		return
	}
	// 调用service层代码
	// 传入id，表示是修改文化在哪个
	// 不传入id，表示创建新LifeLog
	lifeLog, err := a.lifeLogService.Save(ctx.Request.Context(), aDomain.LifeLogDomain{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: aDomain.Author{
			Id: userInfo.Id,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("编辑LifeLog失败",
			loggerx.Error(err),
			loggerx.String("method：", "LifeLogHandler:Edit"))
		return
	}
	ctx.JSON(http.StatusOK, Result[vo.LifeLogVo]{
		Code: 200,
		Msg:  "编辑LifeLog成功",
		Data: vo.LifeLogVo{
			Id:      lifeLog.Id,
			Content: lifeLog.Content,
			// 将毫秒值时间戳，转换为，time.Time类型(2024-09-23 16:00:00 +0800 CST)
			CreateTime: time.UnixMilli(lifeLog.CreateTime),
			Title:      lifeLog.Title,
			UpdateTime: time.UnixMilli(lifeLog.UpdateTime),
			AuthorId:   userInfo.Id,
			AuthorName: userInfo.NickName,
		},
	})
}

// Delete 删除LifeLog（根据Id删除LifeLog，可以删除线上库/制作库中的LifeLog）
func (a *LifeLogHandler) Delete(ctx *gin.Context) {
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
		a.logger.Error("参数bind失败",
			loggerx.String("method:", "LifeLogHandler:Delete"))
		return
	}
	// 调用service层代码
	err = a.lifeLogService.Delete(ctx.Request.Context(), req.Ids, req.Public)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "删除LifeLog失败",
			Data: "error",
		})
		a.logger.Error("删除LifeLog失败", loggerx.Error(err),
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
func (a *LifeLogHandler) SearchByTitle(ctx *gin.Context) {
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
		a.logger.Error("参数bind失败",
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	// 调用service层代码
	lifeLogs, err := a.lifeLogService.SearchByTitle(ctx.Request.Context(),
		req.Title, req.Limit, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "查找LifeLog失败",
			Data: "error",
		})
		a.logger.Error("查找LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.LifeLogVo]{
		Code: 200,
		Msg:  "查找LifeLog成功",
		Data: a.lifeLogDomainsToLifeLogVos(lifeLogs),
	})
}

// Detail 查看LifeLog详情（根据id查找LifeLog）
func (a *LifeLogHandler) Detail(ctx *gin.Context) {
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
		a.logger.Error("参数bind失败",
			loggerx.String("method:", "LifeLogHandler:SearchById"))
		return
	}
	if req.Id <= 0 {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "参数错误",
			Data: "error",
		})
		a.logger.Error("id参数<=0",
			loggerx.String("method:", "LifeLogHandler:SearchById"))
		return
	}
	// 获取token中存储的用户信息
	userInfo, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:Search"))
		return
	}
	// 调用service层代码
	lifeLog, err := a.lifeLogService.Detail(ctx.Request.Context(), req.Id, req.Public)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "查找LifeLog失败",
			Data: "error",
		})
		a.logger.Error("查找LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	// 增加阅读量
	// 异步执行，防止子协程出现错误，导致主协程阻塞
	go func() {
		// 为什么要新建一个ctx
		//    防止因为主协程执行结束，导致ctx被取消，从而异步执行也会被取消
		ctxNew, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := a.interactiveService.IncreaseReadCount(ctxNew, a.biz, lifeLog.Id, userInfo.Id)
		if err != nil {
			a.logger.Error("增加阅读量失败", loggerx.Error(er),
				loggerx.String("method:", "LifeLogHandler:Search"))
		}
	}()
	// 获取点赞数，收藏数，阅读数
	res, err := a.interactiveService.GetInteractiveInfo(ctx, a.biz, lifeLog.Id)
	if err != nil {
		a.logger.Error("获取点赞数，收藏数，阅读数失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Detail"))
	}
	ctx.JSON(http.StatusOK, Result[vo.LifeLogVo]{
		Code: 200,
		Msg:  "查找LifeLog成功",
		Data: vo.LifeLogVo{
			Id:      lifeLog.Id,
			Title:   lifeLog.Title,
			Content: lifeLog.Content,
			// 将毫秒值时间戳，转换为，time.Time类型(2024-09-23 16:00:00 +0800 CST)
			CreateTime:   time.UnixMilli(lifeLog.CreateTime),
			UpdateTime:   time.UnixMilli(lifeLog.UpdateTime),
			AuthorId:     userInfo.Id,
			AuthorName:   userInfo.NickName,
			Status:       lifeLog.Status,
			ReadCount:    res.ReadCount,
			LikeCount:    res.LikeCount,
			CollectCount: res.CollectCount,
		},
	})
}

// DraftList 作者的LifeLog列表（根据作者id分页查询制作库中的LifeLog）
func (a *LifeLogHandler) DraftList(ctx *gin.Context) {
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
		a.logger.Error("参数bind失败",
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	// 获取token中存储的用户信息
	userInfo, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:Search"))
		return
	}
	// 调用service层代码
	lifeLogs, err := a.lifeLogService.SearchByAuthorId(ctx.Request.Context(), userInfo.Id,
		req.Limit, req.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "查找LifeLog列表失败",
			Data: "error",
		})
		a.logger.Error("查找LifeLog列表失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Search"))
		return
	}
	ctx.JSON(http.StatusOK, Result[[]vo.LifeLogVo]{
		Code: 200,
		Msg:  "查找LifeLog列表成功",
		Data: a.lifeLogDomainsToLifeLogVos(lifeLogs),
	})
}

// Revoke 撤销LifeLog（将指定id的LifeLog状态更换为"隐藏"，对读者不可见）
func (a *LifeLogHandler) Revoke(ctx *gin.Context) {
	// 获取token中存储的用户信息
	userInfo, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
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
		a.logger.Error("id string转int64失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Revoke"))
		return
	}
	err = a.lifeLogService.Revoke(ctx.Request.Context(), id, userInfo.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "撤销LifeLog失败",
			Data: "error",
		})
		a.logger.Error("撤销LifeLog失败", loggerx.Error(err),
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
func (a *LifeLogHandler) Publish(ctx *gin.Context) {
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
		a.logger.Error("参数bind失败，Publish方法", loggerx.Error(err),
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
		a.logger.Error("未发现用户的session信息，Publish方法", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Publish"))
		return
	}
	// 调用service层代码
	err = a.lifeLogService.Publish(ctx, aDomain.LifeLogDomain{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: aDomain.Author{
			Id: userInfo.Id,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result[string]{
			Code: 5,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("LifeLog发表失败，Publish方法", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Publish"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "LifeLog发表成功",
		Data: "success",
	})
}

// lifeLogDomainsToLifeLogVos 将domain.LifeLogDomain转换为vo.LifeLogVo
func (a *LifeLogHandler) lifeLogDomainsToLifeLogVos(ads []aDomain.LifeLogDomain) []vo.LifeLogVo {
	avs := make([]vo.LifeLogVo, 0, len(ads))
	for _, ad := range ads {
		avs = append(avs, vo.LifeLogVo{
			Id:      ad.Id,
			Content: a.Abstract(ad.Content),
			Title:   ad.Title,
			// 将毫秒值时间戳，转换为，time.Time类型(2024-09-23 16:00:00 +0800 CST)
			UpdateTime: time.UnixMilli(ad.UpdateTime),
			CreateTime: time.UnixMilli(ad.CreateTime),
			AuthorId:   ad.Author.Id,
			AuthorName: ad.Author.Name,
			Status:     ad.Status,
		})
	}
	return avs
}

// Abstract LifeLog摘要，取前100字
func (a *LifeLogHandler) Abstract(content string) string {
	// 将字符串转为UTF-8的rune数组
	runes := []rune(content)
	// 判断长度，如果长度小于100，直接返回
	if len(runes) < 100 {
		return content
	}
	// 截取前100个字符
	return string(runes[:100]) + "..."
}

// Like 点赞LifeLog
func (a *LifeLogHandler) Like(ctx *gin.Context) {
	// 获取路径参数
	idString := ctx.Param("id")
	// 将string转为int64
	lifeLogId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("id string转int64失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Like"))
		return
	}
	userInfo, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:Like"))
		return
	}
	// 不启动生产者
	/*
		err = a.interactiveService.IncreaseLikeCount(ctx.Request.Context(),
			a.biz, lifeLogId, userInfo.Id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, Result[string]{
				Code: 500,
				Msg:  "点赞失败",
				Data: "error",
			})
			a.logger.Error("点赞失败", loggerx.Error(err),
				loggerx.String("method:", "LifeLogHandler:Like"))
			return
		}
	*/
	// 启动生产者
	err = a.syncProducer.ProduceReadEvent(lifeLogEvent.ReadEvent{
		LifeLogId: lifeLogId,
		UserId:    userInfo.Id,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "点赞失败",
			Data: "error",
		})
		a.logger.Error("点赞失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Like"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "点赞成功",
		Data: "success",
	})
}

// UnLike 取消点赞LifeLog
func (a *LifeLogHandler) UnLike(ctx *gin.Context) {
	// 获取路径参数
	idString := ctx.Param("id")
	// 将string转为int64
	lifeLogId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("id string转int64失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:UnLike"))
		return
	}
	userInfo, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:UnLike"))
		return
	}
	// 减少点赞数
	err = a.interactiveService.DecreaseLikeCount(ctx.Request.Context(),
		a.biz, lifeLogId, userInfo.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "取消点赞失败",
			Data: "error",
		})
		a.logger.Error("取消点赞失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:UnLike"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "取消点赞成功",
		Data: "success",
	})
}

// Collect 收藏LifeLog
func (a *LifeLogHandler) Collect(ctx *gin.Context) {
	type CollectReq struct {
		BidId int64 `json:"biz_id"`
	}
	var req CollectReq
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("参数bind失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:collect"))
		return
	}
	userInfo, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:collect"))
		return
	}
	err = a.interactiveService.IncreaseCollectCount(ctx.Request.Context(),
		a.biz, req.BidId, userInfo.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "收藏失败",
			Data: "error",
		})
		a.logger.Error("收藏失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:collect"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "收藏成功",
		Data: "success",
	})
}

// UnCollect 取消收藏LifeLog
func (a *LifeLogHandler) UnCollect(ctx *gin.Context) {
	// 获取路径参数
	idString := ctx.Param("id")
	// 将string转为int64
	lifeLogId, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result[string]{
			Code: 400,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("id string转int64失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:UnCollect"))
		return
	}
	userInfo, ok := a.GetUserInfo(ctx)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "系统错误",
			Data: "error",
		})
		a.logger.Error("获取用户信息失败，token中不存在用户信息",
			loggerx.String("method：", "LifeLogHandler:UnCollect"))
		return
	}
	// 减少收藏数
	err = a.interactiveService.DecreaseCollectCount(ctx.Request.Context(),
		a.biz, lifeLogId, userInfo.Id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "取消收藏失败",
			Data: "error",
		})
		a.logger.Error("取消收藏失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:UnCollect"))
		return
	}
	ctx.JSON(http.StatusOK, Result[string]{
		Code: 200,
		Msg:  "取消收藏成功",
		Data: "success",
	})
}

func (a *LifeLogHandler) Hot(c *gin.Context) {
	// 获取热榜
	res, err := a.job.Run()
	// 将res转为[]string
	r, _ := res.([]string)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Result[string]{
			Code: 500,
			Msg:  "获取热门LifeLog失败",
			Data: "error",
		})
		a.logger.Error("获取热门LifeLog失败", loggerx.Error(err),
			loggerx.String("method:", "LifeLogHandler:Hot"))
		return
	}
	c.JSON(http.StatusOK, Result[[]string]{
		Code: 200,
		Msg:  "获取热门LifeLog成功",
		Data: r,
	})
}
