//go:build linux
// +build linux

package drive

import (
	"fmt"
	"github.com/hhzhhzhhz/drive/tun"
	"github.com/hhzhhzhhz/gopkg/utils"
	"go.uber.org/multierr"
	"net"
	"strconv"
)

func NewDrive(name, subnet string) (Drive, error) {
	_, _, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("error cidr %s", err.Error())
	}
	dev, err := tun.CreateTUN(name, 1500)
	if err != nil {
		return nil, err
	}
	return &linuxDrive{drive: dev, subnet: subnet}, nil
}

// tunDrive Network Interface Card (NIC)
type linuxDrive struct {
	subnet string
	drive  tun.Tun
}

func (l *linuxDrive) Write(data []byte) (n int, err error) {
	return l.drive.Write(data)
}

func (l *linuxDrive) Read(data []byte) (n int, err error) {
	return l.drive.Read(data)
}

func (l *linuxDrive) Close() error {
	return l.drive.Close()
}

func (l *linuxDrive) Tun() tun.Tun {
	return l.drive
}

func (l *linuxDrive) Up(cfg interface{}) error {
	dName, err := l.drive.Name()
	if err != nil {
		return err
	}
	_, cer := utils.Command("/sbin/ip", timeout, "link", "set", "dev", dName, "mtu", strconv.Itoa(1500))
	err = multierr.Append(err, cer)
	_, cer = utils.Command("/sbin/ip", timeout, "addr", "add", l.subnet, "dev", dName)
	err = multierr.Append(err, cer)
	// sys.Command("/sbin/ip", timeout, "-6", "addr", "add", "fced:9999::9999/64", "dev", dName)
	_, cer = utils.Command("/sbin/ip", timeout, "link", "set", "dev", dName, "up")
	err = multierr.Append(err, cer)
	return err
}

// Route action: add
func (l *linuxDrive) Route(action, ip string) error {
	dName, err := l.drive.Name()
	if err != nil {
		return err
	}
	_, err = utils.Command("/sbin/ip", timeout, "route", action, ip, "dev", dName)
	return err
}
