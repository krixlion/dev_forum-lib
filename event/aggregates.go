package event

type AggregateId string

const (
	AuthAggregate    AggregateId = "auth"
	UserAggregate    AggregateId = "user"
	ArticleAggregate AggregateId = "article"
	CommentAggregate AggregateId = "comment"
)
