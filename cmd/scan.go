package cmd

import (
	"github.com/pete911/arpf/internal"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "list MAC and IPs on the network",
		Run:   scanRun,
	}
)

func scanRun(_ *cobra.Command, _ []string) {

	nis, err := internal.FindNetworkInterfaces()
	if err != nil {
		log.Fatal().Err(err).Msg("find network interfaces")
	}

	for _, ni := range nis {
		src := ni.IpNet.IP
		dst := internal.CIDRIPs(ni.IpNet)
		log.Info().Msg(ni.String())
		log.Info().Msgf("sending ARPs: src %s dst %s ... %s", src, dst[0], dst[len(dst)-1])
		arps, err := internal.Scan(ni.Iface, src, dst)
		if err != nil {
			log.Error().Err(err).Msg("arp scan")
		}
		for _, arp := range arps {
			log.Info().Msg(arp.String())
		}
	}
}
