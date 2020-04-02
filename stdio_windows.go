// +build windows

package main

import (
	"os"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

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
		return errors.Wrap(err, "load dll error")
	}

	setStdHandle, err := syscall.GetProcAddress(syscall.Handle(kernel), "SetStdHandle")
	if err != nil {
		return errors.Wrap(err, "GetProcAddress error")
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
		return errors.Wrap(errno, "syscall error")
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
		return errors.Wrap(errno, "syscall error")
	}

	*os.Stdout = *stdout

	return nil
}
