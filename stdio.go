// +build windows

package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
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

var (
	// to avoid GC
	stdError  FileDescriptor
	stdOutput FileDescriptor
)

func SetStdHandle(stdout, stderr *os.File) error {
	stdError = stderr
	stdOutput = stdout

	kernel, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return fmt.Errorf("LoadLibrary: %s", err)
	}

	setStdHandle, err := syscall.GetProcAddress(syscall.Handle(kernel), "SetStdHandle")
	if err != nil {
		return fmt.Errorf("GetProcAddress: SetStdHandle: %s", err)
	}

	t1 := syscall.STD_ERROR_HANDLE
	errHandle := *(*uint32)(unsafe.Pointer(&t1))
	ret, _, errno := syscall.Syscall(
		uintptr(setStdHandle),
		2,
		uintptr(errHandle),
		uintptr(stdError.Fd()),
		0,
	)

	if ret == 0 || errno != 0 {
		return fmt.Errorf("Syscall: Returned value not 0: %d", ret)
	}
	*os.Stderr = *stderr

	t2 := syscall.STD_OUTPUT_HANDLE
	outHandle := *(*uint32)(unsafe.Pointer(&t2))
	ret, _, errno = syscall.Syscall(
		uintptr(setStdHandle),
		2,
		uintptr(outHandle),
		uintptr(stdOutput.Fd()),
		0,
	)

	if ret == 0 || errno != 0 {
		return fmt.Errorf("Syscall: Returned value not 0: %d", ret)
	}

	*os.Stdout = *stdout

	return nil
}
