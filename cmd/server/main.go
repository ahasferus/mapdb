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
	server := NewServer()

	err := ParseOptions(server)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer server.Close()

	go func() {
		signalHandler()
		server.Close()
		os.Exit(0)
	}()

	errorhandler(server.Start())
	fmt.Println("exit")
}
