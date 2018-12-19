package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/caicloud/nirvana/cli"
	"github.com/caicloud/nirvana/log"
	"github.com/spf13/cobra"
)

const ()

// NewCmdInit ...
func NewCmdInit() *cli.Command {
	cmd := cli.NewCommand(&cobra.Command{
		Use:   "init",
		Short: "init a new Makefile for the project",
		Run:   runHelp,
	})
	cmd.AddFlag(
		cli.BoolFlag{
			Name:       "override",
			Usage:      "override already exists Makefile and .caimake",
			Persistent: true,
		},
		cli.BoolFlag{
			Name:       "offline",
			Usage:      "install full scripts for offline environment",
			Persistent: true,
		},
		cli.StringSliceFlag{
			Name:       "registries",
			Usage:      "docker registries",
			DefValue:   []string{"cargo.caicloudprivatetest.com/caicloud"},
			Persistent: true,
		},
	)

	cmd.AddCommand(newCmdInitGo())
	return cmd
}

func newCmdInitGo() *cli.Command {
	return cli.NewCommand(&cobra.Command{
		Use:   "go",
		Short: "init a new golang project",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatalf("Error get current work dir, %v", err)
				return err
			}
			err = restoreCaimake(dir)
			if err != nil {
				log.Fatal(err)
				return err
			}
			cmds, builds, err := detectGoCmdAndDocker(dir)
			if err != nil {
				log.Fatalf("Error get sub dirs of builds, %v", err)
				return err
			}
			return restoreMakefile(dir, "go", map[string]interface{}{
				"cmds":   cmds,
				"builds": builds,
			})
		},
	})
}

func restoreCaimake(dir string) error {
	caimakePath := filepath.Join(dir, ".caimake")

	_, err := os.Stat(caimakePath)
	if err == nil && !cli.GetBool("override") {
		// exists
		return errors.New(".caimake already exists, use --override to cover it")
	}

	// restore update.sh to ./.caimake/update.sh
	err = bash.Bindata.RestoreAsset(caimakePath, "update.sh")
	if err != nil {
		return fmt.Errorf("Error restore update.sh to ./.caimake/update.sh, %v", err)
	}

	if cli.GetBool("offline") {
		// restore caimake.sh to ./.caimake/caimake.sh
		err = bash.Bindata.RestoreAsset(caimakePath, "caimake.sh")
		if err != nil {
			return fmt.Errorf("Error restore caimake.sh to ./.caimake/caimake.sh, %v", err)
		}
	}

	log.Info("Restore .caimake successfully")
	return nil
}

func restoreMakefile(dir, lang string, additional map[string]interface{}) error {
	_, err := os.Stat(filepath.Join(dir, "Makefile"))
	if err == nil && !cli.GetBool("override") {
		// exists
		return errors.New("Makefile already exists, use --override to cover it")
	}

	// load template
	b, err := bash.Bindata.Asset(fmt.Sprintf("Makefile.%v.tmpl", lang))
	if err != nil {
		return fmt.Errorf("Error load golang Makefile template, %v", err)
	}

	tmpl, err := template.New("Makefile").Parse(string(b))
	if err != nil {
		return fmt.Errorf("Error parse golang Makefile template, %v", err)
	}

	m := map[string]interface{}{
		"dockerImagePrefix": filepath.Base(dir) + "-",
		"registries":        cli.GetStringSlice("registries"),
	}

	for k, v := range additional {
		m[k] = v
	}

	makefile, err := os.OpenFile(filepath.Join(dir, "Makefile"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Error open Makefile, %v", err)
	}

	tmpl.Execute(makefile, m)

	log.Info("Generate Makefile for golang project successfully")
	return nil
}

func subDirs(dir string) ([]string, error) {
	dirs := []string{}

	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return dirs, nil
	}
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("%v is not dir", dir)
	}

	first := true
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if first {
			first = false
			return nil
		}
		if err != nil {
			return err
		}

		if info.IsDir() {
			base := filepath.Base(path)
			dirs = append(dirs, base)
			return filepath.SkipDir
		}
		return nil
	})

	return dirs, err
}

func addprefix(prefix string, l []string) []string {
	ret := make([]string, len(l))
	for i, v := range l {
		ret[i] = prefix + v
	}
	return ret
}

func detectGoCmdAndDocker(dir string) ([]string, []string, error) {
	// detect cmds and builds
	cmds, err := subDirs(filepath.Join(dir, "cmd"))
	if err != nil {
		// log.Fatalf("Error get sub dirs of cmds, %v", err)
		return nil, nil, err
	}
	builds, err := subDirs(filepath.Join(dir, "build"))
	if err != nil {
		// log.Fatalf("Error get sub dirs of builds, %v", err)
		return nil, nil, err
	}
	return addprefix("cmd/", cmds), addprefix("build/", builds), nil
}
