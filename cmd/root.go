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
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/cruious-kitten/forge/cmd/completion"
	"github.com/cruious-kitten/forge/cmd/execute"
	"github.com/cruious-kitten/forge/cmd/run"
	"github.com/cruious-kitten/forge/cmd/version"
	"github.com/cruious-kitten/forge/pkg/forge"
)

var forgeFile string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Root() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "forge",
		Short: "Forge helps developers and system administrators to create command chains that can be easily executed",
		Long: `Forge uses a yaml file structure where you can define complex actions that can be executed from a simple CLI.
	Forgeries defined in forge carry the name "forgery". 
	You can define a forgery by simply providing if with a name and commands to execute, a list of predefined forgeries or a
	combination of the two.
	`,
	}
	setFilePath(rootCmd)
	rootCmd.AddCommand(version.Command())
	rootCmd.AddCommand(completion.Command(rootCmd))
	rootCmd.AddCommand(execute.Command(getForgery(forgeFile)))
	rootCmd.AddCommand(run.Command(getForgery(forgeFile)))

	return rootCmd
}

func setFilePath(cmd *cobra.Command) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory")
	}
	cmd.PersistentFlags().StringVarP(&forgeFile, "file", "f", filepath.Join(wd, "Forge.yaml"), "Forge.yaml file (default is Forge in current directory)")
}

func getForgery(filePath string) func() (forge.Forge, error) {
	return func() (forge.Forge, error) {
		f, err := os.Open(filePath)
		if err != nil {
			return forge.Forge{}, err
		}
		defer f.Close()
		forgery, err := forge.FromFile(f)
		if err != nil {
			return forge.Forge{}, fmt.Errorf("%w: file %s", err, filePath)
		}
		return forgery, nil
	}
}
