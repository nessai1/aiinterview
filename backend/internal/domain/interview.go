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
	Topics         []Topic       `json:"topics"`

	IsComplete bool `json:"complete"` // computed -> time.Now() > StartTimestamp + Timing
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
