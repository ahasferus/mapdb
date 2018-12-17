package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func signalHandler() {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	log.Println("Received signal")
}

func errorhandler(err error) {
	if err != nil {
		log.Panicln(err)
		os.Exit(1)
	}
}

func main() {
	client := NewClient()

	err := ParseOptions(client)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		signalHandler()
		os.Exit(0)
	}()

	errorhandler(client.Start())
	fmt.Println("exit")
}
