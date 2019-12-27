// vim: ts=4 et sts=4 sw=4
package config

import (
	//"bytes"
	//"fmt"
	//"reflect"
	"testing"

	//"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/bazarms/jobflow/job"
	//log "github.com/uthng/golog"
)

func TestReadInventoryFile(t *testing.T) {
	var yamlInventoryFile = []byte(`
global:

hosts:
  host1:
    jobflow_ssh_host: host1.com
    jobflow_ssh_port: 25
    jobflow_ssh_user: userhost1
    jobflow_ssh_pass: passhost1
    hostvar: host1
  host2:
    jobflow_ssh_host: host2.com
    jobflow_ssh_port: 5555
    jobflow_ssh_user: userhost2
    jobflow_ssh_privkey: privatekeyhost2
    hostvar: host2
  host3:
    jobflow_ssh_host: host3.com
    jobflow_ssh_privkey: privatekeyhost3
    hostvar: host3
  host4:
    jobflow_ssh_pass: passhost4

groups:
  group1:
    hosts:
    - host1
    - host2
    vars:
      jobflow_ssh_user: group1user
      group1var1: group1
  group2:
    hosts:
    - host2
    - host3
    vars:
      group2var1: group2
      group2var2: group2
  group3:
    hosts:
    - host3
    - host4
    vars:
      hostvar: group3var1
      jobflow_ssh_privkey: privatekeygroup3
#  group4:
#    groups:
#    - group1
#    - group2
#    hosts:
#    - host4
#    vars:
#      group4var1: group4
`)

	output := &job.Inventory{
		Global: map[string]interface{}{},
		Hosts: map[string]job.Host{
			"host1": job.Host{
				Name:   "host1",
				Groups: []string{"group1"},
				Vars: map[string]interface{}{
					"jobflow_ssh_host":    "host1.com",
					"jobflow_ssh_port":    25,
					"jobflow_ssh_user":    "group1user",
					"jobflow_ssh_pass":    "passhost1",
					"jobflow_ssh_privkey": "",
					"hostvar":             "host1",
					"group1var1":          "group1",
				},
			},
			"host2": job.Host{
				Name:   "host2",
				Groups: []string{"group1", "group2"},
				Vars: map[string]interface{}{
					"jobflow_ssh_host":    "host2.com",
					"jobflow_ssh_port":    5555,
					"jobflow_ssh_user":    "group1user",
					"jobflow_ssh_pass":    "",
					"jobflow_ssh_privkey": "privatekeyhost2",
					"hostvar":             "host2",
					"group1var1":          "group1",
					"group2var1":          "group2",
					"group2var2":          "group2",
				},
			},
			"host3": job.Host{
				Name:   "host3",
				Groups: []string{"group2", "group3"},
				Vars: map[string]interface{}{
					"jobflow_ssh_host":    "host3.com",
					"jobflow_ssh_port":    22,
					"jobflow_ssh_user":    "root",
					"jobflow_ssh_pass":    "",
					"jobflow_ssh_privkey": "privatekeygroup3",
					"hostvar":             "group3var1",
					"group2var1":          "group2",
					"group2var2":          "group2",
				},
			},
			"host4": job.Host{
				Name:   "host4",
				Groups: []string{"group3"},
				Vars: map[string]interface{}{
					"jobflow_ssh_host":    "localhost",
					"jobflow_ssh_port":    22,
					"jobflow_ssh_user":    "root",
					"jobflow_ssh_pass":    "passhost4",
					"jobflow_ssh_privkey": "privatekeygroup3",
					"hostvar":             "group3var1",
				},
			},
		},
		Groups: map[string]job.Group{
			"group1": job.Group{
				Name:  "group1",
				Hosts: []string{"host1", "host2"},
				Vars: map[string]interface{}{
					"group1var1":       "group1",
					"jobflow_ssh_user": "group1user",
				},
			},
			"group2": job.Group{
				Name:  "group2",
				Hosts: []string{"host2", "host3"},
				Vars: map[string]interface{}{
					"group2var1": "group2",
					"group2var2": "group2",
				},
			},
			"group3": job.Group{
				Name:  "group3",
				Hosts: []string{"host3", "host4"},
				Vars: map[string]interface{}{
					"hostvar":             "group3var1",
					"jobflow_ssh_privkey": "privatekeygroup3",
				},
			},
			//"group4": job.Group{
			//Hosts: []string{"host1", "host2", "host3", "host4"},
			//Vars: map[string]interface{}{
			//"group4var1": "group4",
			//},
			//},
		},
	}

	inventory := job.NewInventory()
	ReadInventory(inventory, yamlInventoryFile)
	assert.Equal(t, output.Global, inventory.Global)
	for k, v := range inventory.Hosts {
		expected := output.Hosts[k]
		actual := v

		assert.Equal(t, expected.Name, actual.Name)
		assert.ElementsMatch(t, expected.Groups, actual.Groups)
		assert.Equal(t, expected.Vars, actual.Vars)

	}
	for k, v := range inventory.Groups {
		expected := output.Groups[k]
		actual := v

		assert.Equal(t, expected.Name, actual.Name)
		assert.ElementsMatch(t, expected.Hosts, actual.Hosts)
		assert.Equal(t, expected.Vars, actual.Vars)
	}
}
