package domain

type InteractiveDomain struct {
	Id           int64  // 主键
	Biz          string // 业务类型
	BizId        int64  // 业务id（文章id）
	ReadCount    int64  // 阅读数
	CollectCount int64  // 收藏数
	LikeCount    int64  // 点赞数
	CreateTime   int64  // 创建时间
	UpdateTime   int64  // 更新时间
}

type InteractiveReadDomain struct {
	Id         int64  // 主键
	Biz        string // 业务类型
	BizId      int64  // 业务id（文章id）
	UpdateTime int64  // 更新时间
	CreateTime int64  // 创建时间
	UserId     int64  // 用户id
}

type InteractiveCollectDomain struct {
	Id         int64  // 主键
	Biz        string // 业务类型
	BizId      int64  // 业务id（文章id）
	UpdateTime int64  // 更新时间
	CreateTime int64  // 创建时间
	Status     uint8  // 软删除，1收藏，2取消收藏
	UserId     int64  // 用户id
}

type InteractiveLikeDomain struct {
	Id         int64  // 主键
	Biz        string // 业务类型
	BizId      int64  // 业务id（文章id）
	UpdateTime int64  // 更新时间
	CreateTime int64  // 创建时间
	Status     uint8  // 软删除，1点赞，2取消点赞
	UserId     int64  // 用户id
}

type FollowDomain struct {
	Id         int64 // 主键
	FollowerId int64 // 关注着
	FolloweeId int64 // 被关注着
	CreateTime int64 // 创建时间
}
