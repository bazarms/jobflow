// vim: ts=4 et sts=4 sw=4

package config

import (
	//"fmt"
	//"io/ioutil"

	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"

	log "github.com/uthng/golog"
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
			}
		}
	}

	return inventory
}

////////////// INTERNAL FUNCTIONS ////////////////////////

// readHost parses & fills up Inventory Host structure
func readHost(h *job.Host, data map[string]interface{}) {
	h.Host = cast.ToString(data["host"])
	if h.Host == "" {
		h.Host = "localhost"
	}

	h.Port = cast.ToInt(data["port"])
	if h.Port == 0 {
		h.Port = 22
	}

	h.User = cast.ToString(data["user"])
	if h.User == "" {
		h.User = "root"
	}

	h.Password = cast.ToString(data["password"])
	h.PrivateKey = cast.ToString(data["privatekey"])

	//Read vars
	h.Vars = cast.ToStringMap(data["vars"])
	h.Groups = cast.ToStringSlice(data["groups"])
}

// readGroup parses & fills up Inventory Group structure
func readGroup(g *job.Group, data map[string]interface{}) {
	g.Hosts = cast.ToStringSlice(data["hosts"])
	g.Vars = cast.ToStringMap(data["vars"])
}
