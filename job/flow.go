package job

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"math/rand"
	"os"
	"path/filepath"
	//"time"

	"github.com/spf13/cast"

	"github.com/uthng/gossh"

	log "github.com/uthng/golog"
)

/////// DECLARATION OF ALL TYPES /////////////////////////

// Flow represents job flow YAML file containing
// different sections such as variables, multiple jobs etc.
type Flow struct {
	Variables map[string]interface{}
	Jobs      []*Job
	Inventory *Inventory

	PluginDir     string
	RemoteExecDir string
	InventoryFile string

	// IsOnRemote indicates if the flow file is on remote machine
	// even if it is local
	IsOnRemote bool

	Status int
	Result map[string][]*Job
}

////////// DEFINITION OF ALL FUNCTIONS ///////////////////////////

// NewFlow instancies a new Flow
func NewFlow() *Flow {
	flow := &Flow{
		Variables:     make(map[string]interface{}),
		RemoteExecDir: "$HOME",
		Status:        SUCCESS,
		Result:        make(map[string][]*Job),
	}

	return flow
}

// RunAllJobs executes all jobs
func (f *Flow) RunAllJobs() {
	// Loop jobs and exec job by job.
	for _, j := range f.Jobs {
		log.Infoln("Executing job", j.Name)
		f.execJob(j)
	}
}

// RunJob executes a specified job with the name given
func (f *Flow) RunJob(job string) error {
	if job == "" {
		err := fmt.Errorf("No job name is specified")
		f.Status = FAILED
		log.Errorw(err.Error())
		return err
	}

	// Loop jobs and exec job by job.
	for _, j := range f.Jobs {
		if j.Name == job {
			err := f.execJob(j)
			if err != nil {
				f.Status = FAILED
				log.Errorw(err.Error())
				return err
			}
		}
	}

	return nil
}

/////////// INTERNAL FUNCTIONS /////////////////////////:

func (f *Flow) execJob(job *Job) error {
	if job.Hosts == "" || job.Hosts == "localhost" || job.Hosts == "127.0.0.1" {
		return f.execJobLocal(job)
	}

	return f.execJobRemote(job)
}

// execJobLocal executes job on the current host directly
func (f *Flow) execJobLocal(j *Job) error {
	job := copyJob(j)

	job.Start = job.Tasks[0]

	// Set context to execute job
	job.Context["variables"] = f.Variables

	jobErr := job.Run("")

	// Marshalling job result to print if it is on remote
	// Store job result only when it is local
	if f.IsOnRemote {
		jobBytes, jsErr := json.Marshal(job.Result)
		if jsErr != nil {
			fmt.Println(jsErr)
		} else {
			fmt.Println(string(jobBytes))
		}
	} else {
		job.Hosts = "localhost"
		f.Result["localhost"] = append(f.Result["localhost"], job)
	}

	return jobErr
}

// execJobRemote executes job on remote hosts
func (f *Flow) execJobRemote(j *Job) error {
	var count int

	channel := make(chan *Job)

	// Check if job hosts is a group or only a host
	// If it is a group, loop all hosts to init a ssh client
	group, ok := f.Inventory.Groups[j.Hosts]
	if ok {
		for _, hostname := range group.Hosts {
			job := copyJob(j)
			job.Hosts = hostname

			go f.execJobViaSSH(job, channel)
		}

		count = len(group.Hosts)
	} else {
		job := copyJob(j)

		go f.execJobViaSSH(job, channel)

		count = 1
	}

	for i := 0; i < count; i++ {
		j := <-channel

		// Store job result
		f.Result[j.Hosts] = append(f.Result[j.Hosts], j)
	}

	for k, v := range f.Result {
		fmt.Println(k, ":")
		for _, j := range v {
			fmt.Printf("\t%s:\n", j.Name)
			for k, v := range j.Result {
				fmt.Printf("\t\t%s: %+v\n", k, v)
			}
		}
	}
	return nil
}

