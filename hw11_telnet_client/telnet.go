package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrClosedConnection = errors.New("connection is closed")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *telnetClient) Connect() error {
	var err error
	tc.conn, err = net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	return nil
}

func (tc *telnetClient) Send() error {
	if tc.conn == nil {
		return ErrClosedConnection
	}

	if _, err := io.Copy(tc.conn, tc.in); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	return nil
}

func (tc *telnetClient) Receive() error {
	if tc.conn == nil {
		return ErrClosedConnection
	}

	if _, err := io.Copy(tc.out, tc.conn); err != nil {
		return fmt.Errorf("failed to read: %w", err)
	}
	return nil
}

func (tc *telnetClient) Close() error {
	if tc.conn == nil {
		return ErrClosedConnection
	}

	return tc.conn.Close()
}
