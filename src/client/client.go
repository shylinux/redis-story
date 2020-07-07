package client

import (
	"github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/core/wiki"
	"github.com/shylinux/toolkits"
)

var Index = &ice.Context{Name: "client", Help: "client",
	Caches: map[string]*ice.Cache{},
	Configs: map[string]*ice.Config{
		"client": {Name: "client", Help: "client", Value: kit.Data(kit.MDB_SHORT, "name")},
	},
	Commands: map[string]*ice.Command{
		ice.ICE_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.ICE_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		"do": {Name: "do address cmd arg...", Help: "do", Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if redis, err := NewClient(arg[0]); m.Assert(err) {
				defer redis.Close()

				if res, err := redis.Do(arg[1:]...); m.Assert(err) {
					m.Echo("%v", res)
				}
			}
		}},
	},
}

func init() { wiki.Index.Register(Index, nil) }
