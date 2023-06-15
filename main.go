//
// Copyright (C) 2023 ANSYS, Inc. Unauthorized use, distribution, or duplication is prohibited.
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/judwhite/go-svc"
)

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

type PIDdata struct {
	PID int `json:"pid"`
}

func WriteMetadataToFile(filePath string, metadata *PIDdata) (*PIDdata, error) {
	// Convert the slice of Usser structs to JSON data
	data, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

type program struct {
	LogFile *os.File
	wg      sync.WaitGroup
	quit    chan struct{}
}

func (p *program) Init(env svc.Environment) error {
	log.Printf("is win service? %v", env.IsWindowsService())

	// write to "ansysCSPAgentManagerServiceApp.log" when running as a Windows Service
	if env.IsWindowsService() {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}

		logPath := filepath.Join(dir, "ansysCSPAgentManagerServiceApp.log")
		log.Println("logPath", logPath)

		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		p.LogFile = f

		log.SetOutput(f)
	}

	return nil
}

func (p *program) CheckAgentRunning() (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 5, // set a timeout of 5 seconds
	}

	resp, err := client.Get("http://localhost:9001")
	if err != nil {
		fmt.Printf("Error making request: %s\n", err.Error())

		if pid, err := RunAgentBinaryFile(); pid != 0 && err == nil {
			fmt.Println("RunAgentBinaryFile is ok on pid", pid)

			osServiceManagerAppName := "ansysCSPAgentManagerService"
			fileName := "pid.json"

			// Set the default appData path for Linux, Windows, and macOS systems
			var agentAppDataPath string = GetAnsysCSPAgentManagerServiceAppPathByAppName(osServiceManagerAppName)
			configFileLocation := filepath.Join(agentAppDataPath, fileName)
			WriteMetadataToFile(configFileLocation, &PIDdata{PID: pid})
		} else if err != nil {
			fmt.Printf("Error running binary file: %s\n", err.Error())
		}
	}

	return resp, err
}

func (p *program) Start() error {
	p.quit = make(chan struct{})

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		p.CheckAgentRunning() // first start for agent

		for {
			select {
			case <-ticker.C:
				fmt.Println("Hello, World! by fmt") // stdout
				log.Println("Hello, World! by log") // stderr

				// === check if agent is running - start ===
				// resp, err := p.CheckAgentRunning()
				// if err != nil {
				// 	continue // continue loop instead of exiting
				// 	// return // exit goroutine
				// }

				// defer resp.Body.Close()

				// fmt.Printf("Response status: %d\n", resp.StatusCode)
				// === check if agent is running - end ===
			case <-p.quit:
				return
			}
		}
	}()

	fmt.Println("the start func will end")

	return nil
}

func (p *program) Stop() error {
	close(p.quit)
	p.wg.Wait()
	return nil
}

func main() {
	prg := &program{}

	if err := svc.Run(prg); err != nil {
		log.Fatal(err)
	}
}

func RunAgentBinaryFile() (int, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return 0, err
	}

	binaryFileName := ""
	switch runtime.GOOS {
	case "linux":
		binaryFileName = "ansysCSPAgentApp"
	case "windows":
		binaryFileName = "ansysCSPAgentApp.exe"
	default:
		fmt.Println("Unsupported operating system")
		os.Exit(1)
	}
	binaryFilePath := filepath.Join(dir, binaryFileName)
	log.Println("logPath", binaryFilePath)

	// Set the path to your binary file
	cmd := exec.Command(binaryFilePath)

	// Start the process
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting process:", err)
		return 0, err
	}

	// Get the process ID
	pid := cmd.Process.Pid
	fmt.Println("Process started with PID:", pid)

	// * Wait for the process to finish, we don't need to wait for it because we don't want this process to be blocked
	// if err := cmd.Wait(); err != nil {
	// 	fmt.Println("Process exited with error:", err)
	// } else {
	// 	fmt.Println("Process exited successfully")
	// }

	return pid, nil
}
