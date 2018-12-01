package gojobs_test

import (
	"errors"
	//"fmt"
	"testing"

	"github.com/uthng/gojobs"
)

func TestCheckTasksSuccess(t *testing.T) {

	task1 := &gojobs.Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 2",
	}
	task2 := &gojobs.Task{
		Name: "Task 2",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 3",
	}
	task3 := &gojobs.Task{
		Name: "Task 3",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
	}

	taskErr := &gojobs.Task{
		Name: "Exit Error",
		Func: nil,
	}

	w := gojobs.NewJob("Job 1")
	w.Start = task1
	w.AddTask(task1)
	w.AddTask(task2)
	w.AddTask(task3)
	w.AddTask(taskErr)

	res := w.CheckTasks()
	if res != true {
		t.Fail()
	}
}

func TestCheckTasksOnFailure(t *testing.T) {

	task1 := &gojobs.Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 2",
	}
	task2 := &gojobs.Task{
		Name: "Task 2",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Hmmmmmmm",
		OnSuccess: "Task 3",
	}
	task3 := &gojobs.Task{
		Name: "Task 3",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
	}

	taskErr := &gojobs.Task{
		Name: "Exit Error",
		Func: nil,
	}

	w := gojobs.NewJob("Job 2")
	w.Start = task1
	w.AddTask(task1)
	w.AddTask(task2)
	w.AddTask(task3)
	w.AddTask(taskErr)

	res := w.CheckTasks()
	if res == true {
		t.Fail()
	}
}

func TestCheckTasksOnSuccess(t *testing.T) {

	task1 := &gojobs.Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 2",
	}
	task2 := &gojobs.Task{
		Name: "Task 2",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Exit Error",
		OnSuccess: "Task 4",
	}
	task3 := &gojobs.Task{
		Name: "Task 3",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
	}

	taskErr := &gojobs.Task{
		Name: "Exit Error",
		Func: nil,
	}

	w := gojobs.NewJob("Job 3")
	w.Start = task1
	w.AddTask(task1)
	w.AddTask(task2)
	w.AddTask(task3)
	w.AddTask(taskErr)

	res := w.CheckTasks()
	if res == true {
		t.Fail()
	}
}

func TestOneTaskJob(t *testing.T) {
	var testVar bool

	testVar = false

	task := &gojobs.Task{
		Name: "Task 1",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			testVar = true
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
	}

	w := gojobs.NewJob("Job 4")
	w.Start = task
	w.AddTask(task)

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

	task1 := &gojobs.Task{
		Name: "Task1",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Error",
		OnSuccess: "Task2",
	}
	task3 := &gojobs.Task{
		Name: "Task3",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: errors.New("Task3: NOK"), Result: nil}
		},
		OnFailure: "Rollback1",
		OnSuccess: "Task4",
	}
	task2 := &gojobs.Task{
		Name: "Task2",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnSuccess: "Task3",
		OnFailure: "Error",
	}
	task4 := &gojobs.Task{
		Name: "Task4",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Error",
	}

	rollback1 := &gojobs.Task{
		Name: "Rollback1",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Error",
		OnSuccess: "Rollback2",
	}

	rollback2 := &gojobs.Task{
		Name: "Rollback2",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
		OnFailure: "Error",
	}

	err := &gojobs.Task{
		Name: "Error",
		Func: func(m map[string]interface{}) *gojobs.CmdResult {
			return &gojobs.CmdResult{Error: nil, Result: nil}
		},
	}

	w := gojobs.NewJob("Job 5")
	w.Start = task1
	w.AddTask(task1)
	w.AddTask(task2)
	w.AddTask(task3)
	w.AddTask(task4)
	w.AddTask(rollback1)
	w.AddTask(rollback2)
	w.AddTask(err)

	res := w.Run("")
	if res != nil {
		t.Error()
	}
}
