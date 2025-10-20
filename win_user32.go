//go:build windows

package main

import (
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Load user32.dll and kernel32.dll and prepare procedure pointers for Windows API calls.
var (
	user32               = windows.NewLazySystemDLL("user32.dll")
	kernel32             = windows.NewLazySystemDLL("kernel32.dll")
	procGetModuleHandle  = kernel32.NewProc("GetModuleHandleW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procDispatchMessageW = user32.NewProc("DispatchMessageW")
	procGetMessageW      = user32.NewProc("GetMessageW")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
)

// Windows message and window style constants.
const (
	WM_DESTROY       = 0x0002
	WM_SETTINGCHANGE = 0x001A
	CW_USEDEFAULT    = 0x80000000
	WS_OVERLAPPED    = 0x00000000
)

// HWND and HINSTANCE are Windows handle types.
type (
	HWND      uintptr
	HINSTANCE uintptr
)

// WNDCLASSEXW defines a window class for the Windows API.
type WNDCLASSEXW struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     HINSTANCE
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  *uint16
	LpszClassName *uint16
	HIconSm       uintptr
}

// MSG represents a Windows message.
type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

// getModuleHandle returns a handle to the current module (process).
//
// Returns:
//   - windows.Handle: handle to the current module
//   - error: non-nil if the call fails
func getModuleHandle() (windows.Handle, error) {
	ret, _, err := procGetModuleHandle.Call(0)
	handle := windows.Handle(ret)
	if handle == 0 {
		return 0, err
	}
	return handle, nil
}

// wndProc is the window procedure callback for handling Windows messages.
//
// Arguments:
//   - hwnd: window handle
//   - msg: message identifier
//   - wparam: message parameter
//   - lparam: message parameter
//
// Returns:
//   - uintptr: result of message processing
func wndProc(hwnd HWND, msg uint32, wparam, lparam uintptr) uintptr {
	switch msg {
	case WM_SETTINGCHANGE:
		// WM_SETTINGCHANGE is sent when system settings change, including environment variables.
		if lparam != 0 {
			str := windows.UTF16PtrToString((*uint16)(unsafe.Pointer(lparam)))  //nolint:govet
			if str == "Environment" {
				log.Println("Environment changed")
			}
		}
	case WM_DESTROY:
		// Post a quit message to end the message loop.
		procPostQuitMessage.Call(0)  //nolint:errcheck
		return 0
	}
	// Call the default window procedure for unhandled messages.
	ret, _, _ := procDefWindowProcW.Call(
		uintptr(hwnd), uintptr(msg), wparam, lparam)
	return ret
}

// watch creates a hidden window and enters a message loop to listen for environment changes.
//
// This function blocks until the window receives a WM_DESTROY message.
func watch() {
	instance, err := getModuleHandle()
	if err != nil {
		panic(err)
	}

	className, err := syscall.UTF16PtrFromString("EnvWatcherClass")
	if err != nil {
		panic(err)
	}

	// Register the window class.
	wc := WNDCLASSEXW{
		CbSize:        uint32(unsafe.Sizeof(WNDCLASSEXW{})),
		LpfnWndProc:   syscall.NewCallback(wndProc),
		HInstance:     HINSTANCE(instance),
		LpszClassName: className,
	}

	atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
	if atom == 0 {
		panic(fmt.Errorf("RegisterClassExW failed: %v", err))
	}

	windowName, err := syscall.UTF16PtrFromString("EnvWatcher")
	if err != nil {
		panic(err)
	}

	// Create the window (hidden, message-only).
	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		WS_OVERLAPPED,
		CW_USEDEFAULT, CW_USEDEFAULT,
		CW_USEDEFAULT, CW_USEDEFAULT,
		0, 0, uintptr(instance), 0)

	if hwnd == 0 {
		panic(fmt.Errorf("CreateWindowExW failed: %v", err))
	}

	// Enter the message loop.
	var msg MSG
	for {
		ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if int32(ret) == 0 {
			break // WM_QUIT received, exit loop
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))  //nolint:errcheck
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))  //nolint:errcheck
	}
}
