package shell

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/briandowns/spinner"
)

type Executor interface {
	Execute(command string) (string, error)
}

type shellExecutor struct {
	ctx context.Context
}

func (e *shellExecutor) Execute(command string) (string, error) {

	out := strings.Builder{}
	for _, v := range strings.Split(command, "\n") {
		cout, err := exec.Command(SHELL, EXEC, v).CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("[%w] \n %s", err, cout)
		}
		_, err = out.WriteString(string(cout))
		if err != nil {
			return "", err
		}
	}
	return out.String(), nil
}

func NewExecutor(ctx context.Context) Executor {
	return &shellExecutor{
		ctx: ctx,
	}
}

func VerboseExecutor(stdOut, errOut io.Writer, exec Executor) Executor {
	return &ShellWritter{
		Out:      stdOut,
		Err:      errOut,
		Executor: exec,
	}
}

type ShellWritter struct {
	Out      io.Writer
	Err      io.Writer
	Executor Executor
}

func (s *ShellWritter) Execute(command string) (string, error) {
	sp := spinner.New(spinner.CharSets[78], 100*time.Millisecond)
	spErr := sp.Color("blue")
	if spErr == nil {
		sp.Suffix = color.InBlue(" cmd: ") + command
		sp.Start()
	}

	v, execErr := s.Executor.Execute(command)
	if spErr == nil {
		sp.Stop()
	}
	if execErr != nil {
		fmt.Fprintln(s.Out, color.InRed(command))
		fmt.Fprintln(s.Err, color.InRed("error: "), execErr.Error())
	}
	fmt.Fprintln(s.Out, color.InGreen("✅︎"+command))
	return v, execErr
}
