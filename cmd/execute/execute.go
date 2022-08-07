package execute

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mimatache/forge/internal/forgery"
	"github.com/mimatache/forge/internal/shell"
	"github.com/spf13/cobra"
)

type forgeryReader func() (forgery.Forge, error)

var values []string

func Command(reader forgeryReader) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "execute a tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := reader()
			if err != nil {
				return err
			}
			if len(args) != 1 {
				return fmt.Errorf("%s command expects an argument to be vicen representing the tool name", cmd.Name())
			}
			v, err := f.GetTool(args[0])
			if err != nil {
				return err
			}

			options := []forgery.Option{}
			for _, v := range values {
				vls := strings.SplitN(v, "=", 2)
				if len(vls) != 2 {
					return fmt.Errorf("invalid set argument: %s", v)
				}
				options = append(options, forgery.WithArgument(vls[0], vls[1]))
			}

			c, err := v.Command(options...)
			if err != nil {
				return err
			}
			exec := shell.Executor{
				Ctx: context.Background(),
				Out: os.Stdout,
				Err: os.Stderr,
			}
			err = exec.Execute(c)
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
