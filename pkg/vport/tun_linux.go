package vport

import (
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	devNet  = "/dev/net"
	tunFile = "/dev/net/tun"
)

type Tap struct {
	io.ReadWriteCloser
	name string
}

// New creates a new TAP interface with the given name.
func New(name string) (*Tap, error) {
	var (
		fd  int
		err error
	)

	switch fd, err = syscall.Open(
		"/dev/net/tun", os.O_RDWR|syscall.O_NONBLOCK, 0,
	); {
	// in some linux containers, the /dev/net/tun file does not exist
	// we have to create and mknod it.
	case err == syscall.ENOENT:
		if err := createSpecialFile(tunFile); err != nil {
			return nil, err
		}
		// FIXME: Beware of recursive infinite loop here.
		// If we get ENOENT again, we will loop forever. ;_;
		return New(name)
	case err == nil:
		// continue
	default:
		return nil, err
	}

	name, err = ioctlCreateInterface(uintptr(fd), name)
	if err != nil {
		return nil, fmt.Errorf("ioctlCreateInterface failed: %w", err)
	}

	return &Tap{
		ReadWriteCloser: os.NewFile(uintptr(fd), "tun"),
		name:            name,
	}, nil
}

func createSpecialFile(path string) error {
	if err := os.MkdirAll(devNet, 0755); err != nil {
		return err
	}

	if err := syscall.Mknod(path, syscall.S_IFCHR|0666, int(unix.Mkdev(10, 200))); err != nil {
		return err
	}
	return nil
}
