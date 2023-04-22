package event

// All event names must be lowercase and follow the structure: "{noun}-{action}".
// Eg. article-created, notification-sent, order-accepted.
// For longer names use snake-case naming.
// Eg. changed_password_notification-sent.
type EventType string

const (
	ArticleCreated EventType = "article-created"
	ArticleDeleted EventType = "article-deleted"
	ArticleUpdated EventType = "article-updated"

	UserCreated EventType = "user-created"
	UserDeleted EventType = "user-deleted"
	UserUpdated EventType = "user-updated"

	UserLoggedIn  EventType = "user-logged_in"
	UserLoggedOut EventType = "user-logged_out"

	KeySetUpdated EventType = "key_set-updated"
)
