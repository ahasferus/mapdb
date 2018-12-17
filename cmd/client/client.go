package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	Key   string
	Value string
}

const (
	DefaultPort = 9090
	DefaultHost = "127.0.0.1"
)

type Client struct {
	serverHost string
	serverPort int
}

func NewClient() *Client {
	client := new(Client)
	client.serverHost = DefaultHost
	client.serverPort = DefaultPort
	return client
}

func Set(c *rpc.Client, key, value string) (err error) {
	args := &Args{key, value}
	var reply string

	err = c.Call("Database.Set", args, &reply)
	if err != nil {
		log.Fatal("Call failed:", err)
		return
	}
	fmt.Printf("Call succeed\n")
	return
}

func Get(c *rpc.Client, key string) (err error) {
	args := &Args{Key: key}
	var reply string

	err = c.Call("Database.Get", args, &reply)
	if err != nil {
		log.Fatal("Call failed:", err)
		return
	}
	fmt.Printf("Call succeed. Response: %s\n", reply)
	return
}

func (client *Client) Start() (err error) {
	address := client.serverHost + ":" + strconv.Itoa(client.serverPort)
	c, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		log.Fatal("Dialing error:", err)
		return
	}
	scanner := bufio.NewScanner(os.Stdin)
	var command []string
	for scanner.Scan() {
		command = strings.Split(scanner.Text(), " ")
		if len(command) < 1 || len(command) > 3 {
			log.Fatal("Wrong command provided %v:", command)
			continue
		}
		switch command[0] {
		case "set":
			if len(command) != 3 {
				log.Fatal("Wrong set command provided %v:", command)
				continue
			}
			err := Set(c, command[1], command[2])
			if err != nil {
				return err
			}
		case "get":
			if len(command) != 2 {
				log.Fatal("Wrong get command provided %v:", command)
				continue
			}
			err := Get(c, command[1])
			if err != nil {
				return err
			}
		case "quit":
			return
		default:
			fmt.Println("Supported commands are: set, get, help, quit.")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Reading standard input:", err)
	}

	return
}
