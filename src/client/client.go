package client

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/core/wiki"
	kit "github.com/shylinux/toolkits"
)

const (
	CLIENT = "client"
)

var Index = &ice.Context{Name: CLIENT, Help: "client",
	Configs: map[string]*ice.Config{
		CLIENT: {Name: CLIENT, Help: "client", Value: kit.Data(kit.MDB_SHORT, "name")},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		"do": {Name: "do address cmd key value", Help: "do", Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
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
