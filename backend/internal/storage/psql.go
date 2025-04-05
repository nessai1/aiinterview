package storage

import (
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nessai1/aiinterview/internal/utils"
	"math/rand"
	"slices"
	"strings"
	"time"

	"context"
	"database/sql"
	"fmt"
	"github.com/nessai1/aiinterview/internal/domain"
	"github.com/pressly/goose"
)

type PSQLStorage struct {
	db *sql.DB
}

func NewPSQLStorageFromAddr(addr string) (*PSQLStorage, error) {
	conn, err := sql.Open("pgx", addr)
	if err != nil {
		return nil, fmt.Errorf("cannot create DB connection: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping to created DB connection: %w", err)
	}

	s, err := NewPSQLStorage(conn)
	if err != nil {
		return nil, fmt.Errorf("cannot create PSQL storage with created conn: %w", err)
	}

	return s, nil
}

func NewPSQLStorage(db *sql.DB) (*PSQLStorage, error) {
	err := goose.Up(db, "migrations/psql")
	if err != nil {
		return nil, fmt.Errorf("cannot make migrations for PSQL storage: %w", err)
	}

	s := PSQLStorage{
		db: db,
	}

	return &s, nil
}

func (s *PSQLStorage) GetUserInterviewList(ctx context.Context, userUUID string) ([]*domain.Interview, error) {
	req := `SELECT 
    	i.uuid,
    	i.title,
    	i.start_timestamp,
    	i.timing,
    	i.feedback,
    	s.name,
    	s.grade
	FROM interview i LEFT JOIN section s ON i.uuid = s.interview_uuid WHERE i.owner_uuid = $1`

	rows, err := s.db.QueryContext(ctx, req, userUUID)
	if err != nil {
		return nil, fmt.Errorf("error while query interviews rows: %w", err)
	}

	defer rows.Close()

	var feedback sql.NullString
	var uuid, title, sectionName, sectionGrade string
	var timing int
	var startTimestamp time.Time

	interviews := make(map[string]*domain.Interview)
	for rows.Next() {
		err = rows.Scan(&uuid, &title, &startTimestamp, &timing, &feedback, &sectionName, &sectionGrade)
		if err != nil {
			return nil, fmt.Errorf("cannot scan fields of interview: %w", err)
		}

		_, found := interviews[uuid]
		if found {
			interviews[uuid].Sections = append(interviews[uuid].Sections, domain.Section{Name: sectionName, Grade: domain.TopicGrade(sectionGrade)})
		} else {

			feedbackStr := ""
			if feedback.Valid {
				feedbackStr = feedback.String
			}

			sections := []domain.Section{{Name: sectionName, Grade: domain.TopicGrade(sectionGrade)}}
			timingDuration := time.Duration(timing)
			secsLeft := int(startTimestamp.Add(time.Duration(int(time.Second) * timing)).Sub(time.Now().UTC()).Seconds())
			interviews[uuid] = &domain.Interview{
				UUID:           uuid,
				Title:          title,
				Timing:         timingDuration,
				StartTimestamp: startTimestamp,
				Sections:       sections,
				Feedback:       feedbackStr,
				SecondsLeft:    secsLeft,
				IsComplete:     secsLeft <= 0 || feedbackStr != "",
			}
		}
	}

	interviewsSlice := make([]*domain.Interview, 0, len(interviews))
	for _, v := range interviews {
		interviewsSlice = append(interviewsSlice, v)
	}

	return interviewsSlice, nil
}

func (s *PSQLStorage) RegisterUser(ctx context.Context) (domain.User, error) {
	userUUID, err := utils.GenerateUUIDv7()
	if err != nil {
		return domain.User{}, fmt.Errorf("cannot generate UUIDv7: %w", err)
	}

	_, err = s.db.ExecContext(ctx, "INSERT INTO users (uuid) VALUES ($1)", userUUID)

	if err != nil {
		return domain.User{}, fmt.Errorf("error while exec register query: %w", err)
	}

	return domain.User{UUID: userUUID}, nil
}

