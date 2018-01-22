package cmd

import (
	"github.com/caicloud/build-infra/pkg/basher"
	"github.com/caicloud/build-infra/pkg/version"
	"github.com/caicloud/nirvana/cli"
	"github.com/caicloud/nirvana/log"
	"github.com/spf13/cobra"
)

const (
	dotfile = "~/.caimake"
	project = "make-rules"
)

var (
	bash *basher.Basher
)

// NewCaimakeCommand returns a new caimake command
func NewCaimakeCommand() *cli.Command {
	cmd := cli.NewCommand(&cobra.Command{
		Use:   "caimake",
		Short: `caimake helps you to build your project appropriately`,
		Long: `caimake helps you to build your project appropriately

It follows the specification on Github caicloud/build-infra.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			bindata := basher.NewBindata(
				Asset,
				MustAsset,
				AssetInfo,
				AssetNames,
				AssetDir,
				RestoreAsset,
				RestoreAssets,
			)
			bash, err = basher.NewBasher(dotfile, project, version.Get().Version, bindata)
			if err != nil {
				log.Fatal(err)
			}
		},
		Run: runHelp,
	})

	cmd.AddCommand(NewCmdInit())
	cmd.AddCommand(NewCmdVersion())
	cmd.AddCommand(NewCmdGolang())
	cmd.AddCommand(NewCmdDocker())
	cmd.AddCommand(NewCmdClean())
	cmd.AddCommand(NewCmdUpdate())
	return cmd
}
