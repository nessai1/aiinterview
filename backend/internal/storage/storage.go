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

	CreateInterview(ctx context.Context, owner domain.User, title string, timing int, topics []domain.Topic, thread domain.ChatThread) (domain.Interview, error)

	GetQuestion(ctx context.Context, UUID string, userUUID string) (domain.Question, error)
	AddQuestion(ctx context.Context, question, sectionUUID string) (domain.Question, error)
	AnswerQuestion(ctx context.Context, UUID, userUUID string, answer, feedback string) (domain.Question, error)
	GetSection(ctx context.Context, UUID string) (domain.Section, error)
	GetInterview(ctx context.Context, UUID string, UserUUID string) (domain.Interview, error)
	CompleteSection(ctx context.Context, UUID string, userUUID string) error
	StartSection(ctx context.Context, UUID string, userUUID string) error
	CompleteInterview(ctx context.Context, UUID string, userUUID string, feedback string) error
	IsUserExists(ctx context.Context, userUUID string) (bool, error)
}
