package run

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/cruious-kitten/forge/internal/shell"
	"github.com/cruious-kitten/forge/pkg/forge"
)

type forgeryReader func() (forge.Forge, error)

func argProvider(fr forgeryReader) []string {
	f, err := fr()
	if err != nil {
		return []string{}
	}
	return f.GetForgeries()
}

const longDescription = `
A forgery is a sequenctial execution of multiple tools. 
Outputs of tool execution can be used as arguments to other tools inside the forgery.

Example:

tools:
- name: build
  args:
    - name: version
      default: dev
    - name: hash
      default: 123#456
    - name: date
  description: Go install the local package
  cmd: go build -v -ldflags="-X 'github.com/cruious-kitten/forge/internal/info.Version={{.version}}' -X 'github.com/cruious-kitten/forge/internal/info.CommitHash={{.hash}}' -X 'github.com/cruious-kitten/forge/internal/info.BuildDate={{.date}}'" .
- name: git-tag
  cmd: git describe --tags --dirty --always
- name: build-date
  cmd: date +%FT%T%z
- name: hash
  cmd: git rev-parse --short HEAD 2>/dev/null
forgeries:
  - name: release
    tools:
    - name: build-date
      output: date
    - name: hash
      output: hash
    - name: git-tag
      output: tag
    - name: build
      args:
        version: output.tag
        date: output.date
        hash: output.hash

`

func Command(fr forgeryReader) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "run [forgery]",
		Short:     "execute a forgery",
		Long:      longDescription,
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: argProvider(fr),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := fr()
			if err != nil {
				return err
			}
			exec := shell.VerboseExecutor(os.Stdout, os.Stderr, shell.NewExecutor(context.Background()))

			err = f.RunForgery(exec, args[0])
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			return nil
		},
	}

	return cmd
}
