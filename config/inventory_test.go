// vim: ts=4 et sts=4 sw=4
package config

import (
	//"bytes"
	//"fmt"
	//"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uthng/jobflow/job"
	//log "github.com/uthng/golog"
)

func TestReadInventoryFile(t *testing.T) {
	var yamlInventoryFile = []byte(`
global:

hosts:
  host1:
    host: host1.com
    port: 25
    user: userhost1
    password: passhost1
    vars:
      hostvar: host1
  host2:
    host: host2.com
    port: 5555
    user: userhost2
    privatekey: privatekeyhost2
    vars:
      hostvar: host2
  host3:
    host: host3.com
    privatekey: privatekeyhost3
    vars:
      hostvar: host3
  host4:
    host: host4.com
    password: passhost4

groups:
  group1:
    hosts:
    - host1
    - host2
    vars:
      group1var1: group1
      group1var2: group1
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
      group3var1: group3
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
				Name:     "host1",
				Host:     "host1.com",
				Port:     25,
				User:     "userhost1",
				Password: "passhost1",
				Vars: map[string]interface{}{
					"hostvar": "host1",
				},
			},
			"host2": job.Host{
				Name:       "host2",
				Host:       "host2.com",
				Port:       5555,
				User:       "userhost2",
				PrivateKey: "privatekeyhost2",
				Vars: map[string]interface{}{
					"hostvar": "host2",
				},
			},
			"host3": job.Host{
				Name:       "host3",
				Host:       "host3.com",
				Port:       22,
				User:       "root",
				PrivateKey: "privatekeyhost3",
				Vars: map[string]interface{}{
					"hostvar": "host3",
				},
			},
			"host4": job.Host{
				Name:     "host4",
				Host:     "host4.com",
				Port:     22,
				User:     "root",
				Password: "passhost4",
				Vars:     map[string]interface{}{},
			},
		},
		Groups: map[string]job.Group{
			"group1": job.Group{
				Name:  "group1",
				Hosts: []string{"host1", "host2"},
				Vars: map[string]interface{}{
					"group1var1": "group1",
					"group1var2": "group1",
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
					"group3var1": "group3",
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

	inventory := ReadInventoryFile(yamlInventoryFile)
	assert.Equal(t, output, inventory)
}
