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

func setVerboseFlag(s *server) {
	const (
		usage = "Run server in verbose mode."
	)

	flag.BoolVar(&s.verbose, "verbose", DefaultVerbose, usage)
	flag.BoolVar(&s.verbose, "v", DefaultVerbose, usage+"(shorthand)")
}

func setModeFlag(s *server) {
	modes := strings.Join(AllowedModes, ", ")
	usage := "Set storage mode for server (allowed modes are: " + modes + ")."

	flag.StringVar(&s.mode, "mode", DefaultMode, usage)
	flag.StringVar(&s.mode, "m", DefaultMode, usage+"(shorthand)")
}

func setPortFlag(s *server) {
	const (
		usage = "Port number for server to listen (1-65535)."
	)

	flag.IntVar(&s.port, "port", DefaultPort, usage)
	flag.IntVar(&s.port, "p", DefaultPort, usage+"(shorthand)")
}

func setHostFlag(s *server) {
	const (
		usage = "IP address for server to bind."
	)

	flag.StringVar(&s.host, "host", DefaultHost, usage)
	flag.StringVar(&s.host, "h", DefaultHost, usage+"(shorthand)")
}

func setFlags(s *server) {
	setVerboseFlag(s)
	setModeFlag(s)
	setPortFlag(s)
	setHostFlag(s)
}

func validateMode(s *server) error {
	for _, mode := range AllowedModes {
		if mode == s.mode {
			return nil
		}
	}
	return ParseCliError{time.Now(), fmt.Sprintf("Wrong server mode provided: %v", s.mode)}
}

func validatePort(s *server) error {
	if s.port < 1 || s.port > 65535 {
		return ParseCliError{time.Now(), fmt.Sprintf("Wrong port number provided: %v", s.port)}
	}
	return nil
}

func validateHost(s *server) error {
	ip := net.ParseIP(s.host)
	if ip == nil {
		return ParseCliError{time.Now(), fmt.Sprintf("Wrong host IP provided: %v", s.host)}
	}
	return nil
}

func validateOptions(s *server) error {
	var err error

	err = validateMode(s)
	if err != nil {
		return err
	}
	err = validatePort(s)
	if err != nil {
		return err
	}
	err = validateHost(s)
	if err != nil {
		return err
	}
	return nil
}

func ParseOptions(s *server) error {
	var err error

	setFlags(s)

	flag.Parse()

	err = validateOptions(s)
	if err != nil {
		return err
	}
	return nil
}
