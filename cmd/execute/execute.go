package execute

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cruious-kitten/forge/internal/shell"
	"github.com/cruious-kitten/forge/pkg/forge"
)

type forgeryReader func() (forge.Forge, error)

var values []string

func argProvider(fr forgeryReader) []string {
	f, err := fr()
	if err != nil {
		return []string{}
	}
	return f.GetTools()
}

const longDescription = `
A toll is a single execution "script" that will be run.
A tool can have as little as one shell command or multiple commands that will be executed.
Single command example:

`

func Command(fr forgeryReader) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "exec [tool]",
		Short:     "Execute a tool",
		Long:      longDescription,
		Args:      cobra.ExactValidArgs(1),
		ValidArgs: argProvider(fr),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := fr()
			if err != nil {
				return err
			}
			exec := shell.VerboseExecutor(os.Stdout, os.Stderr, shell.NewExecutor(context.Background()))
			cmdArgs := make(map[string]string, len(values))
			for _, v := range values {
				vls := strings.SplitN(v, "=", 2)
				if len(vls) != 2 {
					return fmt.Errorf("invalid set argument: %s", v)
				}
				cmdArgs[vls[0]] = vls[1]
			}
			_, err = f.RunTool(exec, args[0], cmdArgs)
			if err != nil {
				cmd.SilenceUsage = true
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&values, "set", []string{}, "set argument values for running a specific tool (eg: --set key=value)")

	return cmd
}
