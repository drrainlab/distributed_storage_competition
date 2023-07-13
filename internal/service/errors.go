package service

import "errors"

var (
	ErrNotEnoughSpace = errors.New("not enough space")
	ErrEmptyFile      = errors.New("empty file")
)
