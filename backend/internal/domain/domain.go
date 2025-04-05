package domain

import "time"

type User struct {
	UUID string
}

type Interview struct {
	UUID           string        `json:"uuid"`
	Title          string        `json:"title"`
	Timing         time.Duration `json:"timing"` // in seconds, cuz need to show timer HH:MM:SS
	SecondsLeft    int           `json:"seconds_left"`
	StartTimestamp time.Time     `json:"start_timestamp"`

	IsComplete bool `json:"complete"` // feedback != ""

	Feedback string      `json:"feedback,omitempty"`
	Sections []Section   `json:"sections"`
	Thread   *ChatThread `json:"-"`
}

func (i *Interview) GetActiveSection() *Section {
	if i.Sections == nil {
		return nil
	}

	for _, section := range i.Sections {
		if !section.IsComplete && section.IsStarted {
			return &section
		}
	}

	return nil
}

type Section struct {
	UUID          string     `json:"uuid"`
	InterviewUUID string     `json:"interview_uuid"`
	Name          string     `json:"name"`
	Grade         TopicGrade `json:"grade"`
	Position      int        `json:"position"`
	IsStarted     bool       `json:"isStarted"`
	IsComplete    bool       `json:"isComplete"`
	Questions     []Question `json:"questions"`
	Color         string     `json:"color"`
}

func (s *Section) GetActiveQuestion() *Question {
	if s.Questions == nil {
		return nil
	}

	for _, question := range s.Questions {
		if !question.Done {
			return &question
		}
	}

	return nil
}

type Question struct {
	UUID          string `json:"uuid"`
	SectionUUID   string `json:"section_uuid"`
	InterviewUUID string `json:"interview_uuid"`
	Question      string `json:"question"`
	Answer        string `json:"answer"`
	Feedback      string `json:"feedback"`
	Done          bool   `json:"done"`
	Position      int    `json:"position"`
}

type Summarize struct {
}

type Assistant struct {
	ID         string
	Model      string
	ExternalID string
}

type TopicGrade string

const (
	GradeJunior TopicGrade = "junior"
	GradeMiddle TopicGrade = "middle"
	GradeSenior TopicGrade = "senior"
)

type Topic struct {
	Name  string     `json:"name"`
	Grade TopicGrade `json:"grade"`
}

type ChatThread struct {
	ID     string
	Secret string
}
