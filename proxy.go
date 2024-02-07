package drive

import (
	"github.com/hhzhhzhhz/gopkg/log"
	"io"
	"net"
	"runtime"
	"strings"
)

func NewProxy(l log.Log) *Proxy {
	return &Proxy{
		logger: l,
	}
}

// Proxy forwards a TCP request to a TCP service.
type Proxy struct {
	logger log.Log
}

func (p *Proxy) Copy(dst, src io.ReadWriteCloser) error {
	defer func() {
		if errRecover := recover(); errRecover != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			log.Logger().Error("an error occurred while copying panic error:%v; stack:%s", errRecover, buf[:n])
		}
	}()
	go p.copy(dst, src)
	p.copy(src, dst)
	return nil
}

func (p Proxy) copy(dst io.WriteCloser, src io.Reader) {
	_, err := io.Copy(dst, src)
	if err != nil {
		p.logger.Error("src read failed err_msg=%s", err.Error())
	}
}

func TCPServer(server net.Listener, handle func(conn net.Conn)) {
	for {
		conn, err := server.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				log.Logger().Warn("temporary accept err=%s", err)
				if conn != nil {
					conn.Close()
				}
				runtime.Gosched()
				continue
			}

			// theres no direct way to detect this error because it is no exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Logger().Error("listener.Accept() - %s", err)
			}

			if conn != nil {
				conn.Close()
			}
			break
		}
		go handle(conn)
	}
}
