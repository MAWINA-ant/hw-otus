package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	telnetStruct := TelnetStruct{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
	return &telnetStruct
}

type TelnetStruct struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	con     net.Conn
}

func (t *TelnetStruct) Connect() error {
	con, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.con = con
	return nil
}

func (t *TelnetStruct) Close() error {
	if t.con != nil {
		return t.con.Close()
	}
	return nil
}

func (t *TelnetStruct) Send() error {
	_, err := io.Copy(t.con, t.in)
	return err
}

func (t *TelnetStruct) Receive() error {
	_, err := io.Copy(t.out, t.con)
	return err
}
