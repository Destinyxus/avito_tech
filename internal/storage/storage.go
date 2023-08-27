package storage

import "fmt"

type SegmentAlreadyExistsError struct {
	Slug string
}

type UserNotExists struct {
	Id int
}

type SegmentNotExists struct {
	Slug string
}

func (e SegmentAlreadyExistsError) Error() string {
	return fmt.Sprintf("user already has the segment: %s", e.Slug)
}

func (u UserNotExists) Error() string {
	return fmt.Sprintf("user %d not exists", u.Id)
}

func (s SegmentNotExists) Error() string {
	return fmt.Sprintf("slug %s not exists", s.Slug)
}
