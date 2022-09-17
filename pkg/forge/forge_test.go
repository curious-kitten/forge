package forge_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cruious-kitten/forge/pkg/forge"
)

type TestExecutor struct {
	expectedOutput string
	t              *testing.T
}

func (t TestExecutor) Execute(command string) (string, error) {
	assert.Equal(t.t, t.expectedOutput, command)
	return "", nil
}

func Test_RunTool(t *testing.T) {
	file, err := os.Open("./testdata/forgery_works.yaml")
	assert.NoError(t, err, "could not read test data file")
	frg, err := forge.FromFile(file)
	assert.NoError(t, err)

	tts := []struct {
		name        string
		cmdName     string
		cmd         string
		arg         map[string]string
		expectedErr error
	}{
		{
			name:    "single line cmd",
			cmdName: "simpleCmd",
			cmd:     "echo simpleCmd",
		},
		{
			name:    "multi line cmd",
			cmdName: "multilineCmd",
			cmd:     "echo line 1 \necho line 2 \n",
		},
		{
			name:    "template cmd",
			cmdName: "templateCmd",
			cmd:     "echo hello world",
		},
		{
			name:    "template cmd set value",
			cmdName: "templateCmd",
			cmd:     "echo cool value",
			arg: map[string]string{
				"name": "cool value",
			},
		},
		{
			name:        "missing name",
			cmdName:     "badName",
			expectedErr: forge.ErrNotFound,
		},
		{
			name:        "missing name",
			cmdName:     "badName",
			expectedErr: forge.ErrNotFound,
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			te := TestExecutor{
				t:              t,
				expectedOutput: tt.cmd,
			}
			_, err := frg.RunTool(te, tt.cmdName, tt.arg)
			assert.True(t, errors.Is(err, tt.expectedErr))
		})
	}
}
