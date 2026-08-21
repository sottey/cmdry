package cmd

import (
	"fmt"

	"github.com/sottey/cmdry/internal/buildinfo"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the Cmdry build version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), buildinfo.Version)
	},
}

func init() { rootCmd.AddCommand(versionCmd) }
