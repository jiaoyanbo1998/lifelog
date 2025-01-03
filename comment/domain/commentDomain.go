package domain

type CommentDomain struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	Biz        string `json:"biz"`
	BizId      int64  `json:"biz_id"`
	Content    string `json:"content"`
	ParentId   int64  `json:"parent_id"`
	RootId     int64  `json:"root_id"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
