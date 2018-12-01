package gojobs

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"

	log "github.com/uthng/golog"
	utils "github.com/uthng/goutils"
)

/////////// DECLARATION OF ALL TYPES //////////////////////////

// Job describes structure of a job
type Job struct {
	Name  string
	Start *Task

	Tasks         []*Task
	ValueRegistry *ValueRegistry
}

// Task describes attributes of a task
type Task struct {
	Name      string
	Func      CmdFunc
	Params    map[string]interface{}
	OnSuccess string
	OnFailure string
	Result    *CmdResult
}

////////// DEFINITION OF ALL FUNCTIONS ///////////////////////////

// NewJob instancies a new Job
func NewJob(name string) *Job {
	job := &Job{
		Name: name,
	}

	job.ValueRegistry = NewValueRegistry()

	return job
}

// Run a job throught all tasks
//
// Firstly, it checks to ensure that all task's names
// are valid task.fmt_unicode
//
func (job *Job) Run(tasks string) error {
	var res bool
	var err error

	log.Infoln("JOB RUN STARTED")

	// Check the name of all tasks indicated in taskflow
	res = job.CheckTasks()
	if res == false {
		return fmt.Errorf("Error while checking task names")
	}

	// Run certain tasks given in parameter
	if tasks != "" {
		err = job.RunTaskByTask(tasks)
	} else {
		// Run complete taskflow by running the first task
		err = job.RunAllTasks(job.Start)
	}

	if err != nil {
		log.Errorln("JOB RUN FAILED")
		return err
	}

	log.Infoln("JOB RUN COMPLETED")

	log.Debugw("Value Registry", "registry", job.ValueRegistry)

	return nil
}

// AddTask adds a new task to the job
func (job *Job) AddTask(task *Task) {
	if task == nil {
		return
	}

	job.Tasks = append(job.Tasks, task)

	return
}

// RunTaskByTask executes only task function of tasks specified
// in command line parameter. It return error if a task fails
func (job *Job) RunTaskByTask(tasks string) error {
	for _, task := range strings.Split(tasks, ",") {
		log.Infow("Running", "task", task)

		s, err := job.GetTaskByName(task)
		if s == nil {
			log.Errorln(err)
			return err
		}

		if s.Func == nil {
			log.Warnln("Ignored", "task", task, "reason", "func is nil")
			continue
		}

		// Before execute command func, we must render each param template
		// if it exists with  Value registry
		err = job.RenderTaskTemplate(s, job.ValueRegistry.ValueList)
		if err != nil {
			log.Errorw("Template rendering error", "task", s.Name, "err", err)
			return err
		}

		s.Result = s.Func(s.Params)

		if s.Result.Error != nil {
			log.Errorln("Execution error", "task", s.Name, "err", s.Result.Error)
			log.Errorw("Result", "task", s.Name, "status", "KO")
			return s.Result.Error
		}

		// In all cases, add task result to value registry
		job.ValueRegistry.AddValue(s.Name, s.Result.Result)
		log.Infow("Result", "task", s.Name, "status", "OK")
	}

	return nil
}

// RunAllTasks executes task function. If it returns error,
// check if task on failure is specified and then go on it.
// Otherwise, check and go on task on Success if specified.
func (job *Job) RunAllTasks(task *Task) error {
	log.Infow("Running", "task", task.Name)

	if task.Func == nil {
		log.Warnw("Ignored", "task", task.Name, "reason", "func is nil")
		return nil
	}

	// Before execute command func, we must render each param template
	// if it exists with  Value registry
	err := job.RenderTaskTemplate(task, job.ValueRegistry.ValueList)
	if err != nil {
		log.Errorw("Error templating", "task", task.Name, "err", err)
		return err
	}

	task.Result = task.Func(task.Params)

	if task.Result.Error != nil {
		log.Errorw("Execution error", "task", task.Name, "err", task.Result.Error)
		log.Errorw("Result", "task", task.Name, "status", "KO")
		// Go the task of failure if specified
		if len(task.OnFailure) > 0 {
			taskFailure, _ := job.GetTaskByName(task.OnFailure)
			job.RunAllTasks(taskFailure)
		}
	} else {
		// In all cases, add task result to value registry
		job.ValueRegistry.AddValue(task.Name, task.Result.Result)

		log.Infow("Result", "task", task.Name, "status", "OK")
		// Go the task of Success if specified
		if len(task.OnSuccess) > 0 {
			taskSuccess, _ := job.GetTaskByName(task.OnSuccess)
			job.RunAllTasks(taskSuccess)
		}
	}

	return task.Result.Error
}

// GetTaskByName returns task by its name in the task list of the job
func (job *Job) GetTaskByName(name string) (*Task, error) {
	for _, task := range job.Tasks {
		if task.Name == name {
			return task, nil
		}
	}

	return nil, fmt.Errorf("Task doesnot exist: %v", name)
}

// CheckTasks checks all tasks to see if the name given for task on
// failure or on success matches valid task names
func (job *Job) CheckTasks() bool {
	var taskNames []interface{}
	var res bool

	// Comparaison function of 2 strings
	fn := func(str1 interface{}, str2 interface{}) bool {
		if str1 == str2 {
			return true
		}
		return false
	}

	// Loop tasks to get a list of task names
	for _, task := range job.Tasks {
		taskNames = append(taskNames, task.Name)
	}

	// Loop again all tasks and check for each task, the name specified
	// in Task On Success or Task On Failure exists in the list of task names
	for _, task := range job.Tasks {
		if task.OnSuccess != "" {
			res = utils.IsElementInArray(task.OnSuccess, taskNames, fn)
			if res == false {
				log.Errorw("Task does not exist !", "task", task.OnSuccess)
				return false
			}
		}
		if task.OnFailure != "" {
			res = utils.IsElementInArray(task.OnFailure, taskNames, fn)
			if res == false {
				log.Errorw("Task does not exist !", "task", task.OnFailure)
				return false
			}
		}
	}

	return true
}

// RenderTaskTemplate renders go template in each param with
// the values in ValueRegistry
func (job *Job) RenderTaskTemplate(task *Task, data map[string]interface{}) error {
	var err error
	var tpl bytes.Buffer

	for key, value := range task.Params {
		tpl.Reset()
		// Render only string value
		switch v := value.(type) {
		case string:
			// Create a new template with name : task name + key
			log.Debugw("Template rendering", "task", task.Name, "value", value.(string), "type", v)
			t := template.New(task.Name + "-" + key).Funcs(sprig.TxtFuncMap())
			t, err = t.Parse(value.(string))
			if err != nil {
				log.Errorw("Template parsing error", "task", task.Name, "key", key)
				return err
			}

			err = t.Execute(&tpl, data)
			if err != nil {
				log.Errorw("Template rendering error", "task", task.Name, "key", key)
				return err
			}

			// Assign new rendered value to param key
			log.Debugw("New value rendered", "task", task.Name, "key", key, "tpl", tpl.String())
			task.Params[key] = tpl.String()

		default:
			log.Warnw("Param type ignored", "type", v)
		}
	}

	return nil
}
