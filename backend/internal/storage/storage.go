package storage

import (
	"context"
	"fmt"
	"github.com/nessai1/aiinterview/internal/domain"
)

var ErrNotFound = fmt.Errorf("entity not found")

type Storage interface {
	GetUserInterviewList(ctx context.Context, UserUUID string) ([]*domain.Interview, error)
	RegisterUser(ctx context.Context) (domain.User, error)

	GetAssistant(ctx context.Context, ID string) (domain.Assistant, error)
	SetAssistant(ctx context.Context, assistant domain.Assistant) error

	CreateInterview(ctx context.Context, owner domain.User, title string, timing int, topics []domain.Topic) (domain.Interview, error)
}
