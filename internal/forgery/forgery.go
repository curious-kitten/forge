package forgery

import (
	"errors"
	"fmt"
	"strings"
)

const (
	outputPrefix = "output."
)

var (
	ErrTmplBuildFailuer = errors.New("error while building the template")
	ErrTmplExecFailuer  = errors.New("error while executing the template")
)

type RunTool func(name string, args map[string]string) (string, error)

type Forgery struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Tools       []Tool `yaml:"tools"`
}

type Tool struct {
	Args   map[string]string `yaml:"args"`
	Name   string            `yaml:"name"`
	Output string            `yaml:"output"`
}

type Argument struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (f Forgery) RunForgery(run RunTool) error {
	outputs := map[string]string{}
	for _, tls := range f.Tools {
		for k, v := range tls.Args {
			if !strings.HasPrefix(v, outputPrefix) {
				continue
			}

			val, ok := outputs[strings.TrimPrefix(v, outputPrefix)]
			if !ok {
				continue
			}
			tls.Args[k] = val
		}
		output, err := run(tls.Name, tls.Args)
		outputs[tls.Output] = output
		if err != nil {
			return fmt.Errorf("forgery exited with errors")
		}
	}
	return nil
}
