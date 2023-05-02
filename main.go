package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/judwhite/go-svc"
)

type program struct {
	LogFile *os.File
	wg      sync.WaitGroup
	quit    chan struct{}
}

func (p *program) Init(env svc.Environment) error {
	log.Printf("is win service? %v", env.IsWindowsService())

	// Check env and open log file for linux and windows

	// write to "HelloWorldGoOsService.log" when running as a Windows Service
	if env.IsWindowsService() {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}

		logPath := filepath.Join(dir, "HelloWorldGoOsService.log")
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

func (p *program) Start() error {
	p.quit = make(chan struct{})

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("Hello, World! by fmt") // stdout
				log.Println("Hello, World! by log") // stderr
			case <-p.quit:
				return
			}
		}
	}()

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
