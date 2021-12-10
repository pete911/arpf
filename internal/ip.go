package internal

import (
	"encoding/binary"
	"net"
)

func BroadcastIPv4(ip *net.IPNet) net.IP {
	maskOnes, _ := ip.Mask.Size()
	mask := uint32(0xFFFFFFFF >> maskOnes)
	out := binary.BigEndian.Uint32(ip.IP.To4()) | mask
	return intToIP(out)
}

func HostMinIPv4(ip *net.IPNet) net.IP {
	maskOnes, _ := ip.Mask.Size()
	mask := uint32(0xFFFFFFFF >> maskOnes)
	// network address
	out := binary.BigEndian.Uint32(ip.IP.To4()) &^ mask
	// increment by one
	out++
	return intToIP(out)
}

func HostMaxIPv4(ip *net.IPNet) net.IP {
	bcastInt := binary.BigEndian.Uint32(BroadcastIPv4(ip).To4())
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
	for i := binary.BigEndian.Uint32(from.To4()); i <= binary.BigEndian.Uint32(to.To4()); i++ {
		b := make([]byte, 4)
		binary.BigEndian.PutUint32(b, i)
		out = append(out, b)
	}
	return out
}

func intToIP(i uint32) net.IP {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}
