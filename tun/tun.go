package tun

import (
	"io"
)

type Event int

const (
	EventUp = 1 << iota
	EventDown
	EventMTUUpdate
)

type Tun interface {
	io.ReadWriteCloser
	Name() (string, error)
}
