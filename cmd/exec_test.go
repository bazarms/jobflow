// +build unit

package cmd

import (
	//"fmt"
	"os"
	"testing"
	//    "reflect"
	"runtime"

	"github.com/stretchr/testify/assert"

	"github.com/bazarms/jobflow/job"
	//log "github.com/uthng/golog"
)

func TestExec(t *testing.T) {
	testCases := []struct {
		name     string
		yamlFile string
		output   map[string]map[string]*job.CmdResult
	}{
		{
			"ShellExec",
			"./data/exec.yml",
			map[string]map[string]*job.CmdResult{
				"job1": map[string]*job.CmdResult{
					"shell11": &job.CmdResult{
						//Error: nil,
						Result: map[string]interface{}{
							"result": "var1\n",
						},
					},
					"shell12": &job.CmdResult{
						Error: nil,
						Result: map[string]interface{}{
							"result": "var2\n",
						},
					},
					"shell13": &job.CmdResult{
						Error: nil,
						Result: map[string]interface{}{
							"result": "var1+var2\n",
						},
					},
				},
				"job2": map[string]*job.CmdResult{
					"shell21": &job.CmdResult{
						Error: nil,
						Result: map[string]interface{}{
							"result": "var1/var2\n",
						},
					},
					"shell22": &job.CmdResult{
						Error: nil,
						Result: map[string]interface{}{
							"result": "var1*var2\n",
						},
					},
				},
			},
		},
	}

	pluginDir = "../bin/" + runtime.GOOS + "_" + runtime.GOARCH + "/plugins"
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var args []string

			//log.SetVerbosity(log.DEBUG)

			os.Setenv("VAR1", "var1")
			os.Setenv("VAR2", "var2")
			args = append(args, tc.yamlFile)

			jf := exec(args)
			//assert.Equal(t, tc.output, jf.Result)
			for _, j := range jf.Result["localhost"] {
				assert.Equal(t, tc.output[j.Name], j.Result)
			}
		})
	}

}
