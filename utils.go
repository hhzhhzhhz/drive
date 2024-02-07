package drive

import (
	"fmt"
	"github.com/hhzhhzhhz/gopkg/utils"
	"net"
	"runtime"
	"strings"
	"time"
)

func Ipv4Gateway() (net.IP, error) {
	timeout := 2 * time.Second
	switch runtime.GOOS {
	case "linux":
		ip, err := utils.Command("sh", timeout, "-c", "route -n | grep 'UG[ \t]' | awk 'NR==1{print $2}'")
		if err != nil {
			return nil, err
		}
		ipv4 := net.ParseIP(ip)
		if ipv4 == nil {
			return nil, err
		}
		return ipv4, nil
	case "darwin":
		ip, err := utils.Command("sh", timeout, "-c", "route -n get default | grep 'gateway' | awk 'NR==1{print $2}'")
		if err != nil {
			return nil, err
		}
		ipv4 := net.ParseIP(ip)
		if ipv4 == nil {
			return nil, err
		}
		return ipv4, nil
	case "windows":
		output, err := utils.Command("route", timeout, "print", "0.0.0.0")
		if err != nil {
			return nil, err
		}
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) > 2 && fields[0] == "0.0.0.0" {
				return net.ParseIP(fields[2]), nil
			}
		}
		return nil, fmt.Errorf("unable to find default gateway")
	}

	return nil, fmt.Errorf("%s system does not support", runtime.GOOS)
}
