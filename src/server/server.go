package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"path"
)

const SERVER = "server"
const REDIS = "redis"

var Index = &ice.Context{Name: REDIS, Help: "redis",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			"source", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		SERVER: {Name: "server port path auto 压测 启动 构建 下载", Help: "服务器", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(SERVER, kit.META_SOURCE))
			}},
			"build": {Name: "build", Help: "构建", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv("prepare", func(p string) {})
				m.Cmdy(code.INSTALL, "build", m.Conf(SERVER, kit.META_SOURCE))
			}},
			"start": {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				pp := ""
				m.Optionv("prepare", func(p string) []string {
					pp = p
					return []string{"--port", path.Base(p)}
				})
				m.Cmdy(code.INSTALL, "start", m.Conf(SERVER, kit.META_SOURCE), "bin/redis-server")

				m.Sleep("1s")
				m.Cmd("client", mdb.CREATE, kit.MDB_HOST, "localhost", kit.MDB_PORT, path.Base(pp))
			}},
			"bench": {Name: "bench nconn=100 nreq=1000 host=localhost port=10001@key cmds=", Help: "压测", Hand: func(m *ice.Message, arg ...string) {
				for _, k := range kit.Split(kit.Select(m.Option("cmds"), "get,set")) {
					if s, e := Bench(kit.Int64(m.Option("nconn")), kit.Int64(m.Option("nreq")), []string{m.Option("host") + ":" + m.Option("port")}, []string{k}, func(cmd string, arg []interface{}, res interface{}) {

					}); m.Assert(e) {
						m.Push("cmd", k)
						m.Push("cost", kit.Format(s.EndTime.Sub(s.BeginTime)))
						m.Push("nreqs", s.NReq)
						m.Push("nerr", s.NErr)
						m.Push("qps", s.QPS)
						m.Push("avg", s.AVG)
					}
				}
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				switch arg[0] {
				case kit.MDB_PORT:
					m.Cmdy(SERVER)
				}
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Cmdy(code.INSTALL, path.Base(m.Conf(SERVER, kit.META_SOURCE)), arg)
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
