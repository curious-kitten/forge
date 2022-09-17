package tool

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

var (
	ErrTmplBuildFailuer = errors.New("error while building the template")
	ErrTmplExecFailuer  = errors.New("error while executing the template")
)

type Executor func(command string) (string, error)

type Tool struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Args        []Argument `yaml:"args"`
	Cmd         string     `yaml:"cmd"`
}

func (t Tool) Run(exec Executor, o ...Option) (string, error) {
	opts := options{
		args: map[string]string{},
	}
	for _, arg := range t.Args {
		opts.args[arg.Name] = arg.Default
	}
	for _, opt := range o {
		opts = opt(opts)
	}
	tmpl, err := template.New("cmd").Parse(t.Cmd)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrTmplBuildFailuer, err.Error())
	}
	buff := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buff, "cmd", opts.args)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrTmplExecFailuer, err.Error())
	}
	return exec(buff.String())
}

type Argument struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
}
