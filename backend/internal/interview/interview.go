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

var ErrSmallTiming = fmt.Errorf("timing is too small: one section >= 5 minutes")

func NewService(str storage.Storage, aiService *ai.Service, logger *zap.Logger) (*Service, error) {
	s := Service{
		storage:   str,
		aiService: aiService,
		logger:    logger,
	}

	return &s, nil
}

func (s *Service) CreateInterview(ctx context.Context, user domain.User, title string, timingMins int, topics []domain.Topic) (domain.Interview, error) {
	if len(topics)*5 > timingMins {
		return domain.Interview{}, ErrSmallTiming
	}

	thread, firstQuestion, err := s.aiService.Start(ctx, topics, timingMins)
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot start AI: %w", err)
	}

	interview, err := s.storage.CreateInterview(ctx, user, title, timingMins*60, topics, thread)

	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot create new interview in storage: %w", err)
	}

	question, err := s.storage.AddQuestion(ctx, firstQuestion, interview.Sections[0].UUID)
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot add first question to storage: %w", err)
	}

	interview.Sections[0].Questions = append(interview.Sections[0].Questions, question)

	return interview, nil
}

func (s *Service) GetInterview(ctx context.Context, user domain.User, interviewUUID string) (domain.Interview, error) {

	interview, err := s.storage.GetInterview(ctx, interviewUUID, user.UUID)
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot get interview from storage: %w", err)
	}

	return interview, nil
}

// Flow section

var ErrSectionOver = fmt.Errorf("section is over")
var ErrInterviewOver = fmt.Errorf("interview is over")

var ErrAlreadyAnswered = fmt.Errorf("question already answered")

func (s *Service) AnswerQuestionAndGetNewQuestion(ctx context.Context, user domain.User, questionUUID string, answer string) (domain.Question, error) {
	// TODO
	return domain.Question{}, nil
}