func (s *PSQLStorage) GetAssistant(ctx context.Context, ID string) (domain.Assistant, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, external_id, model FROM assistants WHERE id = $1", ID)
	if err != nil {
		return domain.Assistant{}, fmt.Errorf("error while query assistant: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var assistant domain.Assistant
		err = rows.Scan(&assistant.ID, &assistant.ExternalID, &assistant.Model)
		if err != nil {
			return domain.Assistant{}, fmt.Errorf("error while scan assistant: %w", err)
		}

		return assistant, nil
	}

	return domain.Assistant{}, ErrNotFound
}

func (s *PSQLStorage) SetAssistant(ctx context.Context, assistant domain.Assistant) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO assistants (id, external_id, model) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET external_id = $2, model = $3", assistant.ID, assistant.ExternalID, assistant.Model)
	if err != nil {
		return fmt.Errorf("error while exec set assistant query: %w", err)
	}

	return nil
}

func (s *PSQLStorage) CreateInterview(ctx context.Context, owner domain.User, title string, timing int, topics []domain.Topic, thread domain.ChatThread) (domain.Interview, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot begin transaction: %w", err)
	}

	interviewUUID, err := utils.GenerateUUIDv7()
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot generate UUIDv7: %w", err)
	}

	startTime := time.Now()
	_, err = tx.ExecContext(ctx, "INSERT INTO interview (uuid, owner_uuid, title, start_timestamp, timing, thread) VALUES ($1, $2, $3, $4, $5, $6)", interviewUUID, owner.UUID, title, startTime.UTC(), timing, thread.ID+"||"+thread.Secret)
	if err != nil {
		tx.Rollback()
		return domain.Interview{}, fmt.Errorf("error while exec insert interview query: %w", err)
	}

	sections := make([]domain.Section, 0, len(topics))
	for i, topic := range topics {
		sectionUUID, err := utils.GenerateUUIDv7()
		if err != nil {
			tx.Rollback()
			return domain.Interview{}, fmt.Errorf("cannot generate UUIDv7 for section: %w", err)
		}

		isStarted := false
		if i == 0 {
			isStarted = true
		}

		color := utils.GenerateColor(rand.Int())

		_, err = tx.ExecContext(ctx, "INSERT INTO section (uuid, name, grade, position, interview_uuid, color, is_started) VALUES ($1, $2, $3, $4, $5, $6, $7)", sectionUUID, topic.Name, topic.Grade, i, interviewUUID, color, isStarted)

		if err != nil {
			tx.Rollback()
			return domain.Interview{}, fmt.Errorf("error while exec insert section query: %w", err)
		}

		sections = append(sections, domain.Section{
			UUID:          sectionUUID,
			InterviewUUID: interviewUUID,
			Name:          topic.Name,
			Grade:         topic.Grade,
			Position:      i,
			IsStarted:     false,
			IsComplete:    false,
			Questions:     make([]domain.Question, 0),
			Color:         color,
		})
	}

	err = tx.Commit()
	if err != nil {
		return domain.Interview{}, fmt.Errorf("cannot commit transaction: %w", err)
	}

	return domain.Interview{
		UUID:           interviewUUID,
		Title:          title,
		Timing:         time.Duration(timing),
		StartTimestamp: startTime,
		IsComplete:     false,
		Sections:       sections,
	}, nil
}

// TODO: не хватает времени делать длинный JOIN до владельца, поэтому каждый может отвечать на вопрос каждого :)
func (s *PSQLStorage) GetQuestion(ctx context.Context, UUID string, _ string) (domain.Question, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT q.section_uuid, s.interview_uuid, q.question, q.answer, q.feedback, q.done, q.position FROM question q LEFT JOIN section s ON q.section_uuid = s.uuid WHERE q.uuid = $1", UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("error while query question: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var question domain.Question
		err = rows.Scan(&question.SectionUUID, &question.InterviewUUID, &question.Question, &question.Answer, &question.Feedback, &question.Done, &question.Position)
		if err != nil {
			return domain.Question{}, fmt.Errorf("error while scan question: %w", err)
		}

		question.UUID = UUID
		return question, nil
	}

	return domain.Question{}, ErrNotFound
}

