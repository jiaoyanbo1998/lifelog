package vo

type CollectVo struct {
	Id         int64  `json:"Id"`
	Name       string `json:"name"`
	Status     uint8  `json:"status"`
	UserId     int64  `json:"userId"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

type CollectDetailVo struct {
	CollectId       int64
	UpdateTime      int64
	CreateTime      int64
	Status          uint8
	PublicLifeLogVo []PublicLifeLogVo
}

type PublicLifeLogVo struct {
	PublicLifeLogId int64
	Title           string
	Content         string
	AuthorId        int64
}
