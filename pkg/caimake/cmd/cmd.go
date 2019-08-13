package cmd

import (
	"github.com/caicloud/build-infra/pkg/basher"
	"github.com/caicloud/nirvana/cli"
	"github.com/caicloud/nirvana/log"
	"github.com/spf13/cobra"
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
			b, err := bindata.Asset("caimake.sh")
			if err != nil {
				log.Fatal("Error load script caimake.sh")
				return err
			}
			bash, err = basher.NewBasher("caimake", b, bindata)
			if err != nil {
				log.Fatal(err)
				return err
			}
			return nil
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
