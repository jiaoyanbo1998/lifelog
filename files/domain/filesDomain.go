package domain

type FileDomain struct {
	Id         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	Url        string `json:"url"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
