package config

import (
	//"fmt"
	"io/ioutil"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"

	log "github.com/uthng/golog"
	"github.com/uthng/jobflow/job"
)

// ReadFlowFile reads the flow content from a file and
// create a new instance Flow
func ReadFlowFile(file string) *job.Flow {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalw("Cannot read jobflow file", "file", file, "err", err)
	}

	jf := job.NewFlow()
	jf.InventoryFile = file

	jf.IsOnRemote = false

	ReadFlow(jf, content)

	return jf
}

// ReadFlow unmarshals the content of the configuration file
// into Config struct
func ReadFlow(jf *job.Flow, content []byte) {
	config := make(map[string]interface{})

	// Get config under map[string]interface{}
	//content, err := ioutil.ReadFile("testdata/hello")
	//if err != nil {
	//log.Fatalw("Cannot read flow file", "file", file, "err", err)
	//}

	err := yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatalw("Cannot unmarshal flow file content", "err", err)
	}

	// Tip to keep order while parsing config map
	var keys []string
	for k := range config {
		keys = append(keys, k)
	}

	for _, k := range keys {
		v := config[k]

		if k == "on_remote" {
			jf.IsOnRemote = cast.ToBool(v)
		} else if k == "variables" {
			jf.Variables = cast.ToStringMap(v)
		} else {
			j := job.NewJob(k)

			// Parse tasks
			readJob(j, cast.ToStringMap(v))

			// Add job to job list
			jf.Jobs = append(jf.Jobs, j)
		}
	}
}

////////////// INTERNAL FUNCTIONS ////////////////////////

// readJob parses & fills up Job structure
func readJob(j *job.Job, data map[string]interface{}) {
	hosts := cast.ToString(data["hosts"])
	if hosts == "" {
		j.Hosts = "localhost"
	} else {
		j.Hosts = hosts
	}

	//Read tasks
	tasks := cast.ToSlice(data["tasks"])
	if len(tasks) <= 0 {
		log.Warnw("No tasks specified", "job", j.Name)
	}

	for i, t := range tasks {
		task := &job.Task{}

		tm := cast.ToStringMap(t)
		// Check name
		n, ok := tm["name"]
		if ok {
			task.Name = cast.ToString(n)
			delete(tm, "name")
		} else {
			task.Name = "task-" + cast.ToString(i+1)
		}

		for k, v := range tm {
			vm := cast.ToStringMap(v)

			plugin := k

			// Check plugin's mandatory parameters
			// Check cmd parameter
			cmd := cast.ToString(vm["cmd"])
			if cmd == "" {
				log.Fatalw("No command is specified", "plugin", plugin)
			}

			// Check params parameter
			task.Params = cast.ToStringMap(vm["params"])
			if len(task.Params) == 0 {
				log.Fatalw("No parameter is specified", "plugin", plugin)
			}

			// If OnSuccess of the previous task is not specified
			// so set it to the current task. Like that, all tasks
			// can be executed in case of onsuccess not specified
			if i > 0 && j.Tasks[i-1].OnSuccess == "" {
				j.Tasks[i-1].OnSuccess = task.Name
			}

			c, ok := job.GetCmdByName(plugin + "." + cmd)
			if !ok {
				log.Fatalw("No command found", "cmd", cmd, "plugin", plugin)
			}

			task.Cmd = c
		}

		j.AddTask(task)
	}
}
