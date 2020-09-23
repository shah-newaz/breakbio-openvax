package cmd

import (
	"breakbio-openvax/pkg/breakbio"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serve)
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Start BreakBio-OpenVax server.",
	Long:  "Start BreakBio-OpenVax server on container to receive commands.",
	Run: func(cmd *cobra.Command, args []string) {
		breakbio.Serve()
	},
}
