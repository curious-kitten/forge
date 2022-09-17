package tool

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
