package executer

import (
	"io"
	"os/exec"
	"path/filepath"
)

func ExecCommandWithName(name string, command []string, w io.Writer) error {

	rebootLimit := 10
	var err error
	for rebootLimit > 0 {
		if err = ExecCommand(command, w); err != nil {
			rebootLimit--
		}
	}
	return err
}

func ExecCommand(command []string, w io.Writer) error {
	ap, err := filepath.Abs(command[0])
	if err != nil {
		return err
	}
	bp := filepath.Dir(ap)
	command[0] = ap
	cmd := &exec.Cmd{
		Path:   ap,
		Dir:    bp,
		Args:   command,
		Stderr: w,
		Stdout: w,
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	defer func() {
		cmd.Process.Kill()
	}()
	return cmd.Wait()
}
