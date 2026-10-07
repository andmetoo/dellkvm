//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

func userConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "dellkvm", "config.toml"), nil
}

func openFile(path string) error {
	verb, _ := syscall.UTF16PtrFromString("open")
	file, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	shellExecute := syscall.NewLazyDLL("shell32.dll").NewProc("ShellExecuteW")
	code, _, callErr := shellExecute.Call(0, uintptr(unsafe.Pointer(verb)), uintptr(unsafe.Pointer(file)), 0, 0, 1)
	if code <= 32 {
		return fmt.Errorf("open %s: Windows error %d: %v", path, code, callErr)
	}
	return nil
}
