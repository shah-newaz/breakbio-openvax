package cmd

import (
	"breakbio-openvax/pkg/breakbio/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of BreakBio OpenVax Server.",
	Long:  `All software has versions. This is BreakBio OpenVax Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Cyan("BreakBio OpenVax Server v" + VERSION)
	},
}
