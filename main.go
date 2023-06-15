//
// Copyright (C) 2023 ANSYS, Inc. Unauthorized use, distribution, or duplication is prohibited.
//

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	TTools "AnsysCSPAgentManagerService/src/tools"

	"github.com/judwhite/go-svc"
)

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

func (p *program) StartNewAgentApp(agentManagerServiceConfigFileLocation string) {
	if pid, err := RunAgentBinaryFile(); pid != 0 && err == nil {
		fmt.Println("RunAgentBinaryFile is ok on pid", pid)
		TTools.WritePidDataToFile(agentManagerServiceConfigFileLocation, &TTools.PIDdata{PID: pid})
	} else if err != nil {
		fmt.Printf("Error running binary file: %s\n", err.Error())
	}
}

func (p *program) ProcessExists(pid int) (bool, error) {
	if runtime.GOOS == "windows" {
		isProcessExists, err := TTools.ProcessExists_windows(pid)
		return isProcessExists, err
	} else if runtime.GOOS == "linux" {
		isProcessExists, err := TTools.ProcessExists_linux(pid)
		return isProcessExists, err
	}

	return false, nil
}

func (p *program) CheckAgentRunning() (bool, error) {

	osServiceManagerAppName := "ansysCSPAgentManagerService"
	fileName := "pid.json"

	// Set the default appData path for Linux, Windows, and macOS systems
	var agentManagerServiceAppDataPath string = TTools.GetAnsysCSPAgentManagerServiceAppPathByAppName(osServiceManagerAppName)
	agentManagerServiceConfigFileLocation := filepath.Join(agentManagerServiceAppDataPath, fileName)

	pidData, err := TTools.ReadPidDataFromFile(agentManagerServiceConfigFileLocation)
	if err != nil {
		fmt.Printf("Error reading pid data from file: %s\n", err.Error())
	}

	// === See pid - start ===
	isProcessExists, err := p.ProcessExists(pidData.PID)
	if err != nil {
		fmt.Printf("Failed to find process: %s\n", err)
		fmt.Println("Starting new agent...")
		p.StartNewAgentApp(agentManagerServiceConfigFileLocation)
	}
	// === See pid - end ===

	return isProcessExists, err
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
