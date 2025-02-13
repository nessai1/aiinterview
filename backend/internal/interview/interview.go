package interview

import (
	"context"
	"fmt"
	"github.com/nessai1/aiinterview/internal/ai"
	"github.com/nessai1/aiinterview/internal/domain"
	"github.com/nessai1/aiinterview/internal/storage"
	"go.uber.org/zap"
)

type Service struct {
	storage   storage.Storage
	aiService *ai.Service
	logger    *zap.Logger
}

const ErrSmallTiming = fmt.Errorf("timing is too small: one section >= 5 minutes")

func NewService(str storage.Storage, aiService *ai.Service, logger *zap.Logger) (*Service, error) {
	s := Service{
		storage:   str,
		aiService: aiService,
		logger:    logger,
	}

	return &s, nil
}

func (s *Service) CreateInterview(ctx context.Context, user domain.User, title string, timing int, topics []domain.Topic) (domain.Interview, error) {
	mins := timing / 60
	if len(topics)*5 < mins {
		return domain.Interview{}, ErrSmallTiming
	}

	interview, err := s.storage.CreateInterview(ctx, user, title, timing, topics)
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot create new interview in storage: %w", err)
	}

	return interview, nil
}
