package forge

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/cruious-kitten/forge/internal/forgery"
	"github.com/cruious-kitten/forge/internal/tool"
)

var (
	ErrNotFound     = errors.New("could not find item")
	ErrBuildFailure = errors.New("failed to build tool")
)

func IsInternalError(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, ErrBuildFailure)
}

type Executor interface {
	Execute(command string) (string, error)
}

type Forge struct {
	Tools     []*tool.Tool       `yaml:"tools"`
	Forgeries []*forgery.Forgery `yaml:"forgeries"`
}

func (f Forge) RunTool(exec Executor, name string, args map[string]string) (string, error) {
	return f.runTool(exec)(name, args)
}

func (f Forge) runTool(exec Executor) func(name string, args map[string]string) (string, error) {
	return func(name string, args map[string]string) (string, error) {
		execTool, ok := Find(f.Tools, func(v *tool.Tool) bool { return v.Name == name })
		if !ok {
			return "", fmt.Errorf("%w: tool %s ", ErrNotFound, name)
		}
		options := []tool.Option{}
		for k, v := range args {
			options = append(options, tool.WithArgument(k, v))
		}
		value, err := execTool.Run(exec.Execute, options...)

		if err != nil {
			return "", fmt.Errorf("tool exited with errors")
		}
		return strings.TrimSpace(value), nil
	}
}

func (f Forge) GetTools() []string {
	tools := make([]string, len(f.Tools))
	for i, v := range f.Tools {
		tools[i] = v.Name
	}
	return tools
}

func (f Forge) GetForgeries() []string {
	forgeries := make([]string, len(f.Forgeries))
	for i, v := range f.Forgeries {
		forgeries[i] = v.Name
	}
	return forgeries
}

func (f Forge) RunForgery(exec Executor, name string) error {
	execForge, ok := Find(f.Forgeries, func(v *forgery.Forgery) bool { return v.Name == name })
	if !ok {
		return fmt.Errorf("%w: forgery %s ", ErrNotFound, name)
	}
	for _, t := range execForge.Tools {
		_, ok := Find(f.Tools, func(v *tool.Tool) bool { return v.Name == t.Name })
		if !ok {
			return fmt.Errorf("%w: tool %s ", ErrNotFound, t.Name)
		}
	}
	return execForge.RunForgery(f.runTool(exec))
}

func FromFile(file io.Reader) (Forge, error) {
	frg := Forge{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&frg); err != nil {
		return frg, err
	}
	return frg, nil
}

func Find[T any](items []T, query func(v T) bool) (T, bool) {
	for _, v := range items {
		if query(v) {
			return v, true
		}
	}
	var v T
	return v, false
}
