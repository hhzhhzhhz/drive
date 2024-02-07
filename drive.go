package drive

import (
	"github.com/hhzhhzhhz/drive/tun"
	"io"
	"time"
)

var (
	timeout = 2 * time.Second
)

type Drive interface {
	io.ReadWriteCloser
	Tun() tun.Tun
	Up(cfg interface{}) error
}
