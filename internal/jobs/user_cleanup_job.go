package jobs

import (
	"log"

	"github.com/JIeeiroSst/hex/internal/core/ports"
)

type UserCleanupJob struct {
	userService ports.UserService
}

func NewUserCleanupJob(userService ports.UserService) *UserCleanupJob {
	return &UserCleanupJob{
		userService: userService,
	}
}

func (j *UserCleanupJob) Execute() {
	log.Println("Starting user cleanup job")
	// Add your cleanup logic here
	// For example, delete inactive users
	j.cleanupInactiveUsers()
	log.Println("Finished user cleanup job")
}

func (j *UserCleanupJob) cleanupInactiveUsers() {
	// Implementation of cleanup logic
}
