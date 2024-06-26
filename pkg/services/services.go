package services

import "errors"

var (
	ErrVideoNotFound = errors.New("video not found")
	ErrBadURL = errors.New("bad video url")
)
