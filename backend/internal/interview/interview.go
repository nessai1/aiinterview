package interview

import (
	"context"
	"fmt"
	"github.com/nessai1/aiinterview/internal/ai"
	"github.com/nessai1/aiinterview/internal/domain"
	"github.com/nessai1/aiinterview/internal/message"
	"github.com/nessai1/aiinterview/internal/storage"
	"go.uber.org/zap"
)

type Service struct {
	storage       storage.Storage
	aiService     *ai.Service
	logger        *zap.Logger
	messageParser *message.Parser
}

var ErrSmallTiming = fmt.Errorf("timing is too small: one section >= 5 minutes")

func NewService(str storage.Storage, aiService *ai.Service, logger *zap.Logger, messageParser *message.Parser) (*Service, error) {
	s := Service{
		storage:       str,
		aiService:     aiService,
		logger:        logger,
		messageParser: messageParser,
	}

	return &s, nil
}

func (s *Service) GetUserInterviewList(ctx context.Context, user domain.User) ([]*domain.Interview, error) {
	interviews, err := s.storage.GetUserInterviewList(ctx, user.UUID)
	if err != nil {
		return nil, fmt.Errorf("cannot get user interview list from storage: %w", err)
	}

	for _, interview := range interviews {
		if interview.IsComplete && interview.Feedback == "" {
			// внутри идет закрытие интервью
			_, err = s.GetInterview(ctx, user, interview.UUID)
			if err != nil {
				return nil, fmt.Errorf("cannot close interview while get interview in list load: %w", err)
			}
		}
	}

	return interviews, nil
}

func (s *Service) CreateInterview(ctx context.Context, user domain.User, title string, timingMins int, topics []domain.Topic) (domain.Interview, error) {
	if len(topics)*5 > timingMins {
		return domain.Interview{}, ErrSmallTiming
	}

	thread, firstQuestion, err := s.aiService.Start(ctx, topics, timingMins)
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot start AI: %w", err)
	}

	firstQuestionParsed, err := s.messageParser.Parse([]byte(firstQuestion))
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot parse first question: %w", err)
	}

	firstQuestion = string(firstQuestionParsed)

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

	if interview.IsComplete && interview.Feedback == "" {
		err := s.closeInterview(ctx, &interview, user)
		if err != nil {
			return domain.Interview{}, fmt.Errorf("cannot close interview while get interview: %w", err)
		}
	}

	return interview, nil
}

func (s *Service) closeInterview(ctx context.Context, interview *domain.Interview, user domain.User) error {

	activeSection := interview.GetActiveSection()
	if activeSection == nil {
		return fmt.Errorf("cannot get active section for close interview")
	}

	question := activeSection.GetActiveQuestion()

	if question != nil {
		// Интервью всегда будет закрыто с одним незаконченным вопросом. Потому что заканчивается оно только по таймеру
		_, err := s.answerQuestion(ctx, "", *question, *interview.Thread)
		if err != nil {
			return fmt.Errorf("cannot answer question for close interview: %w", err)
		}
	}

	feedback, err := s.aiService.Feedback(ctx, *interview.Thread)
	if err != nil {
		return fmt.Errorf("cannot get feedback for close interview: %w", err)
	}

	feedbackParsed, err := s.messageParser.Parse([]byte(feedback))
	if err != nil {
		return fmt.Errorf("cannot parse feedback: %w", err)
	}

	feedback = string(feedbackParsed)

	err = s.storage.CompleteInterview(ctx, interview.UUID, user.UUID, feedback)
	if err != nil {
		return fmt.Errorf("cannot complete interview in storage: %w", err)
	}

	return nil
}

// Flow section

var ErrSectionOver = fmt.Errorf("section is over")
var ErrInterviewOver = fmt.Errorf("interview is over")

var ErrAlreadyAnswered = fmt.Errorf("question already answered")

func (s *Service) AnswerQuestion(ctx context.Context, user domain.User, questionUUID string, answer string) (domain.Question, error) {

	question, err := s.storage.GetQuestion(ctx, questionUUID, user.UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot get question from storage: %w", err)
	}

	if question.Done {
		return domain.Question{}, ErrAlreadyAnswered
	}

	interview, err := s.storage.GetInterview(ctx, question.InterviewUUID, user.UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot get interview from storage: %w", err)
	}

	if interview.Thread == nil {
		return domain.Question{}, fmt.Errorf("interview thread is nil")
	}

	question, err = s.answerQuestion(ctx, answer, question, *interview.Thread)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot answer question: %w", err)
	}

	if interview.IsComplete {
		return domain.Question{}, ErrInterviewOver
	}

	questionSection := domain.Section{}
	for _, section := range interview.Sections {
		if section.UUID == question.SectionUUID {
			questionSection = section
			break
		}
	}

	sectionTime := int(interview.Timing) / len(interview.Sections)
	if (len(interview.Sections)-(questionSection.Position+1))*sectionTime >= interview.SecondsLeft {
		err = s.storage.CompleteSection(ctx, questionSection.UUID, user.UUID)
		if err != nil {
			return domain.Question{}, fmt.Errorf("cannot complete section in storage: %w", err)
		}

		return question, ErrSectionOver
	}

	return question, nil
}

