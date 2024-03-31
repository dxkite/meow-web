package executer

import (
	"io"
	"os/exec"
	"path/filepath"

	"dxkite.cn/log"

	"dxkite.cn/meownest/src/utils"
)

func ExecInstance(name string, command []string) error {
	w := utils.MakeNameLoggerWriter(name)
	rebootLimit := 10
	var err error
	for rebootLimit > 0 {
		if err = ExecCommand(command, w); err != nil {
			log.Error("exec error, reboot", command, err)
			rebootLimit--
		}
	}
	return err
}

func ExecCommand(command []string, w io.Writer) error {
	ap, err := filepath.Abs(command[0])
	if err != nil {
		log.Error("exec", command, err)
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
		log.Error("exec", command, err)
		return err
	}

	log.Info("exec", command, "pid", cmd.Process.Pid)

	defer func() {
		cmd.Process.Kill()
	}()
	return cmd.Wait()
}
