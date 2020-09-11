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
		CLIENT: {Name: CLIENT, Help: "client", Value: kit.Data(
			kit.MDB_FIELD, "time,hash,host,port",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { m.Load() }},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) { m.Save() }},

		CLIENT: {Name: "client hash 执行:button 返回 连接 cmd:textarea", Help: "client", Action: map[string]*ice.Action{
			mdb.CREATE: {Name: "create host=localhost port=10000@key", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},
			mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.MODIFY, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH), arg)
			}},
			mdb.DELETE: {Name: "delete", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.DELETE, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH))
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				switch arg[0] {
				case "port":
					m.Cmdy(server.SERVER)
				}
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if m.Option(mdb.FIELDS, m.Conf(CLIENT, kit.META_FIELD)); len(arg) == 0 || arg[0] == "" {
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
				m.PushAction("删除")
				return
			}

			msg := m.Cmd(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, arg[0])
			if redis, err := NewClient(kit.Format("%s:%s", msg.Append("host"), msg.Append("port"))); m.Assert(err) {
				defer redis.Close()

				if res, err := redis.Do(kit.Split(kit.Select("info CPU", arg, 1))...); m.Assert(err) {
					m.Echo("%v", res)
				}
			}
		}},
	},
}

func init() { server.Index.Merge(Index, nil) }