func (s *PSQLStorage) AnswerQuestion(ctx context.Context, UUID, _ string, answer, feedback string) (domain.Question, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE question SET answer = $1, feedback = $2, done = true WHERE uuid = $3", answer, feedback, UUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("error while exec update question query: %w", err)
	}

	return s.GetQuestion(ctx, UUID, "")
}

func (s *PSQLStorage) getSectionQuestions(ctx context.Context, sectionUUID string) ([]domain.Question, error) {
	if sectionUUID == "" {
		return nil, errors.New("cannot get section questions: section UUID is empty")
	}

	rows, err := s.db.QueryContext(ctx, "SELECT question, answer, feedback, uuid, position, done FROM question WHERE section_uuid = $1", sectionUUID)
	if err != nil {
		return nil, fmt.Errorf("error while query questions: %w", err)
	}

	defer rows.Close()

	questions := make([]domain.Question, 0)
	for rows.Next() {
		var question domain.Question
		err = rows.Scan(&question.Question, &question.Answer, &question.Feedback, &question.UUID, &question.Position, &question.Done)
		if err != nil {
			return nil, fmt.Errorf("error while scan question: %w", err)
		}

		question.SectionUUID = sectionUUID
		questions = append(questions, question)
	}

	slices.SortFunc(questions, func(a, b domain.Question) int {
		if a.Position < b.Position {
			return -1
		} else if a.Position > b.Position {
			return 1
		}

		return 0
	})

	return questions, nil
}

func (s *PSQLStorage) AddQuestion(ctx context.Context, question, sectionUUID string) (domain.Question, error) {
	questionUUID, err := utils.GenerateUUIDv7()
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot generate UUIDv7: %w", err)
	}

	questions, err := s.getSectionQuestions(ctx, sectionUUID)
	if err != nil {
		return domain.Question{}, fmt.Errorf("cannot get questions for compute position: %w", err)
	}

	position := len(questions) + 1

	_, err = s.db.ExecContext(ctx, "INSERT INTO question (uuid, section_uuid, question, answer, feedback, done, position) VALUES ($1, $2, $3, $4, $5, $6, $7)", questionUUID, sectionUUID, question, "", "", false, position)
	if err != nil {
		return domain.Question{}, fmt.Errorf("error while exec insert question query: %w", err)
	}

	return domain.Question{
		UUID:        questionUUID,
		SectionUUID: sectionUUID,
		Question:    question,
		Answer:      "",
		Feedback:    "",
		Done:        false,
	}, nil
}

