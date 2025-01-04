package vo

type CollectVo struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Status     uint8  `json:"status"`
	UserId     int64  `json:"userId"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

type CollectDetailVo struct {
	Id              int64
	CollectId       int64
	LifeLogId       int64
	UpdateTime      int64
	CreateTime      int64
	Status          uint8
	PublicLifeLogVo PublicLifeLogVo
}

type PublicLifeLogVo struct {
	Title    string
	Content  string
	AuthorId int64
}
