package internal

import (
	"encoding/binary"
	"net"
)

func BroadcastIPv4(ip *net.IPNet) net.IP {
	out := ipToInt(ip.IP) | getMask(ip)
	return intToIP(out)
}

func HostMinIPv4(ip *net.IPNet) net.IP {
	// network address
	out := ipToInt(ip.IP) &^ getMask(ip)
	// increment by one
	out++
	return intToIP(out)
}

func HostMaxIPv4(ip *net.IPNet) net.IP {
	bcastInt := ipToInt(BroadcastIPv4(ip))
	bcastInt--
	return intToIP(bcastInt)
}

// CIDRIPs returns list of IPs not including network and broadcast IP
func CIDRIPs(ip *net.IPNet) []net.IP {
	return IPs(HostMinIPv4(ip), HostMaxIPv4(ip))
}

// IPs returns list of IPs including from and to arguments
func IPs(from net.IP, to net.IP) []net.IP {
	var out []net.IP
	for i := ipToInt(from); i <= ipToInt(to); i++ {
		out = append(out, intToIP(i))
	}
	return out
}

func getMask(ip *net.IPNet) uint32 {
	return 0xFFFFFFFF ^ binary.BigEndian.Uint32(ip.Mask)
}

func intToIP(i uint32) net.IP {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func ipToInt(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}
