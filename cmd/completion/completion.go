/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package completion

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Command(command *cobra.Command) *cobra.Command {
	// cmd represents the completion command
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

	. <(bitbucket completion)

	To configure your bash shell to load completions for each session add to your bashrc

	# ~/.bashrc or ~/.profile
	. <(bitbucket completion)
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := command.GenBashCompletion(os.Stdout)
			if err != nil {
				return fmt.Errorf("completion error: %w", err)
			}
			return nil
		},
	}

	return cmd
}
