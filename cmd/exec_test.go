package cmd

import (
	//"fmt"
	"os"
	"testing"
	//    "reflect"

	"github.com/stretchr/testify/assert"
	//log "github.com/uthng/golog"
)

func TestExec(t *testing.T) {
	testCases := []struct {
		name     string
		yamlFile string
		output   []map[string]interface{}
	}{
		{
			"ShellExec",
			"./data/exec.yml",
			[]map[string]interface{}{
				map[string]interface{}{
					"shell11": map[string]interface{}{
						"result": "var1\n",
					},
					"shell12": map[string]interface{}{
						"result": "var2\n",
					},
					"shell13": map[string]interface{}{
						"result": "var1+var2\n",
					},
				},
				map[string]interface{}{
					"shell21": map[string]interface{}{
						"result": "var1/var2\n",
					},
					"shell22": map[string]interface{}{
						"result": "var1*var2\n",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var args []string

			//log.SetVerbosity(log.DEBUG)

			os.Setenv("VAR1", "var1")
			os.Setenv("VAR2", "var2")
			args = append(args, tc.yamlFile)

			jf := exec(args)
			for i, j := range jf.Jobs {
				assert.Equal(t, tc.output[i], j.Result)
			}
		})
	}

}
