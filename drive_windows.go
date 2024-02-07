//go:build windows
// +build windows

package drive

import (
	"github.com/hhzhhzhhz/drive/tun"
	"github.com/hhzhhzhhz/gopkg/utils"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
	"net/netip"
)

func NewDrive(name, subnet string) (Drive, error) {
	//id := &windows.GUID{
	//	0x0000000,
	//	0xFFFF,
	//	0xFFFF,
	//	[8]byte{0xFF, 0xe9, 0x76, 0xe5, 0x8c, 0x74, 0x06, 0x3e},
	//}
	dev, err := tun.CreateTUN(name, 0)
	if err != nil {
		return nil, err
	}
	// 保存原始设备句柄
	nativeTunDevice := dev.(*tun.NativeTun)

	// 获取LUID用于配置网络
	link := winipcfg.LUID(nativeTunDevice.LUID())
	// 设置网卡网络
	addr, err := netip.ParsePrefix(subnet)
	if err != nil {
		return nil, err
	}
	err = link.SetIPAddresses([]netip.Prefix{addr})
	if err != nil {
		return nil, err
	}
	return &winDrive{drive: dev}, nil
}

// tunDrive Network Interface Card (NIC)
type winDrive struct {
	drive tun.Tun
}

func (w *winDrive) Write(data []byte) (int, error) {
	return w.drive.Write(data)
}

func (w *winDrive) Read(data []byte) (n int, err error) {
	return w.drive.Read(data)
}

func (w *winDrive) Close() error {
	return w.drive.Close()
}

func (w *winDrive) Tun() tun.Tun {
	return nil
}

func (w *winDrive) Up(cfg interface{}) error {
	return nil
}

// Route action: add,delete
func (w *winDrive) Route(action, ip, mask, gateway string) error {
	_, err := utils.Command("cmd", timeout, "/C", "route", action, ip, "mask", mask, gateway)
	return err
}
