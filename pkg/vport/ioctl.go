package vport

import (
	"bytes"
	"os"
	"syscall"
	"unsafe"
)

// ioctl is a bit of a rabbit hole to go down if you're not using packages to
// help out with the various bit-shifts, masks, and constants. So for the sake
// of keeping this repository dependency-free, I'm going to include some magic
// constants that will help out with the ioctl calls that we need.
const (
	TAPSETIFFflag = 0x1002
)

// TODO(alexb) - Add some commentary around the above magic function, and how it
// works with ioctl if you're brave enough.
func ioctlCreateInterface(fd uintptr, name string) (string, error) {
	type ioctlReq struct {
		Name  [0x10]byte
		Flags uint16
		_     [0x28 - 0x10 - 2]byte
	}
	var req = ioctlReq{
		Flags: TAPSETIFFflag,
	}
	copy(req.Name[:], name)
	// FIXME: the below is a bit nasty to look at, fix this up at some point
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(syscall.TUNSETIFF),
		uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		return "", os.NewSyscallError("ioctl", errno)
	}
	return string(bytes.Trim(req.Name[:], "\x00")), nil
}
