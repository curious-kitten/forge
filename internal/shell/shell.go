package shell

import (
	"context"
	"io"
	"os/exec"
	"strings"
)

type Executor struct {
	Ctx context.Context
	Out io.Writer
	Err io.Writer
}

func (e *Executor) Execute(command string) error {
	for _, v := range strings.Split(command, "\n") {
		cmd := exec.Command(SHELL, EXEC, v)
		cmd.Stdout = e.Out
		cmd.Stderr = e.Err
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
