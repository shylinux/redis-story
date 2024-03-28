package server

import (
	"path"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/core/code"
	"shylinux.com/x/redis-story/src/client"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code
	source string `data:"http://mirrors.tencent.com/macports/distfiles/redis/redis-5.0.8.tar.gz"`
	action string `data:"bench"`

	start string `name:"start port*=10002 password*=demo"`
	bench string `name:"bench port*=10002 nconn*=100 nreq*=1000 cmdList" help:"压测"`
	list  string `name:"list port path auto start build download" help:"缓存"`
}

func (s server) Init(m *ice.Message, arg ...string) {
	code.PackageCreate(m.Message, nfs.SOURCE, "redis", "", "src/redis.png", s.Link(m))
}
func (s server) Build(m *ice.Message, arg ...string) {
	s.Code.Build(m, "", func(p string) {})
	m.Cmdy(nfs.DIR, path.Join(s.Path(m, ""), "_install/bin/redis-server"))
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "bin/redis-server", func(p string) []string {
		return []string{"--port", path.Base(p), "--requirepass", m.Option(aaa.PASSWORD)}
	})
}
func (s server) Bench(m *ice.Message, arg ...string) {
	defer m.ProcessInner()
	for _, k := range kit.Split(kit.Select(m.Option("cmdList"), "get,set")) {
		begin := time.Now()
		if s, e := client.Bench(m.FormatTaskMeta(), kit.Int64(m.Option("nconn")), kit.Int64(m.Option("nreq")), []string{tcp.LOCALHOST + ice.FS + m.Option(tcp.PORT)}, []string{k}, func(cmd string, arg []string, res ice.Any) {
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
}
func (s server) List(m *ice.Message, arg ...string) { s.Code.List(m, "", arg...) }

func init() { ice.CodeModCmd(server{}) }
