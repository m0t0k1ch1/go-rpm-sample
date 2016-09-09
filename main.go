package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const DefaultPidFilePath = "/var/run/go-rpm-sample.pid"

func main() {
	var pidFilePath = flag.String("pidfile", DefaultPidFilePath, "pid file path")
	flag.Parse()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	hupChan := make(chan bool)
	exitChan := make(chan bool)

	// signal handling
	go func() {
		for s := range signalChan {
			if s == syscall.SIGHUP {
				log.Println("received SIGHUP")
				hupChan <- true
			} else {
				log.Println("received SIGTERM or SIGQUIT")
				exitChan <- true
			}
		}
	}()

	// pid file handling
	if err := createPidFile(*pidFilePath); err != nil {
		log.Fatal(err)
	}
	defer removePidFile(*pidFilePath)

	doneChan := make(chan bool)

	// main loop
	for {
		doSomething(doneChan)

		select {
		case <-doneChan:
		case <-hupChan:
			<-doneChan
			log.Println("hup")
		case <-exitChan:
			<-doneChan
			log.Println("exit")
			return
		}
	}
}

func createPidFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d", os.Getpid())

	return err
}

func removePidFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Fatal(err)
	}
}

func doSomething(doneChan chan bool) {
	go func() {
		log.Println("do something...")
		time.Sleep(3 * time.Second)
		log.Println("done")
		doneChan <- true
	}()
}
