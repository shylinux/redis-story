package server

import (
	"os"
	"runtime"

	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/gdb"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"path"
	"strings"
)

const REDIS = "redis"
const (
	SERVER = "server"
	CLIENT = "client"
	BENCH  = "bench"
)

var Index = &ice.Context{Name: REDIS, Help: "redis",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			"windows", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
			"darwin", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
			"linux", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		SERVER: {Name: "server port=auto auto 启动:button 编译:button 下载:button", Help: "服务器", Action: map[string]*ice.Action{
			"install": {Name: "install", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy("web.code.install", "download", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"compile": {Name: "compile", Help: "编译", Hand: func(m *ice.Message, arg ...string) {
				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				m.Option(cli.CMD_DIR, path.Join("usr/install", name))
				m.Cmdy(cli.SYSTEM, "make")
			}},
			gdb.START: {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				if m.Option("port") == "" {
					m.Option("port", m.Cmdx(tcp.PORT, "get"))
				}
				p := kit.Format("var/daemon/%s", m.Option("port"))
				os.MkdirAll(path.Join(p, "logs"), ice.MOD_DIR)
				os.MkdirAll(path.Join(p, "bin"), ice.MOD_DIR)
				os.MkdirAll(p, ice.MOD_DIR)

				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join("usr/install", name, "src/redis-cli"), path.Join(p, "bin"))
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join("usr/install", name, "src/redis-server"), path.Join(p, "bin"))
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join("usr/install", name, "src/redis-benchmark"), path.Join(p, "bin"))

				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.DAEMON, "bin/redis-server", "--port", m.Option("port"))
			}},
			gdb.STOP: {Name: "stop", Help: "停止", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(cli.SYSTEM, "kill", m.Option("pid"))
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 {
				m.Cmd(cli.DAEMON).Table(func(index int, value map[string]string, head []string) {
					if strings.HasPrefix(value["name"], "bin/redis") {
						m.Push("time", value["time"])
						m.Push("port", path.Base(value["dir"]))
						m.Push("status", value["status"])
						m.Push("name", value["name"])
						m.Push("pid", value["pid"])

					}
				})
				m.PushAction("启动", "停止")
				return
			}
			m.Cmdy(cli.SYSTEM, "bin/redis-cli", "-p", arg[0], "info")
		}},
		CLIENT: {Name: "client port cmd key arg 执行:button 返回:button", Help: "客户端", Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 || arg[0] == "" || arg[0] == "0" {
				m.Cmdy(SERVER)
				return
			}
			m.Cmdy(cli.SYSTEM, "bin/redis-cli", "-p", arg[0], arg[1:])
			m.Set(ice.MSG_APPEND)
		}},
		BENCH: {Name: "bench port cmd 执行:button 返回:button", Help: "压测", Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 || arg[0] == "" || arg[0] == "0" {
				m.Cmdy(SERVER)
				return
			}
			for _, k := range arg[1:] {
				m.Push("cmd", k)
				m.Push("res", m.Cmdx(cli.SYSTEM, "bin/redis-benchmark", "-p", arg[0], "-t", k))
			}
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
