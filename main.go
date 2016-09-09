package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	hupChan := make(chan bool)
	exitChan := make(chan bool)

	// signal handler
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

func doSomething(doneChan chan bool) {
	go func() {
		log.Println("do something...")
		time.Sleep(3 * time.Second)
		log.Println("done")
		doneChan <- true
	}()
}
