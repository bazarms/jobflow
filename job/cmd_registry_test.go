package job_test

import (
	//"fmt"
	"testing"

	log "github.com/uthng/golog"
	"github.com/bazarms/jobflow/job"
)

var module = job.Module{
	Name:        "ModTest",
	Version:     "0.1",
	Description: "ModTest",
}

var fn = func(map[string]interface{}) *job.CmdResult {
	log.Debugln("CmdFunc test")
	return &job.CmdResult{Error: nil, Result: nil}
}

var cmds = []job.Cmd{
	{
		Name:   "cmd1",
		Func:   fn,
		Module: module,
	},
	{
		Name:   "cmd2",
		Func:   fn,
		Module: module,
	},
	{
		Name:   "cmd3",
		Func:   fn,
		Module: module,
	},
}

func TestCmdRegister(t *testing.T) {

	for _, cmd := range cmds {
		log.Debugln(cmd)
		job.CmdRegister(cmd)
	}

	registry := job.GetCmdRegistry()
	log.Debugf("Registry %#v\n", registry)

	nb := job.GetNbOfCmds()
	log.Debugf("nb of commands %v\n", nb)

	if nb != 3 {
		t.Fail()
	}
}

func TestCmdUnregister(t *testing.T) {

	job.CmdUnregister(cmds[1])

	registry := job.GetCmdRegistry()
	log.Debugf("Registry %#v\n", registry)

	nb := job.GetNbOfCmds()
	log.Debugf("nb of commands %v\n", nb)

	if nb != 2 {
		t.Fail()
	}

	_, ok := registry.CmdList["cmd1"]
	if ok {
		t.Fail()
	}

	_, ok = registry.CmdList["cmd2"]
	if ok {
		t.Fail()
	}

}

func TestGetCmdByName(t *testing.T) {
	cmd, ok := job.GetCmdByName("ModTest.cmd3")
	if ok {
		cmd.Func(nil)
	} else {
		t.Fail()
	}
}
