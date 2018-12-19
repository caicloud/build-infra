package cmd

import (
	"github.com/Masterminds/semver"
	"github.com/caicloud/build-infra/pkg/update"
	"github.com/caicloud/build-infra/pkg/version"
	"github.com/caicloud/nirvana/cli"
	"github.com/caicloud/nirvana/log"
	"github.com/spf13/cobra"
)

// NewCmdUpdate ...
func NewCmdUpdate() *cli.Command {
	cmd := cli.NewCommand(&cobra.Command{
		Use:   "update",
		Short: "update caimake",
		RunE: func(cmd *cobra.Command, args []string) error {
			latest, err := update.GetGithubLatestRelease(cli.GetString("repo"))
			if err != nil {
				log.Fatal(err)
				return err
			}

			current := version.Get().Version

			if cli.GetBool("check") {
				log.Infof("Current version: %v", current)
				log.Infof("Latest  version: %v", latest.Version)
				return nil
			}

			latestSemver, err := semver.NewVersion(latest.Version)
			if err != nil {
				log.Fatalf("latest version [%v] does not follow semantic version, %v", latest.Version, err)
				return err
			}
			currentSemver, err := semver.NewVersion(current)
			if err != nil {
				log.Fatalf("current version [%v] does not follow semantic version, %v", current, err)
				return err
			}

			if currentSemver.GreaterThan(latestSemver) || currentSemver.Equal(latestSemver) {
				log.Infof("The binary is up to date!")
				log.Infof("Caimake %v is currently the newest version available.", current)
				return nil
			}
			err = update.DoUpdate(latest.DownloadURL, latest.Size)
			if err != nil {
				log.Fatal(err)
				return err
			}
			return nil
		},
	})

	cmd.AddFlag(
		cli.BoolFlag{
			Name:      "check",
			Shorthand: "c",
			Usage:     "Check latest version only",
		},
		cli.StringFlag{
			Name:      "repo",
			Shorthand: "r",
			Usage:     "Github repo name",
			DefValue:  "caicloud/build-infra",
		},
	)

	return cmd
}
