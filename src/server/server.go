package server

import (
	"path"
	"time"

	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/gdb"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/icebergs/base/web"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"
)

const (
	REDIS_SERVER_START = "redis_server_start"
)
const (
	SERVER = "server"
)
const REDIS = "redis"

var Index = &ice.Context{Name: REDIS, Help: "redis",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			"source", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		SERVER: {Name: "server port path auto start build download", Help: "服务器", Action: map[string]*ice.Action{
			web.DOWNLOAD: {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(SERVER, kit.META_SOURCE))
			}},
			gdb.BUILD: {Name: gdb.BUILD, Help: "构建", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv(code.PREPARE, func(p string) {})
				m.Cmdy(code.INSTALL, gdb.BUILD, m.Conf(SERVER, kit.META_SOURCE))
			}},
			gdb.START: {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				pp := ""
				m.Optionv(code.PREPARE, func(p string) []string {
					pp = p
					return []string{"--port", path.Base(p)}
				})
				m.Cmdy(code.INSTALL, gdb.START, m.Conf(SERVER, kit.META_SOURCE), "bin/redis-server")

				m.Sleep("1s")
				m.Event(REDIS_SERVER_START, tcp.HOST, tcp.LOCALHOST, tcp.PORT, path.Base(pp))
			}},
			gdb.BENCH: {Name: "bench nconn=100 nreq=1000 cmdList=", Help: "压测", Hand: func(m *ice.Message, arg ...string) {
				for _, k := range kit.Split(kit.Select(m.Option("cmdList"), "get,set")) {
					begin := time.Now()
					if s, e := Bench(kit.Int64(m.Option("nconn")), kit.Int64(m.Option("nreq")), []string{"localhost:" + m.Option(tcp.PORT)}, []string{k}, func(cmd string, arg []interface{}, res interface{}) {
						// 检查结果

					}); m.Assert(e) {
						m.Push("time", begin)
						m.Push("cmd", k)
						m.Push("cost", kit.Format(s.EndTime.Sub(s.BeginTime)))
						m.Push("nreqs", s.NReq)
						m.Push("nerr", s.NErr)
						m.Push("qps", s.QPS)
						m.Push("avg", s.AVG)
					}
				}
				m.Option(ice.MSG_PROCESS, ice.PROCESS_INNER)
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				switch arg[0] {
				case tcp.PORT:
					m.Cmdy(SERVER).Appendv(ice.MSG_APPEND, tcp.PORT, kit.MDB_STATUS, kit.MDB_TIME)
				}
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Cmdy(code.INSTALL, path.Base(m.Conf(SERVER, kit.META_SOURCE)), arg)
			if len(arg) == 0 {
				m.PushAction(gdb.BENCH)
			}
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
