package internal

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestBroadcastIPv4(t *testing.T) {
	_, netIp, err := net.ParseCIDR("192.168.86.31/24")
	require.NoError(t, err)
	bNetIp := BroadcastIPv4(netIp)
	assert.Equal(t, "192.168.86.255", bNetIp.String())
}

func TestHostMinIPv4(t *testing.T) {
	_, netIp, err := net.ParseCIDR("192.168.86.31/24")
	require.NoError(t, err)
	ip := HostMinIPv4(netIp)
	assert.Equal(t, "192.168.86.1", ip.String())
}

func TestHostMaxIPv4(t *testing.T) {
	_, netIp, err := net.ParseCIDR("192.168.0.31/24")
	require.NoError(t, err)
	ip := HostMaxIPv4(netIp)
	assert.Equal(t, "192.168.0.254", ip.String())
}

func TestCIDRIPs(t *testing.T) {
	_, netIp, err := net.ParseCIDR("192.168.86.31/24")
	require.NoError(t, err)
	ips := CIDRIPs(netIp)
	require.NoError(t, err)
	assert.Equal(t, 254, len(ips))
	assert.Equal(t, "192.168.86.1", ips[0].String())
	assert.Equal(t, "192.168.86.254", ips[len(ips)-1].String())
}

func TestIPs(t *testing.T) {
	fromIP := net.ParseIP("192.168.0.255")
	toIP := net.ParseIP("192.168.1.5")
	ips := IPs(fromIP, toIP)
	assert.Equal(t, 7, len(ips))
	assert.Equal(t, "192.168.0.255", ips[0].String())
	assert.Equal(t, "192.168.1.5", ips[len(ips)-1].String())
}
