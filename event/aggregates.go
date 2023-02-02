package event

type AggregateId string

const (
	UserAggregate    AggregateId = "user"
	ArticleAggregate AggregateId = "article"
	CommentAggregate AggregateId = "comment"
)
