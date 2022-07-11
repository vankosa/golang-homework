package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &NetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type NetClient struct {
	address string
	conn    net.Conn
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	endMsg  string
}

func (c *NetClient) Connect() error {
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)
	return err
}

func (c *NetClient) Send() error {
	content := make([]byte, 100)

	// get data from stdin
	_, err := bufio.NewReader(c.in).Read(content)

	// check EOF error
	if errors.Is(err, io.EOF) {
		c.endMsg = "...EOF"
		return c.Close()
	}
	// try to write to connection
	_, err = c.conn.Write(content)

	// detect nc close connection
	if _, ok := err.(*net.OpError); ok {
		return c.Close()
	}

	return err
}

func (c *NetClient) Receive() error {
	_, err := io.Copy(c.out, c.conn)

	return err
}

func (c *NetClient) Close() error {
	// set final word
	endMsg := "...Connection was closed by peer"
	if c.endMsg != "" {
		endMsg = c.endMsg
	}

	// send final word to stderr
	_, err := fmt.Fprintf(os.Stderr, endMsg)
	if err != nil {
		return err
	}

	return c.conn.Close()
}
