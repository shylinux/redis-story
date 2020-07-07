package data

import (
	"github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/core/wiki"
	"github.com/shylinux/toolkits"
)

var Index = &ice.Context{Name: "data", Help: "data",
	Caches: map[string]*ice.Cache{},
	Configs: map[string]*ice.Config{
		"data": {Name: "data", Help: "data", Value: kit.Data(kit.MDB_SHORT, "name")},
	},
	Commands: map[string]*ice.Command{
		ice.ICE_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.ICE_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		"data": {Name: "data", Help: "data", Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Echo("hello world")
		}},
	},
}

func init() { wiki.Index.Register(Index, nil) }
