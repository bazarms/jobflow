// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	//"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	log "github.com/uthng/golog"

	"github.com/uthng/jobflow/config"
	// import all jobflow builtin modules
	_ "github.com/uthng/jobflow/plugins/all"
)

var (
	jobexec string
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Exec command is to execute jobs",
	Long:  `Exec command is to execute a specific job. If no job specified, all jobs will get executed in the order.`,
	Run: func(cmd *cobra.Command, args []string) {
		//log.SetVerbosity(log.DEBUG)
		exec(args)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	execCmd.PersistentFlags().StringVar(&jobexec, "job", "all", "Job's name. Default: all")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func exec(args []string) {
	// Check if a flow file is specified
	if len(args) < 1 {
		log.Fatalln("No jobflow file is specified")
	}

	// Read flow file
	content, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatalw("Cannot read jobflow file", "file", args[0], "err", err)
	}

	jf := config.ReadFlowFile(content)

	//Execute all jobs
	if jobexec == "all" {
		log.Debugw("List of jobs", "jobs", jf.Jobs)
		// Loop jobs and exec job by job.
		for _, j := range jf.Jobs {
			log.Infow("Executing job:", "job", j.Name)
			// Loop tasks and exec task by task
			j.Start = j.Tasks[0]

			res := j.Run("")
			if res != nil {
				log.Fatalw("Job error", "job", j.Name, "err", res.Error)
			}

			log.Infow("Job result", "job", j.Name, "res", j.ValueRegistry.ValueList)
		}
	}
}
