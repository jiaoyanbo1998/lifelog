package vo

type FindCommentFeedVo struct {
	UserId          int64  `json:"user_id"`
	Biz             string `json:"biz"`
	BizId           int64  `json:"biz_id"`
	CommentedUserId int64  `json:"commented_user_id"`
	Content         string `json:"content"`
}

type FindLikeFeedVo struct {
	UserId      int64  `json:"user_id"`
	Biz         string `json:"biz"`
	BizId       int64  `json:"biz_id"`
	LikedUserId int64  `json:"liked_user_id"`
}

type FindFollowFeedVo struct {
	UserId         int64  `json:"user_id"`
	Biz            string `json:"biz"`
	FolloweeUserId int64  `json:"followee_user_id"` // 被关注的用户id
	FollowerUserId int64  `json:"follower_user_id"`
}

type FindCollectFeedVo struct {
	UserId          int64  `json:"user_id"`
	Biz             string `json:"biz"`
	BizId           int64  `json:"biz_id"`
	CollectedUserId int64  `json:"collected_user_id"`
}

type FindReadFeedVo struct {
	UserId       int64  `json:"user_id"`
	Biz          string `json:"biz"`
	BizId        int64  `json:"biz_id"`
	ReadedUserId int64  `json:"readed_user_id"`
}
