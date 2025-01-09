package domain

type FileDomain struct {
	Biz        string `json:"biz"`
	BizId      int64  `json:"biz_id"`
	Url        string `json:"url"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
