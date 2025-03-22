package domain

import "time"

type User struct {
	UUID string
}

type Interview struct {
	UUID           string        `json:"uuid"`
	Title          string        `json:"title"`
	Timing         time.Duration `json:"timing"` // in seconds, cuz need to show timer HH:MM:SS
	StartTimestamp time.Time     `json:"start_timestamp"`

	IsComplete bool `json:"complete"` // computed -> time.Now() > StartTimestamp + Timing

	Summarize Summarize   `json:"summarize,omitempty"`
	Sections  []Section   `json:"sections"`
	Thread    *ChatThread `json:"-"`
}

type Section struct {
	UUID          string     `json:"uuid"`
	InterviewUUID string     `json:"interview_uuid"`
	Name          string     `json:"name"`
	Grade         TopicGrade `json:"grade"`
	ActualGrade   TopicGrade `json:"actualGrade"`
	Position      int        `json:"position"`
	IsStarted     bool       `json:"isStarted"`
	IsComplete    bool       `json:"isComplete"`
	Questions     []Question `json:"questions"`
	Color         string
}

type Question struct {
	UUID          string `json:"uuid"`
	SectionUUID   string `json:"section_uuid"`
	InterviewUUID string `json:"interview_uuid"`
	Question      string `json:"question"`
	Answer        string `json:"answer"`
	Feedback      string `json:"feedback"`
	Done          bool   `json:"done"`
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
