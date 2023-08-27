package storage

import (
	"database/sql"
	"fmt"
	"time"
)

type Segment struct {
	ID        int
	IsActive  bool
	CreatedAt time.Time
	DeletedAt sql.NullTime
}
type SegmentAlreadyExistsForUserError struct {
	Slug string
}

type UserNotExists struct {
	Id int
}

type SegmentNotExists struct {
	Slug string
}

type SegmentsNotFound struct {
}

type SegmentAlreadyExistsError struct {
	Slug string
}

type CSVError struct {
}

func (c CSVError) Error() string {
	return fmt.Sprintf("The record not found for this period of time")
}

func (o SegmentAlreadyExistsError) Error() string {
	return fmt.Sprintf("segment %s already exists", o.Slug)
}
func (s SegmentsNotFound) Error() string {
	return fmt.Sprintf("no any active segment for this user was found")
}

func (e SegmentAlreadyExistsForUserError) Error() string {
	return fmt.Sprintf("user already has the segment: %s", e.Slug)
}

func (u UserNotExists) Error() string {
	return fmt.Sprintf("user %d not exists", u.Id)
}

func (s SegmentNotExists) Error() string {
	return fmt.Sprintf("slug %s not exists", s.Slug)
}
