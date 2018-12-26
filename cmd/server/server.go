package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const (
	MemoryMode     = "memory"
	DiskMode       = "disk"
	DefaultPort    = 9090
	DefaultHost    = "127.0.0.1"
	DefaultVerbose = false
	DefaultMode    = MemoryMode
)

var AllowedModes = []string{MemoryMode, DiskMode}

type server struct {
	verbose  bool
	mode     string
	port     int
	host     string
	listener net.Listener
}

func newServer() *server {
	s := new(server)
	s.mode = MemoryMode
	s.port = DefaultPort
	s.host = DefaultHost
	return s
}

func (s *server) close() (err error) {
	if s.listener != nil {
		err = s.listener.Close()
	}
	return
}

func handleCommand(conn net.Conn, commands chan chan string) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		command := scanner.Text()
		command_chan := make(chan string)
		log.Printf("Received command %s, from %s", command, conn.RemoteAddr())
		commands <- command_chan
		command_chan <- command
		io.WriteString(conn, <-command_chan)
	}
}

func handleStorage(db *database, commands chan chan string) {
	var command_str string
	var reply string

	for command_chan := range commands {
		select {
		case command_str = <-command_chan:
			command := strings.Split(command_str, " ")
			if len(command) < 1 {
				command_chan <- "Please provide a command\n"
				break
			}
			switch command[0] {
			case "set":
				if len(command) != 3 {
					command_chan <- "Wrong format for command set. Please use set [key] [value]\n"
					break
				}
				a := &args{command[1], command[2]}
				db.set(a, &reply)
				command_chan <- reply
			case "get":
				if len(command) != 2 {
					command_chan <- "Wrong format for command get. Please use get [key]\n"
					break
				}
				a := &args{command[1], ""}
				db.get(a, &reply)
				command_chan <- reply
			case "del":
				if len(command) != 2 {
					command_chan <- "Wrong format for command del. Please use del [key]\n"
					break
				}
				a := &args{command[1], ""}
				db.del(a, &reply)
				command_chan <- reply
			case "keys":
				if len(command) != 2 {
					command_chan <- "Wrong format for command keys. Please use keys [regex]\n"
					break
				}
				db.keys(command[1], &reply)
				command_chan <- reply
			default:
				log.Println("Wrong command received")
				command_chan <- "Supported commands are: set, get, del\n"
			}
		}
	}
}

func (s *server) start() (err error) {
	db := newDatabase()

	address := s.host + ":" + strconv.Itoa(s.port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalln("Listen error:", err)
	}
	s.listener = l

	log.Printf("Server is running on %s\n", address)

	commands := make(chan chan string)
	go handleStorage(db, commands)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalln("Accept error:", err)
		}

		go handleCommand(conn, commands)
	}
}