func (s *PSQLStorage) GetSection(ctx context.Context, UUID string) (domain.Section, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT name, grade, position, is_started, is_complete, color FROM section WHERE uuid = $1", UUID)
	if err != nil {
		return domain.Section{}, fmt.Errorf("error while query section: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var section domain.Section
		err = rows.Scan(&section.Name, &section.Grade, &section.Position, &section.IsStarted, &section.IsComplete, &section.Color)
		if err != nil {
			return domain.Section{}, fmt.Errorf("error while scan section: %w", err)
		}

		return section, nil
	}

	return domain.Section{}, ErrNotFound
}

func (s *PSQLStorage) GetInterview(ctx context.Context, UUID string, UserUUID string) (domain.Interview, error) {

	rows, err := s.db.QueryContext(ctx, "SELECT title, start_timestamp, timing, thread, feedback  FROM interview WHERE uuid = $1 AND owner_uuid = $2", UUID, UserUUID)
	if err != nil {
		return domain.Interview{}, fmt.Errorf("error while query interview: %w", err)
	}

	defer rows.Close()

	var feedback sql.NullString
	for rows.Next() {
		var interview domain.Interview
		var timing int
		var thread string
		err = rows.Scan(&interview.Title, &interview.StartTimestamp, &timing, &thread, &feedback)

		if err != nil {
			return domain.Interview{}, fmt.Errorf("error while scan interview: %w", err)
		}

		interview.UUID = UUID
		interview.Timing = time.Duration(timing)
		interview.SecondsLeft = int(interview.StartTimestamp.Add(time.Duration(int(time.Second) * timing)).Sub(time.Now().UTC()).Seconds())

		if feedback.Valid {
			interview.Feedback = feedback.String
		}

		interview.IsComplete = interview.Feedback != ""

		if thread != "" {
			strs := strings.Split(thread, "||")
			if len(strs) != 2 {
				return domain.Interview{}, errors.New("invalid thread format - it must have delimiter '||' and contain 2 parts")
			}

			interview.Thread = &domain.ChatThread{
				ID:     strs[0],
				Secret: strs[1],
			}
		}

		sections, err := s.getInterviewSections(ctx, UUID)
		if err != nil {
			return domain.Interview{}, fmt.Errorf("cannot get sections: %w", err)
		}

		interview.Sections = sections

		return interview, nil
	}

	return domain.Interview{}, ErrNotFound
}

func (s *PSQLStorage) CompleteSection(ctx context.Context, UUID string, _ string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE section SET is_complete = true WHERE uuid = $1", UUID)
	if err != nil {
		return fmt.Errorf("error while exec complete section query: %w", err)
	}

	return nil
}

func (s *PSQLStorage) StartSection(ctx context.Context, UUID string, _ string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE section SET is_started = true WHERE uuid = $1", UUID)
	if err != nil {
		return fmt.Errorf("error while exec complete section query: %w", err)
	}

	return nil
}

func (s *PSQLStorage) IsUserExists(ctx context.Context, userUUID string) (bool, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM users WHERE uuid = $1", userUUID)
	if err != nil {
		return false, fmt.Errorf("cannot check user in DB: %w", err)
	}

	defer rows.Close()

	return rows.Next(), nil
}

func (s *PSQLStorage) getInterviewSections(ctx context.Context, interviewUUID string) ([]domain.Section, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT name, grade, position, is_started, is_complete, color, uuid FROM section WHERE interview_uuid = $1", interviewUUID)
	if err != nil {
		return nil, fmt.Errorf("error while query sections: %w", err)
	}

	defer rows.Close()

	sections := make([]domain.Section, 0)
	for rows.Next() {
		var section domain.Section
		err = rows.Scan(&section.Name, &section.Grade, &section.Position, &section.IsStarted, &section.IsComplete, &section.Color, &section.UUID)
		if err != nil {
			return nil, fmt.Errorf("error while scan section: %w", err)
		}

		// TODO: we can load all questions in one query
		questions, err := s.getSectionQuestions(ctx, section.UUID)
		if err != nil {
			return nil, fmt.Errorf("cannot get questions for section: %w", err)
		}

		section.Questions = questions
		section.InterviewUUID = interviewUUID

		sections = append(sections, section)
	}

	return sections, nil
}

func (s *PSQLStorage) CompleteInterview(ctx context.Context, UUID string, userUUID string, feedback string) error {
	_, err := s.db.ExecContext(ctx, "UPDATE interview SET feedback = $1 WHERE uuid = $2 AND owner_uuid = $3", feedback, UUID, userUUID)
	if err != nil {
		return fmt.Errorf("error while exec complete interview query: %w", err)
	}

	_, err = s.db.ExecContext(ctx, "UPDATE section SET is_started = true, is_complete = true WHERE interview_uuid = $1", UUID)
	if err != nil {
		return fmt.Errorf("error while exec complete interview sections query: %w", err)
	}

	return nil
}
