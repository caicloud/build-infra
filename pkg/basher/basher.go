package basher

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// Basher is a context for running bash scripts
type Basher struct {
	name    string
	content []byte
	// The io.Reader given to Bash for STDIN
	Stdin io.Reader
	// The io.Writer given to Bash for STDOUT
	Stdout io.Writer
	// The io.Writer given to Bash for STDERR
	Stderr io.Writer

	Bindata Bindata
}

// NewBasher returns a new basher context
//   - bindata contains all the assets
//   - scriptContent is the all in one script content
// The scripts will be placed in library/name, e.g. ~/.caimake/make-rules
func NewBasher(name string, scriptContent []byte, bindata Bindata) (*Basher, error) {
	c := &Basher{
		name:    name,
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Bindata: bindata,
		content: scriptContent,
	}

	return c, nil
}

// Run1 runs one command
func (c *Basher) Run1(cmd string, args ...string) (int, error) {
	narg := []string{cmd}
	narg = append(narg, args...)
	return c.Run(narg...)
}

// Run2 runs command and subcommand
func (c *Basher) Run2(cmd1, cmd2 string, args ...string) (int, error) {
	narg := []string{cmd1, cmd2}
	narg = append(narg, args...)
	return c.Run(narg...)
}

// Run a command in Bash from this Basher. Standard I/O by
// default is attached to the calling process I/O. You can change this by setting
// the Stdout, Stderr, Stdin variables of the Basher.
func (c *Basher) Run(args ...string) (int, error) {

	argstring := ""
	for _, arg := range args {
		argstring = argstring + " '" + strings.Replace(arg, "'", "'\\''", -1) + "'"
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	envfile, err := c.buildTemp()
	if err != nil {
		return 0, fmt.Errorf("Error create %v temporary script", c.name)
	}
	cmd := exec.Command("bash", "-c", envfile+argstring)
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

func (c *Basher) buildTemp() (string, error) {
	file, err := ioutil.TempFile(os.TempDir(), c.name+".")
	if err != nil {
		return "", err
	}
	defer file.Close()
	file.Write(c.content)
	err = os.Chmod(file.Name(), 0755)
	if err != nil {
		return "", err
	}
	return file.Name(), nil
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
