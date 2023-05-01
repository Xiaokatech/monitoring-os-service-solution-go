package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/judwhite/go-svc"
)

// implements svc.Service
type program struct {
	LogFile *os.File
	svr     *server
	ctx     context.Context
}

func (p *program) Context() context.Context {
	return p.ctx
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	prg := program{
		svr: &server{},
		ctx: ctx,
	}

	defer func() {
		if prg.LogFile != nil {
			if closeErr := prg.LogFile.Close(); closeErr != nil {
				log.Printf("error closing '%s': %v\n", prg.LogFile.Name(), closeErr)
			}
		}
	}()

	// call svc.Run to start your program/service
	// svc.Run will call Init, Start, and Stop
	if err := svc.Run(&prg); err != nil {
		log.Fatal(err)
	}
}

func (p *program) Init(env svc.Environment) error {
	log.Printf("is win service? %v\n", env.IsWindowsService())

	// write to "HelloWorldGoOsService.log" when running as a Windows Service
	if env.IsWindowsService() {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return err
		}

		logPath := filepath.Join(dir, "HelloWorldGoOsService.log")

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
	log.Printf("Starting...\n")

	// Create a new Goroutine to run the server
	go startHTTPServer()

	go p.svr.start()
	return nil
}

func startHTTPServer() {
	// Define the HTTP handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Start the HTTP server on port 9001
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		fmt.Println("Failed to start HTTP server on port 9001:", err)
	}
}

func (p *program) Stop() error {
	log.Printf("Stopping...\n")
	if err := p.svr.stop(); err != nil {
		return err
	}
	log.Printf("Stopped.\n")
	return nil
}