func (f *Flow) execJobViaSSH(j *Job, ch chan *Job) {
	var config *gossh.Config
	var err error

	logger := log.NewLogger()

	logger.Infow("REMOTE JOB RUN STARTED", "job", j.Name, "hosts", j.Hosts)

	host, ok := f.Inventory.Hosts[j.Hosts]
	if !ok {
		logger.Errorw("Host not found", "host", j.Hosts)
		j.Status = FAILED
		return
	}

	logger.Infow("Etablishing ssh connection", "job", j.Name, "hosts", j.Hosts)

	sshUser := cast.ToString(host.Vars["jobflow_ssh_user"])
	sshPass := cast.ToString(host.Vars["jobflow_ssh_pass"])
	sshHost := cast.ToString(host.Vars["jobflow_ssh_host"])
	sshPort := cast.ToInt(host.Vars["jobflow_ssh_port"])
	sshPrivkey := cast.ToString(host.Vars["jobflow_ssh_privkey"])

	if sshPrivkey != "" {
		config, err = gossh.NewClientConfigWithKeyFile(sshUser, sshPrivkey, sshHost, sshPort, false)
		if err != nil {
			logger.Errorw("Error SSH connection", "user", sshUser, "host", sshHost, "port", sshPort, "privkey", sshPrivkey, "err", err)
			j.Status = FAILED
			ch <- j
			return
		}
	} else if sshPass != "" {
		config, err = gossh.NewClientConfigWithUserPass(sshUser, sshPass, sshHost, sshPort, false)
		if err != nil {
			logger.Errorw("Error SSH connection", "user", sshUser, "host", sshHost, "port", sshPort, "pass", "********", "err", err)
			j.Status = FAILED
			ch <- j
			return
		}
	} else {
		logger.Errorw("No ssh password or private key is specified for connection")
		j.Status = FAILED
		ch <- j
		return
	}

	client, err := gossh.NewClient(config)
	if err != nil {
		logger.Errorw("Error creating SSH client", "user", sshUser, "host", sshHost, "port", sshPort, "err", err)
		j.Status = FAILED
		ch <- j
		return
	}

	logger.Infow("Transfering jobflow binary", "job", j.Name, "hosts", j.Hosts)
	// Find location of jobflow binary on the local machine
	//var dirAbsPath string
	exec, err := os.Executable()
	if err != nil {
		//dirAbsPath = filepath.Dir(ex)
		//fmt.Println(ex)
		logger.Errorw("Error getting current binary path", "err", err)
		j.Status = FAILED
		ch <- j
		return
	}

	// Random string
	randStr := randomString(10)
	remoteDir := f.RemoteExecDir + "/." + randStr
	binExec := filepath.Base(exec)

	// Create a tmp on remote machine
	_, err = client.ExecCommand("mkdir -p " + remoteDir)
	if err != nil {
		logger.Errorw("Failed to create a remote folder", "dir", remoteDir, "err", err)
		j.Status = FAILED
		ch <- j
		return
	}

	// Defer function to clean up remote machine
	// and send final job to channel
	defer func() {
		logger.Infow("Clean up remote machine", "job", j.Name, "hosts", j.Hosts, "dir", remoteDir)
		//Remove tmp folder on remote machine
		_, err = client.ExecCommand("rm -rf " + remoteDir)
		if err != nil {
			logger.Errorw("Failed to remove folder on remote machine", "dir", remoteDir, "err", err)
			j.Status = FAILED
			ch <- j
			return
		}

		ch <- j
	}()

	// SCP jobflow binary from local machine to remote machine
	err = client.SCPFile(exec, remoteDir+"/"+binExec, "0755")
	if err != nil {
		logger.Errorw("Failed to scp file to remote machine", "exec", exec, "err", err)
		j.Status = FAILED
		ch <- j
		return
	}

	logger.Infow("Generating local flow file", "job", j.Name, "hosts", j.Hosts)
	// SCP other files: jobflow yaml containing only
	// the current job to remote machine
	newFlow, err := f.generateLocalFlowRemoteMachine(j)
	if err != nil {
		logger.Errorw("Failed to generate new local flow file for remote machine", "err", err)
		j.Status = FAILED
		ch <- j
		return
	}

	logger.Infow("Transfering local flow file", "job", j.Name, "hosts", j.Hosts)
	err = client.SCPBytes(newFlow, remoteDir+"/flow.yml", "0755")
	if err != nil {
		logger.Errorw("Failed to scp new flow file to remote machine", "err", err)
		j.Status = FAILED
		ch <- j
		return
	}

	logger.Infow("Transfering local plugin folder", "job", j.Name, "hosts", j.Hosts)
	err = client.SCPBytes([]byte(f.PluginDir), remoteDir+"/plugins", "0755")
	if err != nil {
		logger.Errorw("Failed to scp plugin folder to remote machine", "err", err)
		j.Status = FAILED
		ch <- j
		return
	}

	//time.Sleep(time.Second * 5)

	logger.Infow("Executing remote jobflow", "job", j.Name, "hosts", j.Hosts)
	// Execute jobflow on remote machine with new location
	remoteCmd := remoteDir + "/" + binExec + " exec --verbosity 0 --plugin-dir " + f.PluginDir + " " + remoteDir + "/flow.yml"
	remoteRes, err := client.ExecCommand(remoteCmd)
	if err != nil {
		logger.Errorw("Failed to execute flow file on remote machine", "job", j.Name, "hosts", j.Hosts, "err", err)
		j.Status = FAILED
		ch <- j
		return
	}

	// Unmarshalling remote result to store in current job locally
	err = json.Unmarshal(remoteRes, &j.Result)
	if err != nil {
		logger.Errorw("Failed to unmarshal remote job result", "res", string(remoteRes))
		j.Status = FAILED
		ch <- j
		return
	}
}

func (f *Flow) generateLocalFlowRemoteMachine(j *Job) ([]byte, error) {
	mFlow := make(map[string]interface{})
	job := make(map[string]interface{})
	tasks := []map[string]interface{}{}
	jobs := []interface{}{}

	mFlow["on_remote"] = "true"
	mFlow["variables"] = f.Variables

	// Tasks
	for _, t := range j.Tasks {
		task := make(map[string]interface{})
		task["name"] = t.Name

		// Extract plugin & cmd name from cmd
		plugin := make(map[string]interface{})
		plugin["cmd"] = t.Cmd.Name
		plugin["params"] = t.Params
		if t.OnSuccess != "" {
			plugin["on_success"] = t.OnSuccess
		}
		if t.OnFailure != "" {
			plugin["on_failure"] = t.OnFailure
		}

		task[t.Cmd.Plugin.Name] = plugin

		tasks = append(tasks, task)
	}

	job["tasks"] = tasks
	jobs = append(jobs, job)
	mFlow["jobs"] = jobs

	return yaml.Marshal(mFlow)
}

func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func copyJob(j *Job) *Job {
	job := NewJob(j.Name)

	job.Hosts = j.Hosts
	job.Start = j.Start
	job.Tasks = j.Tasks

	job.Context = j.Context

	return job
}
