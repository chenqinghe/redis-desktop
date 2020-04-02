package main

import (
	"github.com/pkg/errors"
	"os"
)

type FileDescriptor interface {
	Fd() uintptr
}

func NewFileDescriptor(file string) (*os.File, error) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "open file error")
	}

	return f, nil
}
