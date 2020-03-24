package server

import (
	"github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/core/wiki"
	"github.com/shylinux/toolkits"
)

var Index = &ice.Context{Name: "server", Help: "服务器",
	Caches: map[string]*ice.Cache{},
	Configs: map[string]*ice.Config{
		"server": {Name: "server", Help: "服务器", Value: kit.Data(kit.MDB_SHORT, "name", "pid")},
	},
	Commands: map[string]*ice.Command{
		ice.ICE_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.ICE_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		"server": {Name: "server", Help: "服务器", List: kit.List(
			kit.MDB_INPUT, "button", "name", "启动",
		), Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			msg := m.Cmd("web.code.docker.command", m.Conf(ice.CLI_RUNTIME, "node.name"), "ps")
			msg.Split(msg.Result(), "", "", "\n").Table(func(index int, value map[string]string, head []string) {
				if value["COMMAND"] == "redis-server" {
					m.Push("CONTAINER", m.Conf(ice.CLI_RUNTIME, "node.name"))
					m.Push("COMMAND", value["COMMAND"])
					m.Push("PID", value["PID"])
				}
			})

			if m.Append("PID") == "" {
				m.Option("cmd_type", "daemon")
				m.Cmdy("web.code.docker.command", m.Conf(ice.CLI_RUNTIME, "node.name"), "redis-server")
			}
		}},
		"client": {Name: "client", Help: "命令行", List: kit.List(
			kit.MDB_INPUT, "text", "name", "cmd", "value", "get",
			kit.MDB_INPUT, "text", "name", "name", "value", "hi",
			kit.MDB_INPUT, "button", "name", "查看",
		), Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Cmdy("web.code.docker.command", m.Conf(ice.CLI_RUNTIME, "node.name"), "redis-cli", arg)
		}},
	},
}

func init() { wiki.Index.Register(Index, nil) }
