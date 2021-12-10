package cmd

import (
	"github.com/spf13/cobra"
)

var (
	verboseFlag bool

	RootCmd = &cobra.Command{
		Use:   "arpf",
		Short: "arp utils",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initLog(verboseFlag)
		},
	}
)

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "print debug messages")
	RootCmd.AddCommand(scanCmd)
}
