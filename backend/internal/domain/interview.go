package domain

import "time"

type User struct {
	UUID string
}

type Interview struct {
	UUID           string        `json:"-"`
	Title          string        `json:"title"`
	Timing         time.Duration `json:"timing"`
	StartTimestamp time.Time     `json:"start_timestamp"`
	Topics         []Topic       `json:"topics"`
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
