package storage

import (
	"context"
	"github.com/nessai1/aiinterview/internal/domain"
)

type Storage interface {
	GetUserInterviewList(ctx context.Context, UserUUID string) ([]domain.Interview, error)
	RegisterUser(ctx context.Context, UserUUID string) error
}
