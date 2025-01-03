package commentEvent

import "database/sql"

// CommentEvent ”评论插入数据库“事件
type CommentEvent struct {
	UserId   int64         `json:"user_id"`
	BizId    int64         `json:"biz_id"` // 文章id
	Content  string        `json:"content"`
	ParentId sql.NullInt64 `json:"parent_id"`
	RootId   sql.NullInt64 `json:"root_id"`
}

// Consumer 消费者接口
type Consumer interface {
	// Start 开始消费文章阅读事件
	Start() error
}

// Producer 生产者接口
type Producer interface {
	// ProduceCommentEvent 生产”评论插入数据库“事件
	ProduceCommentEvent(evt CommentEvent) error
}
