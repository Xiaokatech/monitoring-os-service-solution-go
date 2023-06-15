package TTools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

func ProcessExists_windows(pid int) (bool, error) {
	_, err := os.FindProcess(pid)
	return err == nil, err
}

func ProcessExists_linux(pid int) (bool, error) {
	err := syscall.Kill(pid, 0)
	if err == nil {
		return true, err
	}
	if err == syscall.ESRCH {
		return false, err
	}
	return false, nil

}

type PIDdata struct {
	PID int `json:"pid"`
}

func GetAnsysCSPAgentManagerServiceAppPathByAppName(osServiceManagerAppName string) string {
	fmt.Println("GetAnsysCSPAgentManagerServiceAppPathByAppName - start")

	var appDataByAppNamePath string
	switch runtime.GOOS {
	case "linux":
		appDataByAppNamePath = filepath.Join("/usr/local/go", osServiceManagerAppName)
	case "windows":
		appDataByAppNamePath = filepath.Join("C:\\go", osServiceManagerAppName)
	case "darwin":
		appDataByAppNamePath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", osServiceManagerAppName)
	default:
		fmt.Println("Unsupported operating system")
		os.Exit(1)
	}

	fmt.Println("GetAnsysCSPAgentManagerServiceAppPathByAppName:", appDataByAppNamePath)

	return appDataByAppNamePath
}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func ReadPidDataFromFile(filePath string) (*PIDdata, error) {
	// Read the contents of the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse the JSON data in to a slice of User structs
	var pidData *PIDdata
	err = json.Unmarshal(data, &pidData)
	if err != nil {
		return nil, err
	}

	return pidData, nil
}

func WritePidDataToFile(filePath string, pidData *PIDdata) (*PIDdata, error) {
	// Convert the slice of Usser structs to JSON data
	data, err := json.Marshal(pidData)
	if err != nil {
		return nil, err
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return nil, err
	}

	return pidData, nil
}
