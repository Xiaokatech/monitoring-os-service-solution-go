package ProcessCheck

import (
	"os"
	"syscall"
)

type ProcessChecker interface {
	ProcessExists(pid int) (bool, error)
}

type ProcessCheckerLinux struct{}

func (p ProcessCheckerLinux) ProcessExists(pid int) (bool, error) {
	err := syscall.Kill(pid, 0)
	if err == nil {
		return true, nil
	}
	if err == syscall.ESRCH {
		return false, nil
	}
	return false, err
}

type ProcessCheckerWindows struct{}

func (p ProcessCheckerWindows) ProcessExists(pid int) (bool, error) {
	_, err := os.FindProcess(pid)
	if err != nil {
		return false, nil
	}
	return true, nil
}
