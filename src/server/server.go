package server

import (
	"path"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code

	source string `data:"http://mirrors.tencent.com/macports/distfiles/redis/redis-5.0.8.tar.gz"`
	start  string `name:"start port password" help:"启动"`
	bench  string `name:"bench port nconn=100 nreq=1000 cmdList" help:"压测"`
}

func (s server) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		if s.List(m); m.Length() > 0 {
			m.Cut("port,status,time")
		} else {
			m.Cmdy(tcp.PORT)
		}
	}
}
func (s server) Download(m *ice.Message, arg ...string) {
	s.Code.Download(m, m.Config(nfs.SOURCE), arg...)
}
func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Prepare(m, func(p string) {})
	s.Code.Build(m, m.Config(nfs.SOURCE), arg...)
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Prepare(m, func(p string) []string {
		return []string{"--port", path.Base(p), "--requirepass", m.Option(aaa.PASSWORD)}
	})
	s.Code.Start(m, m.Config(nfs.SOURCE), "bin/redis-server")
}
func (s server) Bench(m *ice.Message, arg ...string) {
	for _, k := range kit.Split(kit.Select(m.Option("cmdList"), "get,set")) {
		begin := time.Now()
		if s, e := Bench(kit.Int64(m.Option("nconn")), kit.Int64(m.Option("nreq")), []string{tcp.LOCALHOST + ice.FS + m.Option(tcp.PORT)}, []string{k}, func(cmd string, arg []interface{}, res interface{}) {
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
	if s.Code.List(m, m.Config(nfs.SOURCE), arg...); len(arg) == 0 {
		m.PushAction(s.Bench)
	}
}
func init() { ice.CodeModCmd(server{}) }
