package domain

type FeedEvent struct {
	ID         int64  `json:"id"`
	UserId     int64  `json:"user_id"`
	Content    string `json:"content"`     // 事件的内容(业务内容)
	Type       string `json:"type"`        // 事件的类型，根据不同的事件类型，调用不同的handler
	CreateTime int64  `json:"create_time"` // 创建时间
}

type LifeLogCommentEvent struct {
	Biz             string `json:"biz"`
	BizId           int64  `json:"biz_id"`
	CommentedUserId int64  `json:"commented_user_id"`
	Content         string `json:"content"`
}

type LikeFeedEvent struct {
	Biz         string `json:"biz"`
	BizId       int64  `json:"biz_id"`
	LikedUserId int64  `json:"liked_user_id"`
}

type FollowFeedEvent struct {
	Biz            string `json:"biz"`
	BizId          int64  `json:"biz_id"`
	FollowedUserId int64  `json:"followed_user_id"`
}
