package storage

import (
	"context"
	"github.com/nessai1/aiinterview/internal/domain"
)

type Storage interface {
	GetUserInterviewList(ctx context.Context, UserUUID string) ([]*domain.Interview, error)
	RegisterUser(ctx context.Context) (domain.User, error)

	GetAssistant(ctx context.Context, ID string) (domain.Assistant, error)
	SetAssistant(ctx context.Context, assistant domain.Assistant) error
}
