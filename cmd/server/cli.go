package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
)

type ParseCliError struct {
	time   time.Time
	reason string
}

func (err ParseCliError) Error() string {
	return fmt.Sprintf("%v: %v", err.time.Format(time.RFC3339), err.reason)
}

func setVerboseFlag(server *Server) {
	const (
		usage = "Run server in verbose mode."
	)

	flag.BoolVar(&server.verbose, "verbose", DefaultVerbose, usage)
	flag.BoolVar(&server.verbose, "v", DefaultVerbose, usage+"(shorthand)")
}

func setModeFlag(server *Server) {
	modes := strings.Join(AllowedModes, ", ")
	usage := "Set storage mode for server (allowed modes are: " + modes + ")."

	flag.StringVar(&server.mode, "mode", DefaultMode, usage)
	flag.StringVar(&server.mode, "m", DefaultMode, usage+"(shorthand)")
}

func setPortFlag(server *Server) {
	const (
		usage = "Port number for server to listen (1-65535)."
	)

	flag.IntVar(&server.port, "port", DefaultPort, usage)
	flag.IntVar(&server.port, "p", DefaultPort, usage+"(shorthand)")
}

func setHostFlag(server *Server) {
	const (
		usage = "IP address for server to bind."
	)

	flag.StringVar(&server.host, "host", DefaultHost, usage)
	flag.StringVar(&server.host, "h", DefaultHost, usage+"(shorthand)")
}

func setFlags(server *Server) {
	setVerboseFlag(server)
	setModeFlag(server)
	setPortFlag(server)
	setHostFlag(server)
}

func validateMode(server *Server) error {
	for _, mode := range AllowedModes {
		if mode == server.mode {
			return nil
		}
	}
	return ParseCliError{time.Now(), fmt.Sprintf("Wrong server mode provided: %v", server.mode)}
}

func validatePort(server *Server) error {
	if server.port < 1 || server.port > 65535 {
		return ParseCliError{time.Now(), fmt.Sprintf("Wrong port number provided: %v", server.port)}
	}
	return nil
}

func validateHost(server *Server) error {
	ip := net.ParseIP(server.host)
	if ip == nil {
		return ParseCliError{time.Now(), fmt.Sprintf("Wrong host IP provided: %v", server.host)}
	}
	return nil
}

func validateOptions(server *Server) error {
	var err error

	err = validateMode(server)
	if err != nil {
		return err
	}
	err = validatePort(server)
	if err != nil {
		return err
	}
	err = validateHost(server)
	if err != nil {
		return err
	}
	return nil
}

func ParseOptions(server *Server) error {
	var err error

	setFlags(server)

	flag.Parse()

	err = validateOptions(server)
	if err != nil {
		return err
	}
	return nil
}
