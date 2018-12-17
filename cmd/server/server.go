package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
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

type Server struct {
	verbose  bool
	mode     string
	port     int
	host     string
	listener net.Listener
}

type Database struct {
	storage map[string]string
}

type Args struct {
	Key   string
	Value string
}

func (db *Database) Set(args *Args, reply *string) error {
	db.storage[args.Key] = args.Value
	*reply = ""
	return nil
}

func (db *Database) Get(args *Args, reply *string) error {
	*reply = db.storage[args.Key]
	return nil
}

func NewServer() *Server {
	server := new(Server)
	server.mode = MemoryMode
	server.port = DefaultPort
	server.host = DefaultHost
	return server
}

var DB Database

func NewDatabase() *Database {
	database := new(Database)
	database.storage = make(map[string]string)
	return database
}

func (server *Server) Close() (err error) {
	if server.listener != nil {
		err = server.listener.Close()
	}
	return
}

func (server *Server) Start() (err error) {
	DB := NewDatabase()

	rpc.Register(DB)

	rpc.HandleHTTP()

	address := server.host + ":" + strconv.Itoa(server.port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	server.listener = l

	http.Serve(server.listener, nil)
	return
}
