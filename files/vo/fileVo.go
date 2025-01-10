package vo

type FileVo struct {
	Name       string `json:"name"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
}
