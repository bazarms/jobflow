package cmd

import (
	"os"
	"testing"
	//    "reflect"
	//"github.com/stretchr/testify/assert"
	//log "github.com/uthng/golog"
)

func TestExec(t *testing.T) {
	testCases := []struct {
		name     string
		yamlFile string
	}{
		{
			"ShellExec",
			"./data/exec.yml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var args []string

			//log.SetVerbosity(log.DEBUG)

			os.Setenv("VAR1", "var1")
			os.Setenv("VAR2", "var2")
			args = append(args, tc.yamlFile)

			exec(args)
		})
	}

}
