package cmd

import (
	"breakbio-openvax/pkg/breakbio"
	"github.com/spf13/cobra"
)

var runParams = breakbio.RunParams{}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringVarP(&runParams.Shell, "shell", "s", "", "Shell for execution.")
	err := runCmd.MarkFlagRequired("shell")
	handleRequiredFlagErrors(err)

	runCmd.Flags().StringVarP(&runParams.Command, "command", "c", "", "Command to be executed.")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run command.",
	Long:  "Run command in OS shell.",
	Run: func(cmd *cobra.Command, args []string) {
		breakbio.RunCommand(runParams)
	},
}
