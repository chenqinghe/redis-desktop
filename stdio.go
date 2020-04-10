package main

import (
	"fmt"
	"os"
)

type FileDescriptor interface {
	Fd() uintptr
}

func NewFileDescriptor(file string) (*os.File, error) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("OpenFile error: %s", err)
	}

	return f, nil
}
