package domain

type UserDomain struct {
	Id            int64
	Email         string
	Password      string
	Phone         string
	WechatUnionId string
	WechatOpenId  string
	NickName      string
	NewPassword   string
	Code          string
	UpdateTime    int64
	Authority     int64
	Avatar        string
}
