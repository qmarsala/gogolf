// +build windows

package ui

import (
	"syscall"
	"unsafe"
)

var (
	kernel32        = syscall.NewLazyDLL("kernel32.dll")
	procReadConsole = kernel32.NewProc("ReadConsoleInputW")
	procGetStdHandle = kernel32.NewProc("GetStdHandle")
)

const (
	stdInputHandle = uint32(4294967286) // -10 in uint32
	keyEvent       = 0x0001
	keyDown        = 0x0001
)

type inputRecord struct {
	EventType uint16
	_         uint16
	KeyEvent  keyEventRecord
}

type keyEventRecord struct {
	KeyDown         uint32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

// waitForSpacebar waits for a spacebar keypress on Windows
func (pm *PowerMeter) waitForSpacebar() {
	handle, _, _ := procGetStdHandle.Call(uintptr(stdInputHandle))

	for {
		var record inputRecord
		var numRead uint32

		ret, _, _ := procReadConsole.Call(
			handle,
			uintptr(unsafe.Pointer(&record)),
			1,
			uintptr(unsafe.Pointer(&numRead)),
			0,
		)

		if ret == 0 || numRead == 0 {
			continue
		}

		// Check for spacebar (VK_SPACE = 0x20) and key down event
		if record.EventType == keyEvent &&
			record.KeyEvent.KeyDown == keyDown &&
			record.KeyEvent.VirtualKeyCode == 0x20 {
			return
		}
	}
}
