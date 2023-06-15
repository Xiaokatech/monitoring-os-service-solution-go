package ProcessCheck

import "os"

type ProcessCheckerWindows struct{}

func (p ProcessCheckerWindows) ProcessExists(pid int) (bool, error) {
	_, err := os.FindProcess(pid)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func NewProcessChecker() ProcessChecker {
	return ProcessCheckerWindows{}
}
