package messaging

import (
	"encoding/json"
	"fmt"

	"github.com/JIeeiroSst/hex/internal/core/domain"
	"github.com/JIeeiroSst/hex/internal/core/ports"
)

type UserEvent struct {
	Type string      `json:"type"`
	Data domain.User `json:"data"`
}

type UserEventHandler struct {
	userService ports.UserService
}

func NewUserEventHandler(userService ports.UserService) *UserEventHandler {
	return &UserEventHandler{
		userService: userService,
	}
}

func (h *UserEventHandler) HandleMessage(message []byte) error {
	var event UserEvent
	if err := json.Unmarshal(message, &event); err != nil {
		return err
	}

	switch event.Type {
	case "user_created":
		return h.handleUserCreated(event.Data)
	case "user_updated":
		return h.handleUserUpdated(event.Data)
	case "user_deleted":
		return h.handleUserDeleted(event.Data.ID)
	default:
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

func (h *UserEventHandler) handleUserCreated(user domain.User) error {
	return h.userService.CreateUser(&user)
}

func (h *UserEventHandler) handleUserUpdated(user domain.User) error {
	return h.userService.UpdateUser(&user)
}

func (h *UserEventHandler) handleUserDeleted(userID string) error {
	return h.userService.DeleteUser(userID)
}
