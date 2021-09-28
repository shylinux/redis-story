package server

import (
	"path"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	source string `data:"http://mirrors.tencent.com/macports/distfiles/redis/redis-5.0.8.tar.gz"`

	inputs   string `name:"inputs" help:"补全"`
	download string `name:"download" help:"下载"`
	build    string `name:"build" help:"构建"`
	start    string `name:"start" help:"启动"`
	bench    string `name:"bench port nconn=100 nreq=1000 cmdList" help:"压测"`
	list     string `name:"list port path auto bench start build download" help:"服务器"`
}

func (s server) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER)
	}
}
func (s server) Download(m *ice.Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(tcp.SERVER, kit.META_SOURCE))
}
func (s server) Build(m *ice.Message, arg ...string) {
	m.Optionv(code.PREPARE, func(p string) {})
	m.Cmdy(code.INSTALL, cli.BUILD, m.Conf(tcp.SERVER, kit.META_SOURCE))
}
func (s server) Start(m *ice.Message, arg ...string) {
	m.Optionv(code.PREPARE, func(p string) []string { return []string{"--port", path.Base(p)} })
	m.Cmdy(code.INSTALL, cli.START, m.Conf(tcp.SERVER, kit.META_SOURCE), "bin/redis-server")
	// m.Sleep("1s").Event(REDIS_SERVER_START, tcp.HOST, tcp.LOCALHOST, tcp.PORT, path.Base(m.Option(cli.CMD_DIR)))
}
func (s server) Bench(m *ice.Message, arg ...string) {
	for _, k := range kit.Split(kit.Select(m.Option("cmdList"), "get,set")) {
		begin := time.Now()
		if s, e := Bench(kit.Int64(m.Option("nconn")), kit.Int64(m.Option("nreq")), []string{tcp.LOCALHOST + ":" + m.Option(tcp.PORT)}, []string{k}, func(cmd string, arg []interface{}, res interface{}) {
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
	m.ProcessInner()
}
func (s server) List(m *ice.Message, arg ...string) {
	if m.Cmdy(code.INSTALL, path.Base(m.Conf(tcp.SERVER, kit.META_SOURCE)), arg); len(arg) == 0 {
		m.PushAction(cli.BENCH)
	}
}
func init() { ice.Cmd("web.code.redis.server", server{}) }
