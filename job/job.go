package job

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/spf13/cast"

	log "github.com/uthng/golog"
	utils "github.com/uthng/goutils"
)

/////////// DECLARATION OF ALL TYPES //////////////////////////

// Job describes structure of a job
type Job struct {
	Name  string
	Hosts string
	Start *Task

	Tasks   []*Task
	Context map[string]interface{}
	Result  map[string]interface{}
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
		Name:    name,
		Context: make(map[string]interface{}),
		Result:  make(map[string]interface{}),
	}

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

	log.Infow("JOB RUN STARTED", "job", job.Name, "hosts", job.Hosts)
	log.Debugw("Job context", "context", job.Context)

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
		log.Errorw("JOB RUN FAILED", "job", job.Name)
		return err
	}

	log.Debugw("Job result", "result", job.Result)
	log.Infow("JOB RUN COMPLETED", "job", job.Name, "hosts", job.Hosts)

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
		err = job.RenderTaskTemplate(s)
		if err != nil {
			log.Errorw("Template rendering error", "task", s.Name, "err", err)
			return err
		}

		s.Result = s.Func(s.Params)

		if s.Result.Error != nil {
			log.Errorw("Task result", "task", s.Name, "err", s.Result.Error)
			job.Result[s.Name] = s.Result.Error
			return s.Result.Error
		}

		// In all cases, add task result to value registry
		job.Result[s.Name] = s.Result.Result
		log.Infow("Result", "task", s.Name, "result", s.Result.Result)
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
	err := job.RenderTaskTemplate(task)
	if err != nil {
		log.Errorw("Error templating", "task", task.Name, "err", err)
		return err
	}

	task.Result = task.Func(task.Params)

	if task.Result.Error != nil {
		log.Errorw("Result", "task", task.Name, "err", task.Result.Error)
		job.Result[task.Name] = task.Result.Error

		// Go the task of failure if specified
		if len(task.OnFailure) > 0 {
			taskFailure, _ := job.GetTaskByName(task.OnFailure)
			job.RunAllTasks(taskFailure)
		}
	} else {
		// In all cases, add task result to value registry
		log.Infow("Result", "task", task.Name, "result", task.Result.Result)
		job.Result[task.Name] = task.Result.Result

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
// the values in Job Context & Result
func (job *Job) RenderTaskTemplate(task *Task) error {
	var tpl bytes.Buffer

	// Combine Job Context & Result into one map to render template
	data := make(map[string]interface{})
	data["context"] = job.Context
	data["result"] = job.Result

	// Expand env vars for context
	d := expandEnvContext(data)

	for key, value := range task.Params {
		tpl.Reset()

		kind := reflect.ValueOf(value).Kind()
		// Render only string value
		// Check if kind is struct or ptr, do nothing
		if kind == reflect.Array || kind == reflect.Slice {
			arr := []string{}

			for idx, it := range cast.ToStringSlice(value) {
				str, err := renderParamTemplate(task.Name, key+"["+cast.ToString(idx)+"]", it, d)
				if err != nil {
					return err
				}

				arr = append(arr, str)
			}

			task.Params[key] = arr
		} else if kind == reflect.Map || kind == reflect.Struct || kind == reflect.Ptr {
			log.Warnw("Param kind ignored", "kind", kind)
		} else {
			str, err := renderParamTemplate(task.Name, key, value, d)
			if err != nil {
				return err
			}

			task.Params[key] = str
		}
	}

	return nil
}

//////////////// INTERNAL FUNCTIONS ////////////////////

func renderParamTemplate(task, key string, value interface{}, data map[string]interface{}) (string, error) {
	var tpl bytes.Buffer

	// Create a new template with name : task name + key
	log.Debugw("Template rendering", "task", task, "value", value.(string), "type", reflect.TypeOf(value).Name())
	t := template.New(task + "-" + key).Funcs(sprig.TxtFuncMap())

	// Expand env vars before template rendering
	v := cast.ToString(value)
	v = os.ExpandEnv(v)
	t, err := t.Parse(cast.ToString(v))
	if err != nil {
		log.Errorw("Template parsing error", "task", task, "key", key)
		return "", err
	}

	err = t.Execute(&tpl, data)
	if err != nil {
		log.Errorw("Template rendering error", "task", task, "key", key)
		return "", err
	}

	// Assign new rendered value to param key
	log.Debugw("New value rendered", "task", task, "key", key, "tpl", tpl.String())

	return tpl.String(), nil
}

//expandEnvContext expands values of env variables
func expandEnvContext(data map[string]interface{}) map[string]interface{} {
	d := make(map[string]interface{})

	for k, v := range data {
		kind := reflect.ValueOf(v).Kind()

		// Check if kind is struct or ptr, do nothing
		if kind == reflect.Array || kind == reflect.Slice {
			arr := []string{}
			for _, it := range cast.ToStringSlice(v) {
				arr = append(arr, os.ExpandEnv(it))
			}
			d[k] = arr
		} else if kind == reflect.Map {
			m := expandEnvContext(cast.ToStringMap(v))
			d[k] = m
		} else if kind == reflect.Struct || kind == reflect.Ptr {
			// Do nothing
			d[k] = v
		} else {
			d[k] = os.ExpandEnv(cast.ToString(v))
		}
	}

	return d
}
