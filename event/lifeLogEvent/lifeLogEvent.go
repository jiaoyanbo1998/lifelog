package lifeLogEvent

type EventLifeLog struct {
	AuthorId int64 `json:"author_id"`
	Limit    int   `json:"limit"`
	Offset   int   `json:"offset"`
}
