package domain

type LifeLogDomain struct {
	Id         int64
	Title      string
	Content    string
	Author     Author
	CreateTime int64
	UpdateTime int64
	Status     uint8
}

type Author struct {
	Id   int64
	Name string
}

type UserDomain struct {
	Id            int64
	Email         string
	Password      string
	CreateTime    int64
	UpdateTime    int64
	Phone         string
	WechatUnionId string
	WechatOpenId  string
	NickName      string
}

const (
	LifeLogStatusPublished = 1 // LifeLog已发布(线上库)
	LifeLogStatusUnPublish = 2 // LifeLog未发布(制作库)
	LifeLogStatusHidden    = 3 // LifeLog隐藏了(不对读者展示)
)
