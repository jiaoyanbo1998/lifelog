package vo

type FindCommentFeedVo struct {
	UserId          int64  `json:"user_id"`
	Biz             string `json:"biz"`
	BizId           int64  `json:"biz_id"`
	CommentedUserId int64  `json:"commented_user_id"`
	Content         string `json:"content"`
}
