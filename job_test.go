package gojobs

import (
	"errors"
	//"fmt"
	"os"
	//"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	//"github.com/uthng/gojobs"
	//log "github.com/uthng/golog"
)

func TestCheckTasksSuccess(t *testing.T) {

	task1 := Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 2",
	}
	task2 := Task{
		Name: "Task 2",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 3",
	}
	task3 := Task{
		Name: "Task 3",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
	}

	taskErr := Task{
		Name: "Exit Error",
		Func: nil,
	}

	w := NewJob("Job 1")
	w.Start = &task1
	w.AddTask(&task1)
	w.AddTask(&task2)
	w.AddTask(&task3)
	w.AddTask(&taskErr)

	res := w.CheckTasks()
	if res != true {
		t.Fail()
	}
}

func TestCheckTasksOnFailure(t *testing.T) {

	task1 := Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 2",
	}
	task2 := Task{
		Name: "Task 2",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Hmmmmmmm",
		OnSuccess: "Task 3",
	}
	task3 := Task{
		Name: "Task 3",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
	}

	taskErr := Task{
		Name: "Exit Error",
		Func: nil,
	}

	w := NewJob("Job 2")
	w.Start = &task1
	w.AddTask(&task1)
	w.AddTask(&task2)
	w.AddTask(&task3)
	w.AddTask(&taskErr)

	res := w.CheckTasks()
	if res == true {
		t.Fail()
	}
}

func TestCheckTasksOnSuccess(t *testing.T) {

	task1 := Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 2",
	}
	task2 := Task{
		Name: "Task 2",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 4",
	}
	task3 := Task{
		Name: "Task 3",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
	}

	taskErr := Task{
		Name: "Exit Error",
		Func: nil,
	}

	w := NewJob("Job 3")
	w.Start = &task1
	w.AddTask(&task1)
	w.AddTask(&task2)
	w.AddTask(&task3)
	w.AddTask(&taskErr)

	res := w.CheckTasks()
	if res == true {
		t.Fail()
	}
}

func TestOneTaskJob(t *testing.T) {
	var testVar bool

	testVar = false

	task := Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *CmdResult {
			testVar = true
			return &CmdResult{Error: nil, Result: nil}
		},
	}

	w := NewJob("Job 4")
	w.Start = &task
	w.AddTask(&task)

	err := w.Run("")
	if err != nil {
		t.Error(err)
	}

	if testVar != true {
		t.Fail()
	}
}

func TestMultipleTasks(t *testing.T) {
	//var testVar bool

	//testVar = false

	task1 := Task{
		Name: "Task1",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Error",
		OnSuccess: "Task2",
	}
	task3 := Task{
		Name: "Task3",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: errors.New("Task3: NOK"), Result: nil}
		},
		OnFailure: "Rollback1",
		OnSuccess: "Task4",
	}
	task2 := Task{
		Name: "Task2",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnSuccess: "Task3",
		OnFailure: "Error",
	}
	task4 := Task{
		Name: "Task4",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Error",
	}

	rollback1 := Task{
		Name: "Rollback1",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Error",
		OnSuccess: "Rollback2",
	}

	rollback2 := Task{
		Name: "Rollback2",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Error",
	}

	err := Task{
		Name: "Error",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: nil}
		},
	}

	w := NewJob("Job 5")
	w.Start = &task1
	w.AddTask(&task1)
	w.AddTask(&task2)
	w.AddTask(&task3)
	w.AddTask(&task4)
	w.AddTask(&rollback1)
	w.AddTask(&rollback2)
	w.AddTask(&err)

	res := w.Run("")
	if res != nil {
		t.Error()
	}
}

func TestExpandEnvContext(t *testing.T) {
	input := []byte(`
var1: $VAR1
var2:
  - 1
  - $VAR21
  - 3
var3:
  - var31
  - $VAR32
  - var33
var4:
  var41:
    var411: var411
    var412: ${VAR412}
  var42: "var42"
  var43: $VAR43
`)

	output := map[string]interface{}{
		"var1": "var1",
		"var2": []string{"1", "2", "3"},
		"var3": []string{"var31", "var32", "var33"},
		"var4": map[string]interface{}{
			"var41": map[string]interface{}{
				"var411": "var411",
				"var412": "var412",
			},
			"var42": "var42",
			"var43": "var43",
		},
	}

	data := make(map[string]interface{})

	err := yaml.Unmarshal(input, data)
	assert.Nil(t, err)

	os.Setenv("VAR1", "var1")
	os.Setenv("VAR21", "2")
	os.Setenv("VAR32", "var32")
	os.Setenv("VAR412", "var412")
	os.Setenv("VAR43", "var43")

	context := expandEnvContext(data)

	assert.Equal(t, context, output)
}

func TestRenderTaskTemplate(t *testing.T) {
	ctx := []byte(`
var1: $VAR1
var2:
  - 1
  - $VAR21
  - 3
var3:
  - var31
  - $VAR32
  - var33
var4:
  var41:
    var411: var411
    var412: ${VAR412}
  var42: "var42"
  var43: $VAR43
`)

	output := map[string]interface{}{
		"param1": "var1 var1",
		"param2": []string{"2", "2", "3"},
		"param3": "var412 var412",
	}

	data := make(map[string]interface{})

	err := yaml.Unmarshal(ctx, data)
	assert.Nil(t, err)

	os.Setenv("VAR1", "var1")
	os.Setenv("VAR21", "2")
	os.Setenv("VAR32", "var32")
	os.Setenv("VAR412", "var412")
	os.Setenv("VAR43", "var43")

	task1 := Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *CmdResult {
			return &CmdResult{Error: nil, Result: m}
		},
		Params: map[string]interface{}{
			"param1": "$VAR1 {{ .context.var1 }}",
			"param2": []string{"$VAR21", "{{ index .context.var2 1 }}", "3"},
			"param3": "$VAR412 {{ .context.var4.var41.var412 }}",
		},
	}

	w := NewJob("Job 1")
	w.ValueRegistry.AddValue("context", data)
	w.Start = &task1

	res := w.Run("")
	assert.Nil(t, res)

	// Check result of task1
	result, ok := w.ValueRegistry.GetValueByKey("Task 1")
	assert.Equal(t, true, ok)
	assert.Equal(t, result, output)
}
