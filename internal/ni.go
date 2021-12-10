package internal

import (
	"fmt"
	"net"
	"strings"
)

type NI struct {
	Iface net.Interface
	IpNet *net.IPNet
}

// String (pretty) representation of NI struct
func (ni NI) String() string {
	out := []string{
		fmt.Sprintf("%s: flags=%s mtu %d", ni.Iface.Name, ni.Iface.Flags, ni.Iface.MTU),
		fmt.Sprintf("\tether %s", ni.Iface.HardwareAddr),
		fmt.Sprintf("\tinet %s netmask 0x%s broadcast %s", ni.IpNet.IP, ni.IpNet.Mask, BroadcastIPv4(ni.IpNet)),
	}
	return strings.Join(out, "\n")
}

// FindNetworkInterfaces returns network IPv4 interfaces and corresponding IP network, this function ignores loopback
func FindNetworkInterfaces() ([]NI, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list of the system's network interfaces: %w", err)
	}

	var out []NI
	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 {
			// interface not up, skip
			continue
		}
		if i.Flags&net.FlagLoopback != 0 {
			// loopback, skip
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			return nil, fmt.Errorf("get list of unicast interface addresses: %w", err)
		}
		if ipv4 := getIPv4(addrs); ipv4 != nil {
			out = append(out, NI{Iface: i, IpNet: ipv4})
		}
	}
	return out, nil
}

// getIPv4 filter out non-IPv4 addresses
func getIPv4(in []net.Addr) *net.IPNet {
	for _, a := range in {
		if ipn, ok := a.(*net.IPNet); ok {
			if ipn.IP.To4() == nil {
				// not ipv4, skip
				continue
			}
			return ipn
		}
	}
	return nil
}
