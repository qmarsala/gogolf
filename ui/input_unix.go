// +build !windows

package ui

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	tcgets = 0x5401
	tcsets = 0x5402
)

type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [32]uint8
	Ispeed uint32
	Ospeed uint32
}

// waitForSpacebar waits for a spacebar keypress on Unix-like systems
func (pm *PowerMeter) waitForSpacebar() {
	// Get current terminal settings
	fd := int(os.Stdin.Fd())
	var oldState termios
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcgets, uintptr(unsafe.Pointer(&oldState)))

	// Set terminal to raw mode (no echo, no line buffering)
	newState := oldState
	newState.Lflag &^= syscall.ECHO | syscall.ICANON
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&newState)))

	// Restore terminal settings when done
	defer func() {
		_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&oldState)))
	}()

	// Read single character
	buf := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buf)
		if err == nil && buf[0] == ' ' {
			return
		}
	}
}

// WaitForAnyKey waits for any key press on Unix-like systems
func WaitForAnyKey() {
	// Get current terminal settings
	fd := int(os.Stdin.Fd())
	var oldState termios
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcgets, uintptr(unsafe.Pointer(&oldState)))

	// Set terminal to raw mode (no echo, no line buffering)
	newState := oldState
	newState.Lflag &^= syscall.ECHO | syscall.ICANON
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&newState)))

	// Restore terminal settings when done
	defer func() {
		_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), tcsets, uintptr(unsafe.Pointer(&oldState)))
	}()

	// Read any single character
	buf := make([]byte, 1)
	os.Stdin.Read(buf)
}
