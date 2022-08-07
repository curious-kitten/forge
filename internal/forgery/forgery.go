package forgery

import (
	"bytes"
	"errors"
	"io"
	"text/template"

	"gopkg.in/yaml.v3"
)

var (
	ErrToolNotFound = errors.New("could not find tool")
)

type forge struct {
	Tools []Tool `yaml:"tools"`
}

func (f forge) GetTool(name string) (Commander, error) {
	for _, v := range f.Tools {
		if v.Name == name {
			return v, nil
		}
	}
	return nil, ErrToolNotFound
}

type Tool struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Args        []Argument `yaml:"args"`
	Cmd         string     `yaml:"cmd"`
}

func (t Tool) Command(o ...Option) (string, error) {
	opts := options{
		args: map[string]string{},
	}
	for _, arg := range t.Args {
		if arg.Default != "" {
			opts.args[arg.Name] = arg.Default
		}
	}
	for _, opt := range o {
		opts = opt(opts)
	}
	tmpl, err := template.New("cmd").Parse(t.Cmd)
	if err != nil {
		return "", err
	}
	buff := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buff, "cmd", opts.args)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

type Argument struct {
	Name    string `yaml:"name"`
	Default string `yaml:"default"`
}

func Read(file io.Reader) (Forge, error) {
	frg := forge{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&frg); err != nil {
		return frg, err
	}
	return frg, nil
}

type Forge interface {
	GetTool(name string) (Commander, error)
}

type options struct {
	args map[string]string
}

type Option func(options) options

func WithArgument(key, value string) Option {
	return func(o options) options {
		o.args[key] = value
		return o
	}
}

type Commander interface {
	Command(opts ...Option) (string, error)
}
