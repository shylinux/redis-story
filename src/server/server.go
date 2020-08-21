package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"os"
	"path"
	"runtime"
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

		SERVER: {Name: "server port=auto auto 启动:button 编译:button 下载:button cmd:textarea", Help: "服务器", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"compile": {Name: "compile", Help: "编译", Hand: func(m *ice.Message, arg ...string) {
				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				m.Option(cli.CMD_DIR, path.Join(m.Conf(code.INSTALL, kit.META_PATH), name))
				m.Cmdy(cli.SYSTEM, "make", "-j8")
			}},
			"start": {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				// 分配
				port, p := "", ""
				for {
					port = m.Cmdx(tcp.PORT, "select", port)
					p = path.Join(m.Conf(cli.DAEMON, kit.META_PATH), port)
					if _, e := os.Stat(p); e != nil && os.IsNotExist(e) {
						break
					}
					port = kit.Format(kit.Int(port) + 1)
				}
				os.MkdirAll(path.Join(p, "logs"), ice.MOD_DIR)
				os.MkdirAll(path.Join(p, "bin"), ice.MOD_DIR)
				os.MkdirAll(p, ice.MOD_DIR)

				// 复制
				name := path.Base(strings.TrimSuffix(strings.TrimSuffix(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)), ".tar.gz"), "zip"))
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join(m.Conf(code.INSTALL, kit.META_PATH), name, "src/redis-cli"), path.Join(p, "bin"))
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join(m.Conf(code.INSTALL, kit.META_PATH), name, "src/redis-server"), path.Join(p, "bin"))
				m.Cmd(cli.SYSTEM, "cp", "-r", path.Join(m.Conf(code.INSTALL, kit.META_PATH), name, "src/redis-benchmark"), path.Join(p, "bin"))

				// 启动
				m.Option(cli.CMD_DIR, p)
				m.Cmdy(cli.DAEMON, "bin/redis-server", "--port", port)
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) > 0 && arg[0] != "" {
				m.Cmdy(cli.SYSTEM, "bin/redis-cli", "-p", arg[0], kit.Split(kit.Select("info", arg, 1)))
				return
			}

			m.Cmd(cli.DAEMON).Table(func(index int, value map[string]string, head []string) {
				if strings.HasPrefix(value[kit.MDB_NAME], "bin/redis") {
					m.Push(kit.MDB_TIME, value[kit.MDB_TIME])
					m.Push(kit.MDB_PORT, path.Base(value[kit.MDB_DIR]))
					m.Push(kit.MDB_DIR, value[kit.MDB_DIR])
					m.Push(kit.MDB_STATUS, value[kit.MDB_STATUS])
					m.Push(kit.MDB_PID, value[kit.MDB_PID])
					m.Push(kit.MDB_NAME, value[kit.MDB_NAME])
				}
			})
			m.Sort("time", "time_r")
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
