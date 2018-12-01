package gox

import "C"
import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cast"
	"github.com/uthng/gojobs"
	log "github.com/uthng/golog"
)

var module = gojobs.Module{
	Name:        "gox",
	Version:     "0.1",
	Description: "Use gox to build multiple platforms",
}

// List of available commands for this module
var commands = []gojobs.Cmd{
	{
		Name:   "build",
		Func:   CmdBuild,
		Module: module,
	},
}

// Init initializes module by registering all its commands
// to command registry
func init() {
	for _, cmd := range commands {
		gojobs.CmdRegister(cmd)
	}
}

// CmdBuild compiles multiple platforms.
// It takes a map of params
func CmdBuild(params map[string]interface{}) *gojobs.CmdResult {
	var res = gojobs.NewCmdResult()
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

	log.Debugln("gox: execute command:", args)

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
