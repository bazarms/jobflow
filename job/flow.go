package job

import (
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

	RemoteExecDir string
	InventoryFile string
}

////////// DEFINITION OF ALL FUNCTIONS ///////////////////////////

// NewFlow instancies a new Flow
func NewFlow() *Flow {
	flow := &Flow{
		Variables:     make(map[string]interface{}),
		RemoteExecDir: "$HOME",
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
			err := f.execJob(j)
			if err != nil {
				log.Errorln(err)
			}
		}
	}
}

/////////// INTERNAL FUNCTIONS /////////////////////////:

func (f *Flow) execJob(job *Job) error {
	if job.Hosts == "" || job.Hosts == "localhost" || job.Hosts == "127.0.0.1" {
		return f.execJobLocal(job)
	}

	return f.execJobRemote(job)
}

// execJobLocal executes job on the current host directly
func (f *Flow) execJobLocal(job *Job) error {
	job.Start = job.Tasks[0]

	// Set context to execute job
	job.Context["variables"] = f.Variables

	return job.Run("")
}

// execJobRemote executes job on remote hosts
func (f *Flow) execJobRemote(job *Job) error {
	//var sshClients []*gossh.Client
	var config *gossh.Config
	var err error

	// Check if job hosts is a group or only a host
	// If it is a group, loop all hosts to init a ssh client

	group, ok := f.Inventory.Groups[job.Hosts]
	if ok {
		for _, hostname := range group.Hosts {

			host, ok := f.Inventory.Hosts[hostname]
			if !ok {
				err := fmt.Errorf("host %s in the group %s not found", hostname, group.Name)
				log.Errorw("Error inventory", "err", err)
				return err
			}

			sshUser := cast.ToString(host.Vars["jobflow_ssh_user"])
			sshPass := cast.ToString(host.Vars["jobflow_ssh_pass"])
			sshHost := cast.ToString(host.Vars["jobflow_ssh_host"])
			sshPort := cast.ToInt(host.Vars["jobflow_ssh_port"])
			sshPrivkey := cast.ToString(host.Vars["jobflow_ssh_privkey"])

			if sshPrivkey != "" {
				config, err = gossh.NewClientConfigWithKeyFile(sshUser, sshPrivkey, sshHost, sshPort, false)
				if err != nil {
					log.Errorw("Error SSH connection", "user", sshUser, "host", sshHost, "port", sshPort, "privkey", sshPrivkey, "err", err)
					return err
				}
			} else if sshPass != "" {
				config, err = gossh.NewClientConfigWithUserPass(sshUser, sshPass, sshHost, sshPort, false)
				if err != nil {
					log.Errorw("Error SSH connection", "user", sshUser, "host", sshHost, "port", sshPort, "pass", "********", "err", err)
					return err
				}
			} else {
				err := fmt.Errorf("no ssh password or private key is specified for connection")
				log.Errorw("Error SSH connection", "err", err)
				return err
			}

			client, err := gossh.NewClient(config)
			if err != nil {
				log.Errorw("Error creating SSH client", "user", sshUser, "host", sshHost, "port", sshPort, "err", err)
				return err
			}

			bytes, err := client.ExecCommand("echo " + f.RemoteExecDir)
			log.Infoln(string(bytes))

			// Find location of jobflow binary on the local machine
			//var dirAbsPath string
			exec, err := os.Executable()
			if err != nil {
				//dirAbsPath = filepath.Dir(ex)
				//fmt.Println(ex)
				log.Errorw("Error getting current binary path", "err", err)
				return err
			}

			// Random string
			randStr := randomString(10)
			remoteDir := f.RemoteExecDir + "/" + randStr
			binExec := filepath.Base(exec)

			// Create a tmp on remote machine
			_, err = client.ExecCommand("mkdir -p " + remoteDir)
			if err != nil {
				log.Errorw("Failed to create a remote folder", "dir", remoteDir, "err", err)
				return err
			}

			// SCP jobflow binary from local machine to remote machine
			err = client.SCPFile(exec, remoteDir+"/"+binExec, "0755")
			if err != nil {
				log.Errorw("Failed to scp file to remote machine", "exec", exec, "err", err)
				return err
			}

			// SCP other files: jobflow yaml containing only
			// the current job to remote machine
			newFlow, err := f.generateLocalFlowRemoteMachine(job)
			if err != nil {
				log.Errorw("Failed to generate new local flow file for remote machine", "err", err)
				return err
			}

			err = client.SCPBytes(newFlow, remoteDir+"/flow.yml", "0755")
			if err != nil {
				log.Errorw("Failed to scp new flow file to remote machine", "err", err)
				return err
			}

			//time.Sleep(time.Second * 5)

			// Execute jobflow on remote machine with new location
			_, err = client.ExecCommand(remoteDir + "/" + binExec + " " + remoteDir + "/flow.yml")
			if err != nil {
				log.Errorw("Failed to execute flow file on remote machine", "err", err)
				return err
			}

			//Remove tmp folder on remote machine
			_, err = client.ExecCommand("rm -rf " + remoteDir)
			if err != nil {
				log.Errorw("Failed to remove folder on remote machine", "dir", remoteDir, "err", err)
				return err
			}

		}
	}

	return nil
}

func (f *Flow) generateLocalFlowRemoteMachine(j *Job) ([]byte, error) {
	mFlow := make(map[string]interface{})
	job := make(map[string]interface{})
	tasks := []map[string]interface{}{}

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
	mFlow[j.Name] = job

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
