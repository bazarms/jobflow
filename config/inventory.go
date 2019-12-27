// vim: ts=4 et sts=4 sw=4

package config

import (
	//"fmt"
	"io/ioutil"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"

	log "github.com/uthng/golog"
	utils "github.com/uthng/goutils"

	"github.com/bazarms/jobflow/job"
)

// ReadInventoryFile reads the flow content from a file and
// create a new instance Inventory
func ReadInventoryFile(file string) *job.Inventory {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalw("Cannot read inventory file", "file", file, "err", err)
	}

	inventory := job.NewInventory()

	ReadInventory(inventory, content)

	return inventory
}

// ReadInventory unmarshals the inventory file
func ReadInventory(inventory *job.Inventory, content []byte) {
	config := make(map[string]interface{})

	err := yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatalw("Cannot unmarshal inventory file content", "err", err)
	}

	v, ok := config["global"]
	if ok {
		inventory.Global = cast.ToStringMap(v)
	}

	v, ok = config["hosts"]
	if ok {
		// Parse hosts
		for n, h := range cast.ToStringMap(v) {
			host := job.Host{
				Name: n,
				Vars: make(map[string]interface{}),
			}

			readHost(&host, cast.ToStringMap(h))
			inventory.Hosts[n] = host
		}
	}

	v, ok = config["groups"]
	if ok {
		// Parse groups
		for n, g := range cast.ToStringMap(v) {
			group := job.Group{
				Name: n,
			}

			readGroup(&group, cast.ToStringMap(g))
			inventory.Groups[n] = group

			// Merge vars & groups in host
			for _, h := range group.Hosts {
				host, ok := inventory.Hosts[h]
				if !ok {
					log.Fatalw("Host in the group not found", "host", h, "group", group.Name)
				}

				vars, err := utils.MapStringMerge(host.Vars, group.Vars)
				if err != nil {
					log.Fatalw("Cannot merge group's vars in host's vars", "err", err)
				}

				host.Vars = cast.ToStringMap(vars)
				host.Groups = append(host.Groups, group.Name)

				// Update host
				inventory.Hosts[h] = host
			}
		}
	}
}

////////////// INTERNAL FUNCTIONS ////////////////////////

// readHost parses & fills up Inventory Host structure
func readHost(h *job.Host, data map[string]interface{}) {
	// Init default variables for host
	h.Vars["jobflow_ssh_host"] = "localhost"
	h.Vars["jobflow_ssh_port"] = 22
	h.Vars["jobflow_ssh_user"] = "root"
	h.Vars["jobflow_ssh_pass"] = ""
	h.Vars["jobflow_ssh_privkey"] = ""

	h.Groups = cast.ToStringSlice(data["groups"])

	for k, v := range data {
		h.Vars[k] = v
	}
}

// readGroup parses & fills up Inventory Group structure
func readGroup(g *job.Group, data map[string]interface{}) {
	g.Hosts = cast.ToStringSlice(data["hosts"])
	g.Vars = cast.ToStringMap(data["vars"])
}
