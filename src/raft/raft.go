package raft

import (
	"github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/core/wiki"
	"github.com/shylinux/toolkits"
)

var Index = &ice.Context{Name: "raft", Help: "raft",
	Caches: map[string]*ice.Cache{},
	Configs: map[string]*ice.Config{
		"raft": {Name: "raft", Help: "raft", Value: kit.Data(kit.MDB_SHORT, "name")},
	},
	Commands: map[string]*ice.Command{
		ice.ICE_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.ICE_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		"raft": {Name: "raft", Help: "raft", Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Echo("hello world")
		}},
	},
}

func init() { wiki.Index.Register(Index, nil) }
