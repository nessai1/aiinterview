package domain

import "time"

type User struct {
	UUID string
}

type Interview struct {
	Title     string
	Timing    time.Duration
	StartTime time.Time
	Topics    []Topic
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
	Name  string
	Grade TopicGrade
}
