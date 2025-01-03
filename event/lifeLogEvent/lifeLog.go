package lifeLogEvent

// ReadEvent 文章阅读事件
type ReadEvent struct {
	LifeLogId int64 `json:"life_log_id"`
	UserId    int64 `json:"user_id"`
}

// Consumer 消费者接口
type Consumer interface {
	// Start 开始消费文章阅读事件
	Start() error
}

// Producer 生产者接口
type Producer interface {
	// ProduceReadEvent 生产文章阅读事件
	ProduceReadEvent(evt ReadEvent) error
}
