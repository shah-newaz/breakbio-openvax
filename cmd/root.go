package cmd

import (
	"breakbio-openvax/pkg/breakbio/log"
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var (
	// VERSION is set during build
	VERSION string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "breakbio-openvax",
	Short: "Tool to run breakbio openvax pipeline.",
	Long:  `CLI for breakbio openvax pipeline`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute(version string) {
	VERSION = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func handleRequiredFlagErrors(err error) {
	if err != nil {
		log.Red("Required flags are missing... " + err.Error())
		os.Exit(126)
	}
}
