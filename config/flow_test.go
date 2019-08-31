// vim: ts=4 et sts=4 sw=4
package config

import (
	//"bytes"
	//"fmt"
	//"reflect"
	"testing"

	//"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"

	"github.com/uthng/jobflow/job"
	// import all GoJobs's builtin modules
	//_ "github.com/uthng/jobflow/plugins/all"
	//log "github.com/uthng/golog"
)

func TestReadFlowFile(t *testing.T) {
	var yamlFlowFile = []byte(`
variables:
  var1: $VAR1
  var2: ${VAR2}

jobs:
- name: build
  tasks:
  - shell:
     cmd: exec
     params:
       cmd: echo 10
  - shell:
      cmd: exec
      params:
        cmd: echo 20

- hosts: swmmng
  tasks:
  - name: "github release"
    github:
      cmd: release
      params:
        target: hello
`)

	//cmdFuncShellExec, _ := job.GetCmdByName("shell.exec")
	//cmdFuncGithubRelease, _ := job.GetCmdByName("github.release")

	flowOK := &job.Flow{
		Variables: map[string]interface{}{
			"var1": "$VAR1",
			"var2": "${VAR2}",
		},
		Jobs: []*job.Job{
			{
				Name:  "build",
				Hosts: "localhost",
				Tasks: []*job.Task{
					{
						Name: "task-1",
						//Func: cmdFuncShellExec.Func,
						Params: map[string]interface{}{
							"cmd": "echo 10",
						},
						OnSuccess: "task-2",
					},
					{
						Name: "task-2",
						//Func: cmdFuncShellExec.Func,
						Params: map[string]interface{}{
							"cmd": "echo 20",
						},
					},
				},
			},
			{
				Name:  "job-2",
				Hosts: "swmmng",
				Tasks: []*job.Task{
					{
						Name: "github release",
						//Func: cmdFuncGithubRelease.Func,
						Params: map[string]interface{}{
							"target": "hello",
						},
					},
				},
			},
		},
	}

	jf := job.NewFlow()

	ReadFlow(jf, yamlFlowFile)

	assert.Equal(t, flowOK.Variables, jf.Variables)

	expectedJobs := make(map[string]interface{})
	actualJobs := make(map[string]interface{})

	for _, job := range jf.Jobs {
		actualJobs[job.Name] = job
	}

	for _, job := range flowOK.Jobs {
		expectedJobs[job.Name] = job
	}

	assert.Equal(t, len(expectedJobs), len(actualJobs))

	for k, v := range expectedJobs {
		expected := v.(*job.Job)
		actual := actualJobs[k].(*job.Job)

		assert.Equal(t, expected.Name, actual.Name)
		assert.Equal(t, expected.Hosts, actual.Hosts)

		expectedTasks := make(map[string]interface{})
		actualTasks := make(map[string]interface{})

		for _, task := range actual.Tasks {
			actualTasks[task.Name] = task
		}

		for _, task := range expected.Tasks {
			expectedTasks[task.Name] = task
		}

		for k, v := range expectedTasks {
			expected := v.(*job.Task)
			actual := actualTasks[k].(*job.Task)

			assert.Equal(t, expected.Name, actual.Name)
			assert.Equal(t, expected.Params, actual.Params)
			assert.Equal(t, expected.OnSuccess, actual.OnSuccess)
			assert.Equal(t, expected.OnFailure, actual.OnFailure)
			//assert.Equal(t, expected.Result, actual.Result)
		}
	}
}
