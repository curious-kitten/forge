package forgery_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mimatache/forge/internal/forgery"
)

func Test_ForgeryTools(t *testing.T) {
	file, err := os.Open("./testdata/forgery_works.yaml")
	assert.NoError(t, err, "could not read test data file")
	frg, err := forgery.Read(file)
	assert.NoError(t, err)

	tts := []struct {
		name    string
		cmdName string
		cmd     string
		arg     map[string]string
	}{
		{
			name:    "single line cmd",
			cmdName: "simpleCmd",
			cmd:     "echo simpleCmd",
		},
		{
			name:    "multi line cmd",
			cmdName: "multilineCmd",
			cmd:     "echo line 1 \n echo line 2",
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
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			tool, err := frg.GetTool(tt.cmdName)
			assert.NoError(t, err)
			opts := []forgery.Option{}
			for k, v := range tt.arg {
				opts = append(opts, forgery.WithArgument(k, v))
			}
			cmd, err := tool.Command(opts...)
			assert.NoError(t, err)
			assert.Equal(t, cmd, cmd)
		})
	}
}

func Test_Forgery_ToolNotFound(t *testing.T) {
	file, err := os.Open("./testdata/forgery_works.yaml")
	assert.NoError(t, err, "could not read test data file")
	frg, err := forgery.Read(file)
	assert.NoError(t, err)
	_, err = frg.GetTool("badName")
	assert.True(t, errors.Is(err, forgery.ErrToolNotFound), "invalid error")
}
