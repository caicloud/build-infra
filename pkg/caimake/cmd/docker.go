package cmd

import (
	"github.com/caicloud/nirvana/cli"
	"github.com/spf13/cobra"
)

// NewCmdDocker ...
func NewCmdDocker() *cli.Command {
	cmd := cli.NewCommand(&cobra.Command{
		Use:   "docker",
		Short: "make for docker",
		Run:   runHelp,
	})
	cmd.AddCommand(
		newCmdDockerBuild(),
		newCmdDockerPush(),
	)
	return cmd
}

func newCmdDockerBuild() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "build",
		Short: "docker build",
		Run: func(cmd *cobra.Command, args []string) {
			bash.Run2("docker", "build", args...)
		},
	})
}

func newCmdDockerPush() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "push",
		Short: "docker push",
		Run: func(cmd *cobra.Command, args []string) {
			bash.Run2("docker", "push", args...)
		},
	})
}
