package domain

type UserEventType string

const (
	UserCreatedEvent UserEventType = "user_created"
	UserUpdatedEvent UserEventType = "user_updated"
	UserDeletedEvent UserEventType = "user_deleted"
)
