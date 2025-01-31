package event

type EventLifeLog struct {
	AuthorId int64 `json:"author_id"`
	Limit    int64 `json:"limit"`
	Offset   int64 `json:"offset"`
}
