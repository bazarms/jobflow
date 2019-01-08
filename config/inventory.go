// vim: ts=4 et sts=4 sw=4

package config

import (
	//"fmt"
	//"io/ioutil"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"

	log "github.com/uthng/golog"
	utils "github.com/uthng/goutils"

	"github.com/uthng/jobflow/job"
)

// ReadInventoryFile unmarshals the inventory file
func ReadInventoryFile(content []byte) *job.Inventory {
	config := make(map[string]interface{})

	inventory := job.NewInventory()

	err := yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatalw("Cannot unmarshal inventory file content", "err", err)
	}

	for k, v := range config {
		if k == "global" {
			inventory.Global = cast.ToStringMap(v)
		}
		if k == "hosts" {
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

		if k == "groups" {
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

	return inventory
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
