package internal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

const (
	scanTime         = 10 * time.Second
	scanIntervalTime = 3 * time.Second
)

type ARPResponse struct {
	Names       []string
	ProtAddress net.IP
	HwAddress   net.HardwareAddr
}

func NewARPResponse(in *layers.ARP) ARPResponse {
	if in == nil {
		return ARPResponse{}
	}
	response := ARPResponse{
		ProtAddress: net.IP(in.SourceProtAddress),
		HwAddress:   net.HardwareAddr(in.SourceHwAddress),
	}
	if names, err := net.LookupAddr(response.ProtAddress.String()); err == nil {
		response.Names = names
	}
	return response
}

func (a ARPResponse) String() string {
	return fmt.Sprintf("%s (%s) at %s", strings.Join(a.Names, " "), a.ProtAddress, a.HwAddress)
}

func Scan(iface net.Interface, src net.IP, dst []net.IP) ([]ARPResponse, error) {
	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, fmt.Errorf("pcap open %s device: %w", iface.Name, err)
	}
	defer handle.Close()

	return readWrite(handle, iface, src, dst)
}

func readWrite(handle *pcap.Handle, iface net.Interface, src net.IP, dst []net.IP) ([]ARPResponse, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), scanTime)
	defer cancelFn()
	arps := readARP(ctx, handle, iface)

	// add src IP to already scanned IPs
	scannedIPs := make(map[string]struct{})
	scannedIPs[src.To4().String()] = struct{}{}
	scanInterval := time.NewTicker(scanIntervalTime)

	// initial arp request
	if err := writeARP(handle, iface, src, dst, scannedIPs); err != nil {
		return nil, fmt.Errorf("write arp on %v: %w", iface.Name, err)
	}

	log.Info().Msgf("sending packets for %s in %s intervals", scanTime, scanIntervalTime)
	var out []ARPResponse
loop:
	for {
		select {
		case arp, ok := <-arps:
			if !ok {
				log.Debug().Msg("arps channel closed")
				break loop
			}
			key := net.IP(arp.SourceProtAddress).String()
			log.Debug().Msgf("read arp dst %s src %s", net.IP(arp.DstProtAddress), key)
			if _, ok := scannedIPs[key]; !ok {
				log.Debug().Msgf("arp src %s not saved yet", key)
				scannedIPs[key] = struct{}{}
				out = append(out, NewARPResponse(arp))
			}
		case <-scanInterval.C:
			if err := writeARP(handle, iface, src, dst, scannedIPs); err != nil {
				return nil, fmt.Errorf("write arp on %v: %w", iface.Name, err)
			}
		}
	}
	return out, nil
}

func readARP(ctx context.Context, handle *pcap.Handle, iface net.Interface) <-chan *layers.ARP {
	arps := make(chan *layers.ARP)

	go func() {
		src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
		in := src.Packets()
		for {
			var packet gopacket.Packet
			select {
			case <-ctx.Done():
				log.Debug().Msg("closing arps channel")
				close(arps)
				return
			case packet = <-in:
				arpLayer := packet.Layer(layers.LayerTypeARP)
				if arpLayer == nil {
					continue
				}
				arp, ok := arpLayer.(*layers.ARP)
				if !ok {
					continue
				}
				if arp.Operation != layers.ARPReply || bytes.Equal(iface.HardwareAddr, arp.SourceHwAddress) {
					continue
				}
				arps <- arp
			}
		}
	}()
	return arps
}

func writeARP(handle *pcap.Handle, iface net.Interface, src net.IP, dst []net.IP, scannedIPs map[string]struct{}) error {
	var count int
	for _, ip := range dst {
		if _, ok := scannedIPs[ip.String()]; ok {
			continue
		}
		buf, err := getArpBuffer(iface.HardwareAddr, src, ip)
		if err != nil {
			return fmt.Errorf("get arp buffer: %w", err)
		}
		if err := handle.WritePacketData(buf.Bytes()); err != nil {
			return fmt.Errorf("write packet data: %w", err)
		}
		count++
	}
	log.Debug().Msgf("send arp to %d IPs", count)
	return nil
}

func getArpBuffer(srcMAC net.HardwareAddr, srcIP, dstIP net.IP) (gopacket.SerializeBuffer, error) {
	eth := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(srcMAC),
		SourceProtAddress: []byte(srcIP.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		DstProtAddress:    []byte(dstIP.To4()),
	}
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
		return nil, fmt.Errorf("serialize layer: %w", err)
	}
	return buf, nil
}
