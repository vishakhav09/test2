package ignition

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"

	"github.com/coreos/ignition/config/types"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"ignition_config": resourceConfig(),
			"ignition_disk":   resourceDisk(),
			"ignition_raid":   resourceRaid(),
			"ignition_user":   resourceUser(),
			"ignition_group":  resourceGroup(),
		},
		ConfigureFunc: func(*schema.ResourceData) (interface{}, error) {
			return &cache{
				disks:  make(map[string]*types.Disk, 0),
				arrays: make(map[string]*types.Raid, 0),
				users:  make(map[string]*types.User, 0),
				groups: make(map[string]*types.Group, 0),
			}, nil
		},
	}
}

type cache struct {
	disks  map[string]*types.Disk
	arrays map[string]*types.Raid
	users  map[string]*types.User
	groups map[string]*types.Group

	sync.Mutex
}

func (c *cache) addDisk(g *types.Disk) string {
	c.Lock()
	defer c.Unlock()

	id := id(g)
	c.disks[id] = g

	return id
}

func (c *cache) addRaid(r *types.Raid) string {
	c.Lock()
	defer c.Unlock()

	id := id(r)
	c.arrays[id] = r

	return id
}

func (c *cache) addUser(u *types.User) string {
	c.Lock()
	defer c.Unlock()

	id := id(u)
	c.users[id] = u

	return id
}

func (c *cache) addGroup(g *types.Group) string {
	c.Lock()
	defer c.Unlock()

	id := id(g)
	c.groups[id] = g

	return id
}

func id(input interface{}) string {
	b, _ := json.Marshal(input)
	return hash(string(b))
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

func castSliceInterface(i []interface{}) []string {
	var o []string
	for _, value := range i {
		o = append(o, value.(string))
	}

	return o
}

func getUInt(d *schema.ResourceData, key string) *uint {
	var uid *uint
	if value, ok := d.GetOk(key); ok {
		u := uint(value.(int))
		uid = &u
	}

	return uid
}
