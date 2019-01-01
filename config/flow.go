package config

import (
	//"fmt"
	//"io/ioutil"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"

	log "github.com/uthng/golog"
	"github.com/uthng/jobflow/job"
)

// JobFlow represents the yaml file describing variables
// and all jobs
type JobFlow struct {
	Variables map[string]interface{}
	Jobs      []*job.Job
}

// ReadFlowFile unmarshals the configuration file into Config struct
func ReadFlowFile(content []byte) *JobFlow {
	config := make(map[string]interface{})

	jf := &JobFlow{}

	// Get config under map[string]interface{}
	//content, err := ioutil.ReadFile("testdata/hello")
	//if err != nil {
	//log.Fatalw("Cannot read flow file", "file", file, "err", err)
	//}

	err := yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatalw("Cannot unmarshal flow file content", "err", err)
	}

	for k, v := range config {
		if k == "variables" {
			jf.Variables = cast.ToStringMap(v)
		} else {
			j := job.NewJob(k)

			// Parse tasks
			readJob(j, cast.ToStringMap(v))

			// Add variables to job context
			j.ValueRegistry.AddValue("context", jf.Variables)

			// Add job to job list
			jf.Jobs = append(jf.Jobs, j)
		}
	}

	return jf
}

////////////// INTERNAL FUNCTIONS ////////////////////////

// readJob parses & fills up Job structure
func readJob(j *job.Job, data map[string]interface{}) {
	//Read tasks
	tasks := cast.ToSlice(data["tasks"])
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

			module := k
			cmd := cast.ToString(vm["cmd"])
			task.Params = cast.ToStringMap(vm["params"])

			c, ok := job.GetCmdByName(module + "." + cmd)
			if !ok {
				log.Fatalw("No command found", "cmd", cmd, "module", module)
			}

			task.Func = c.Func
		}

		j.AddTask(task)
	}
}
