package domain

type CollectDomain struct {
	Id         int64
	Name       string
	CreateTime int64
	UpdateTime int64
	Status     uint8
	AuthorId   int64
}

type CollectDetailDomain struct {
	Id         int64
	CollectId  int64
	LifeLogId  int64
	UpdateTime int64
	CreateTime int64
	Status     uint8
	AuthorId   int64
}
