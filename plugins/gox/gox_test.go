// +build unit

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bazarms/jobflow/job"
	//"github.com/bazarms/jobflow/plugins/gox"
)

func TestCmdBuild(t *testing.T) {
	testCases := []struct {
		name   string
		params map[string]interface{}
		result *job.CmdResult
	}{
		{
			"OSArchMissing",
			map[string]interface{}{
				"output": "output1",
			},
			&job.CmdResult{
				Error:  fmt.Errorf("param osarch missing"),
				Result: map[string]interface{}{},
			},
		},
		{
			"OutputMissing",
			map[string]interface{}{
				"osarch": "osarch1",
			},
			&job.CmdResult{
				Error:  fmt.Errorf("param output missing"),
				Result: map[string]interface{}{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CmdBuild(tc.params)
			assert.Equal(t, result, tc.result)
		})
	}

}
