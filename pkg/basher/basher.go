package basher

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/caicloud/nirvana/log"
	homedir "github.com/mitchellh/go-homedir"
)

// Basher is a context for running bash scripts
type Basher struct {
	// Library path
	Library string
	// library version
	Version string
	// SHA256 checksum
	SHA256sum string
	// Project Name
	Name string
	// full path combines library and name
	FullPath string
	// The io.Reader given to Bash for STDIN
	Stdin io.Reader
	// The io.Writer given to Bash for STDOUT
	Stdout io.Writer
	// The io.Writer given to Bash for STDERR
	Stderr io.Writer

	bindata Bindata
}

// NewBasher returns a new basher context
//   - library is the basher environment path, e.g. ~/.caimake
//   - name is the project name, e.g. make-rules
// The scripts will be placed in library/name, e.g. ~/.caimake/make-rules
func NewBasher(library, name, version string, bindata Bindata) (*Basher, error) {
	path, err := absPath(library)
	if err != nil {
		return nil, err
	}
	c := &Basher{
		Library:  path,
		Version:  version,
		Name:     name,
		FullPath: filepath.Join(path, name),
		Stdin:    os.Stdin,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
		bindata:  bindata,
	}

	b, err := bindata.Asset(filepath.Join(name, "sha256sum"))
	if err != nil {
		return nil, err
	}

	c.SHA256sum = strings.TrimSpace(string(b))

	err = c.ensureAssets()
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Run a command in Bash from this Basher. Standard I/O by
// default is attached to the calling process I/O. You can change this by setting
// the Stdout, Stderr, Stdin variables of the Basher.
func (c *Basher) Run(script string, args []string) (int, error) {

	argstring := ""
	for _, arg := range args {
		argstring = argstring + " '" + strings.Replace(arg, "'", "'\\''", -1) + "'"
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	bash := filepath.Join(c.FullPath, script)
	cmd := exec.Command("bash", "-c", bash+argstring)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Start()
	go func() {
		for sig := range signals {
			cmd.Process.Signal(sig)
			if cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
				cmd.Process.Signal(sig)
			}
		}
	}()
	return exitStatus(cmd.Wait())
}

func (c *Basher) ensureAssets() error {
	err := os.MkdirAll(c.FullPath, 0755)
	if err != nil {
		return err
	}

	sha256sum := filepath.Join(c.FullPath, "sha256sum")
	restore := false
	data, err := ioutil.ReadFile(sha256sum)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// not exists
		restore = true

	} else if strings.TrimSpace(string(data)) != c.SHA256sum {
		restore = true
	}

	if restore {
		c.bindata.RestoreAssets(c.Library, c.Name)
		log.Infof("Restore scripts on %v", c.FullPath)
		log.Infof("Version %v, SHA256 checksum %v", c.Version, c.SHA256sum)
	}
	return nil
}

func absPath(dir string) (string, error) {
	path, err := homedir.Expand(dir)
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

func exitStatus(err error) (int, error) {
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// There is no platform independent way to retrieve
			// the exit code, but the following will work on Unix
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return int(status.ExitStatus()), nil
			}
		}
		return 0, err
	}
	return 0, nil
}
