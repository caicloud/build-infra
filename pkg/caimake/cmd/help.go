package cmd

import "github.com/spf13/cobra"

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
