package vo

import "time"

type LifeLogVo struct {
	Id           int64     `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
	AuthorId     int64     `json:"author_id"`
	AuthorName   string    `json:"author_name"`
	Status       uint8     `json:"status"`
	LikeCount    int64     `json:"like_count"`
	ReadCount    int64     `json:"read_count"`
	CollectCount int64     `json:"collect_count"`
}
