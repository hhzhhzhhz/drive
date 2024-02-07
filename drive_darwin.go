//go:build darwin
// +build darwin

package drive

import (
	"fmt"
	"github.com/hhzhhzhhz/drive/tun"
	"github.com/hhzhhzhhz/gopkg/utils"
	"net"
	"time"
)

var (
	ifc = "ifconfig"
)

func NewDrive(name, subnet string) (Drive, error) {
	dev, err := tun.CreateTUN(name, 0)
	if err != nil {
		return nil, err
	}
	return &darwinDrive{drive: dev}, nil
}

// tunDrive Network Interface Card (NIC)
type darwinDrive struct {
	subnet string
	drive  tun.Tun
}

func (d *darwinDrive) Write(data []byte) (n int, err error) {
	return d.drive.Write(data)
}

func (d *darwinDrive) Read(data []byte) (n int, err error) {
	return d.drive.Read(data)
}

func (d *darwinDrive) Close() error {
	return d.drive.Close()
}

func (d *darwinDrive) Tun() tun.Tun {
	return d.drive
}

func (d *darwinDrive) Up(cfg interface{}) error {
	serverIp, ok := cfg.(string)
	if !ok {
		return fmt.Errorf("cfg invalid parameter")
	}
	ip, _, err := net.ParseCIDR(d.subnet)
	if err != nil {
		return err
	}
	dName, err := d.drive.Name()
	if err != nil {
		return err
	}
	timeout := 2 * time.Second
	_, err = utils.Command(ifc, timeout, dName, "inet", ip.String(), serverIp, "up")
	// utils.Command(ifc, timeout, devName, "inet6", "v6", "ipv6", "up")
	return err
}

// Route action: add,change
func (d *darwinDrive) Route(action, ip string, gateway string) error {
	_, err := utils.Command("route", timeout, action, ip, gateway)
	return err
}
