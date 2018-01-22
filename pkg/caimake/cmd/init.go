package cmd

import (
	"os"

	"github.com/caicloud/nirvana/cli"
	"github.com/spf13/cobra"
)

const (
	// EnvProjectType ...
	EnvProjectType = "PROJECT_TYPE"
)

// NewCmdInit ...
func NewCmdInit() *cli.Command {
	cmd := cli.NewCommand(&cobra.Command{
		Use:   "init",
		Short: "init a new Makefile for the project",
		Run: func(cmd *cobra.Command, args []string) {
			os.Setenv(EnvProjectType, cli.GetString("type"))
			bash.Run("install.sh", args)
		},
	})
	cmd.AddFlag(
		cli.StringFlag{
			Name:     "type",
			EnvKey:   EnvProjectType,
			DefValue: "go",
			Usage:    "specify project type",
		},
	)
	return cmd
}
