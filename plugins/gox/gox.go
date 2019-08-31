package main

import "C"
import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cast"
	log "github.com/uthng/golog"
	"github.com/uthng/jobflow/job"
)

var plugin = job.Plugin{
	Name:        "gox",
	Version:     "0.1",
	Description: "Use gox to build multiple platforms",
}

// List of available commands for this plugin
var commands = []job.Cmd{
	{
		Name:   "build",
		Func:   CmdBuild,
		Plugin: plugin,
	},
}

// Init initializes plugin by registering all its commands
// to command registry
func init() {
	for _, cmd := range commands {
		job.CmdRegister(cmd)
	}
}

// CmdBuild compiles multiple platforms.
// It takes a map of params
func CmdBuild(params map[string]interface{}) *job.CmdResult {
	var res = job.NewCmdResult()
	var args []string

	args = append(args, "gox")

	value, ok := params["osarch"]
	if ok == false {
		res.Error = fmt.Errorf("param osarch missing")
		return res
	}
	v := cast.ToStringSlice(value)
	oa := "-osarch=\"" + strings.Join(v, " ") + "\""
	args = append(args, oa)

	value, ok = params["output"]
	if ok == false {
		res.Error = fmt.Errorf("param output missing")
		return res
	}
	o := "-output=\"" + cast.ToString(value) + "\""
	args = append(args, o)

	log.Debugw("Executing command", "args", args)

	// Execute kubectl command
	cmd := exec.Command(args[0], args[1:len(args)]...)
	//cmd := exec.CommandContext(ctx, "gox", args)

	// Check if error
	output, err := cmd.CombinedOutput()
	if err != nil {
		res.Error = err
		return res
	}

	res.Result["result"] = string(output)
	return res
}
