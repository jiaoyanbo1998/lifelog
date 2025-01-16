package vo

type LifeLogVo struct {
	LifeLogListVo
	LifeLogInterVo
}

type LifeLogListVo struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
	AuthorId   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	Status     uint8  `json:"status"`
}

type LifeLogInterVo struct {
	LikeCount    int64 `json:"like_count"`
	CollectCount int64 `json:"collect_count"`
	ReadCount    int64 `json:"read_count"`
}
