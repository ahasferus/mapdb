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
	server := newServer()

	err := ParseOptions(server)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer server.close()

	go func() {
		signalHandler()
		server.close()
		os.Exit(0)
	}()

	errorhandler(server.start())
	fmt.Println("exit")
}
