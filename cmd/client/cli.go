package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type ParseCliError struct {
	time   time.Time
	reason string
}

func (err ParseCliError) Error() string {
	return fmt.Sprintf("%v: %v", err.time.Format(time.RFC3339), err.reason)
}

func setPortFlag(client *Client) {
	const (
		usage = "Port number for remote server (1-65535)."
	)

	flag.IntVar(&client.serverPort, "port", DefaultPort, usage)
	flag.IntVar(&client.serverPort, "p", DefaultPort, usage+"(shorthand)")
}

func setHostFlag(client *Client) {
	const (
		usage = "IP address of remote server."
	)

	flag.StringVar(&client.serverHost, "host", DefaultHost, usage)
	flag.StringVar(&client.serverHost, "h", DefaultHost, usage+"(shorthand)")
}

func setFlags(client *Client) {
	setPortFlag(client)
	setHostFlag(client)
}

func validatePort(client *Client) error {
	if client.serverPort < 1 || client.serverPort > 65535 {
		return ParseCliError{time.Now(), fmt.Sprintf("Wrong port number provided: %v", client.serverPort)}
	}
	return nil
}

func validateHost(client *Client) error {
	ip := net.ParseIP(client.serverHost)
	if ip == nil {
		return ParseCliError{time.Now(), fmt.Sprintf("Wrong host IP provided: %v", client.serverHost)}
	}
	return nil
}

func validateOptions(client *Client) error {
	var err error

	err = validatePort(client)
	if err != nil {
		return err
	}
	err = validateHost(client)
	if err != nil {
		return err
	}
	return nil
}

func ParseOptions(client *Client) error {
	var err error

	setFlags(client)

	flag.Parse()

	err = validateOptions(client)
	if err != nil {
		return err
	}
	return nil
}
