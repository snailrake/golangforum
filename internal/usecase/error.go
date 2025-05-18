package usecase

import "errors"

var (
	ErrMissingToken    = errors.New("missing token")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInvalidUserID   = errors.New("invalid user_id")
	ErrInvalidUsername = errors.New("invalid username")

	ErrInvalidCommentData = errors.New("invalid comment data")
	ErrCommentNotFound    = errors.New("comment not found")

	ErrInvalidPostData = errors.New("invalid post data")
	ErrPostNotFound    = errors.New("post not found")

	ErrInvalidTopicData = errors.New("invalid topic data")
	ErrTopicNotFound    = errors.New("topic not found")
)
