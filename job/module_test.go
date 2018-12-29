package job_test

import (
	//"fmt"
	"testing"

	"github.com/uthng/jobflow/job"
	_ "github.com/uthng/jobflow/plugins/all"
	log "github.com/uthng/golog"
)

func TestLoadModules(t *testing.T) {
	log.SetVerbosity(log.DEBUG)
	//job.NewCmdRegistry()
	//job.NewModuleRegistry()

	//job.LoadModules("../app/modules")

	//registry := job.GetModuleRegistry()
	//fmt.Printf("Registry %#v\n", registry)

	cmd, ok := job.GetCmdByName("shell.exec")
	if ok {
		result := cmd.Func(map[string]interface{}{"cmd": "ls -la"})
		log.Infoln(result)
	} else {
		t.Fail()
	}
}
