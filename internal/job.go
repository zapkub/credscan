package internal

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type FindingLocationPositionBegin struct {
	Line int `json:"line"`
}
type FindingLocationPosition struct {
	Begin *FindingLocationPositionBegin `json:"begin"`
}
type FindingLocation struct {
	Path     string                   `json:"path"`
	Position *FindingLocationPosition `json:"position"`
}
type Finding struct {
	Type     string           `json:"type"`
	RuleID   string           `json:"ruleId"`
	Location *FindingLocation `json:"location"`
}

type Findings struct {
	Findings []*Finding `json:"findings"`
}

func (f *Findings) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, &f)
	}
	return nil
}

func (f Findings) Value() (driver.Value, error) {
	return json.Marshal(f)
}

type JobStatus int64

var jobStatusName = map[JobStatus]string{
	JobStatusQueued:     "Queued",
	JobStatusInProgress: "InProgress",
	JobStatusSuccess:    "Success",
	JobStatusFailure:    "Failure",
}

func (j JobStatus) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", jobStatusName[j])), nil
}

func (j *JobStatus) UnmarshalJSON([]byte) error {
	panic("unimplemented")
}

const (
	JobStatusQueued JobStatus = iota
	JobStatusInProgress
	JobStatusSuccess
	JobStatusFailure
)

type Job struct {
	ID             int        `json:"id"`
	RepositoryURL  string     `json:"repositoryUrl"`
	RepositoryName string     `json:"repositoryName"`
	Findings       *Findings  `json:"findings"`
	Status         JobStatus  `json:"status"`
	QueuedAt       *time.Time `json:"queuedAt"`
	ScanningAt     *time.Time `json:"scanningAt"`
	FinnishedAt    *time.Time `json:"finnishedAt"`
}
