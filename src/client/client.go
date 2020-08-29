package client

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/redis-story/src/server"
	kit "github.com/shylinux/toolkits"
)

const (
	CLIENT = "client"
)

var Index = &ice.Context{Name: CLIENT, Help: "client",
	Configs: map[string]*ice.Config{
		CLIENT: {Name: CLIENT, Help: "client", Value: kit.Data()},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { m.Load() }},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { m.Save() }},

		CLIENT: {Name: "client hash=@key 执行:button 连接 cmd:textarea", Help: "client", Meta: kit.Dict(
			"连接", kit.List(
				kit.MDB_INPUT, "text", "name", "hostport", "value", "localhost:10035",
			),
		), Action: map[string]*ice.Action{
			"connect": {Name: "connect", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Option(mdb.FIELDS, "time,hash,hostport")
			if len(arg) == 0 || arg[0] == "" {
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
				return
			}

			msg := m.Cmd(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH, arg[0])
			if redis, err := NewClient(msg.Append("hostport")); m.Assert(err) {
				defer redis.Close()

				if res, err := redis.Do(kit.Split(kit.Select("info CPU", arg, 1))...); m.Assert(err) {
					m.Echo("%v", res)
				}
			}
		}},
	},
}

func init() { server.Index.Merge(Index, nil) }
