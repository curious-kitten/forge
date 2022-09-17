package tool_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cruious-kitten/forge/internal/tool"
)

func testExecutor(t *testing.T, expectedOutput string) func(command string) (string, error) {
	return func(command string) (string, error) {
		assert.Equal(t, expectedOutput, command)
		return "", nil
	}
}

func Test_ForgeryTools(t *testing.T) {

	tts := []struct {
		tool           tool.Tool
		arg            map[string]string
		expectedOutput string
	}{
		{
			tool: tool.Tool{
				Name:        "simpleCmd",
				Description: "A single line command",
				Cmd:         "echo simpleCmd",
			},
			expectedOutput: "echo simpleCmd",
		},
		{
			tool: tool.Tool{
				Name:        "multilineCmd",
				Description: "A multi line command",
				Cmd: `echo line 1
echo line 2`,
			},
			expectedOutput: "echo line 1\necho line 2",
		},
		{
			tool: tool.Tool{
				Name:        "templateCmd",
				Description: "Using a template for the CMD",
				Args: []tool.Argument{
					{
						Name:    "name",
						Default: "hello world",
					},
				},
				Cmd: "echo {{.name}}",
			},
			expectedOutput: "echo hello world",
		},
		{
			tool: tool.Tool{
				Name: "templateCmdSetValue",
				Args: []tool.Argument{
					{
						Name:    "name",
						Default: "hello world",
					},
				},
				Cmd: "echo {{.name}}",
			},
			arg: map[string]string{
				"name": "cool value",
			},
			expectedOutput: "echo cool value",
		},
	}
	for _, tt := range tts {
		t.Run(tt.tool.Name, func(t *testing.T) {
			te := testExecutor(t, tt.expectedOutput)
			opts := []tool.Option{}
			for k, v := range tt.arg {
				opts = append(opts, tool.WithArgument(k, v))
			}
			_, err := tt.tool.Run(te, opts...)
			assert.NoError(t, err)
		})
	}
}
