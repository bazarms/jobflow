package job

import (
	log "github.com/uthng/golog"
)

/////// DECLARATION OF ALL TYPES /////////////////////////

// Flow represents job flow YAML file containing
// different sections such as variables, multiple jobs etc.
type Flow struct {
	Variables map[string]interface{}
	Jobs      []*Job
}

////////// DEFINITION OF ALL FUNCTIONS ///////////////////////////

// NewFlow instancies a new Flow
func NewFlow() *Flow {
	flow := &Flow{
		Variables: make(map[string]interface{}),
	}

	return flow
}

// RunAllJobs executes all jobs
func (f *Flow) RunAllJobs() {
	// Loop jobs and exec job by job.
	for _, j := range f.Jobs {
		f.execJob(j)
	}
}

// RunJob executes a specified job with the name given
func (f *Flow) RunJob(job string) {
	if job == "" {
		log.Errorln("No job name is specified")
		return
	}

	// Loop jobs and exec job by job.
	for _, j := range f.Jobs {
		if j.Name == job {
			f.execJob(j)
		}
	}
}

/////////// INTERNAL FUNCTIONS /////////////////////////:

func (f *Flow) execJob(job *Job) error {
	job.Start = job.Tasks[0]

	// Set context to execute job
	job.Context["variables"] = f.Variables

	return job.Run("")
}
