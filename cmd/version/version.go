package version

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cruious-kitten/forge/internal/info"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show forge version",
		Run: func(cmd *cobra.Command, args []string) {
			inf := info.AppInfo()
			fmt.Printf(
				`
version:     %s
build hash:  %s
build date:  %s
`, inf.Version, inf.Hash, inf.BuildDate)
		},
	}

	return cmd
}
