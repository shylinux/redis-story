package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/gdb"
	"github.com/shylinux/icebergs/base/nfs"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/base/web"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"net/http"
	"os"
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
			"source", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		SERVER: {Name: "server port 查看:button=auto 启动:button 编译:button 下载:button", Help: "服务器", Action: map[string]*ice.Action{
			"install": {Name: "install", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				// 下载
				source := m.Conf(SERVER, "meta.source")
				p := path.Join(m.Conf("web.code._install", "meta.path"), path.Base(source))
				if _, e := os.Stat(p); e != nil {
					msg := m.Cmd(web.SPIDE, "dev", web.CACHE, http.MethodGet, source)
					m.Cmd(web.CACHE, web.WATCH, msg.Append(web.DATA), p)
				}

				// 解压
				m.Option(cli.CMD_DIR, m.Conf("web.code._install", "meta.path"))
				m.Cmd(cli.SYSTEM, "tar", "xvf", path.Base(source))
				m.Echo(p)
			}},
			"compile": {Name: "compile", Help: "编译", Hand: func(m *ice.Message, arg ...string) {
				// 编译
				source := m.Conf(SERVER, "meta.source")
				m.Option(cli.CMD_DIR, path.Join(m.Conf("web.code._install", "meta.path"), strings.TrimSuffix(path.Base(source), ".tar.gz")))
				m.Cmdy(cli.SYSTEM, "make")

				// 链接
				m.Cmd(nfs.LINK, "bin/redis-cli", path.Join(m.Option(cli.CMD_DIR), "src/redis-cli"))
				m.Cmd(nfs.LINK, "bin/redis-server", path.Join(m.Option(cli.CMD_DIR), "src/redis-server"))
				m.Cmd(nfs.LINK, "bin/redis-benchmark", path.Join(m.Option(cli.CMD_DIR), "src/redis-benchmark"))
			}},
			gdb.START: {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				if m.Option("port") == "" {
					m.Option("port", m.Cmdx(tcp.PORT, "get"))
				}
				m.Cmdy(cli.DAEMON, "bin/redis-server", "--port", m.Option("port"))
			}},
			gdb.STOP: {Name: "stop", Help: "停止", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(cli.SYSTEM, "kill", m.Option("PID"))
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Split(m.Cmdx(cli.SYSTEM, "sh", "-c", "ps aux|grep redis-server|grep -v grep"),
				"USER PID CPU MEM VSZ RSS TTY STAT START TIME COMMAND", " ", "\n")
			m.Table(func(index int, value map[string]string, head []string) {
				if ls := kit.Split(value["COMMAND"], " ", " "); len(ls) > 1 {
					if ls = kit.Split(ls[1], ":", ":"); len(ls) > 1 {
						m.Push("port", ls[1])
						return
					}
				}
				m.Push("port", "8397")
			})
			m.Appendv(ice.MSG_APPEND, "USER", "PID", "STAT", "START", "port", "COMMAND")
			m.PushAction("停止")
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