func (s *Service) NextQuestion(ctx context.Context, user domain.User, interviewUUID string) (domain.Question, error) {
	interview, err := s.storage.GetInterview(ctx, interviewUUID, user.UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot get interview: %w", err)
	}

	if interview.IsComplete {
		return domain.Question{}, ErrInterviewOver
	}

	activeSection := interview.GetActiveSection()
	if activeSection == nil {
		return domain.Question{}, fmt.Errorf("cannot get active section")
	}

	activeQuestion := activeSection.GetActiveQuestion()
	if activeQuestion != nil {
		return domain.Question{}, fmt.Errorf("active question already exists")
	}

	questionText, err := s.aiService.Next(ctx, *interview.Thread)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot get next question from AI: %w", err)
	}

	question, err := s.storage.AddQuestion(ctx, questionText, activeSection.UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot add question to storage: %w", err)
	}

	return question, nil
}

func (s *Service) NextSectionQuestion(ctx context.Context, user domain.User, interviewUUID string) (domain.Question, error) {
	interview, err := s.storage.GetInterview(ctx, interviewUUID, user.UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot get interview: %w", err)
	}

	if interview.IsComplete {
		return domain.Question{}, ErrInterviewOver
	}

	activeSection := interview.GetActiveSection()
	if activeSection == nil {
		return domain.Question{}, fmt.Errorf("cannot get active section")
	}

	if activeSection.Position >= len(interview.Sections)-1 {
		return domain.Question{}, ErrInterviewOver
	}

	activeQuestion := activeSection.GetActiveQuestion()
	if activeQuestion != nil {
		return domain.Question{}, fmt.Errorf("active question already exists")
	}

	var nextSection *domain.Section
	for _, section := range interview.Sections {
		if section.Position == activeSection.Position+1 {
			nextSection = &section
			break
		}
	}

	if nextSection == nil {
		return domain.Question{}, fmt.Errorf("cannot get next section on position %d", activeSection.Position+1)
	}

	err = s.storage.CompleteSection(ctx, activeSection.UUID, user.UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot complete section in storage: %w", err)
	}

	err = s.storage.StartSection(ctx, nextSection.UUID, user.UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot start section in storage: %w", err)
	}

	// TODO: нужно быть уверенным что следующая секция будет с position+1, хер поймешь что выдаст AI
	questionText, err := s.aiService.Change(ctx, *interview.Thread)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot get next question from AI: %w", err)
	}

	question, err := s.storage.AddQuestion(ctx, questionText, nextSection.UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot add question to storage: %w", err)
	}

	return question, nil
}

func (s *Service) answerQuestion(ctx context.Context, answer string, question domain.Question, thread domain.ChatThread) (domain.Question, error) {
	var feedback string
	var err error

	if answer != "" {
		// TODO: не проверяется что щас отвечается дествительно этот вопрос. Хотя в каждый момент времени должен быть только один вопрос с done=false
		feedback, err = s.aiService.Answer(ctx, thread, answer)
		if err != nil {
			return domain.Question{}, fmt.Errorf("cannot answer question to AI: %w", err)
		}

		answerParsed, err := s.messageParser.Parse([]byte(answer))
		if err != nil {
			return domain.Question{}, fmt.Errorf("cannot parse answer: %w", err)
		}
		answer = string(answerParsed)

	} else {
		feedback, err = s.aiService.Skip(ctx, thread)
		if err != nil {
			return domain.Question{}, fmt.Errorf("cannot skip question to AI: %w", err)
		}
	}

	feedbackParsed, err := s.messageParser.Parse([]byte(feedback))
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot parse feedback: %w", err)
	}

	feedback = string(feedbackParsed)

	// TODO: нужно прокинуть userUUID
	question, err = s.storage.AnswerQuestion(ctx, question.UUID, "", answer, feedback)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot answer question in storage: %w", err)
	}

	return question, err
}
