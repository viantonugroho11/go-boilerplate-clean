package sample

import "errors"

var (
	ErrSampleIDRequired = errors.New("sample id is required")
	ErrSampleNotFound   = errors.New("sample not found")
)
