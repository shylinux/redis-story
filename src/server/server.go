package server

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"

	"shylinux.com/x/redis-story/src/client"
)

type server struct {
	ice.Code
	source  string `data:"http://mirrors.tencent.com/macports/distfiles/redis/redis-5.0.8.tar.gz"`
	start   string `name:"start port*=10001 password cluster-enabled=yes,no"`
	cluster string `name:"start count=6 replicas*=1 port*=10001 password=123"`
}

func (s server) Init(m *ice.Message, arg ...string)  { m.PackageCreateSource("redis") }
func (s server) Build(m *ice.Message, arg ...string) { s.Code.Build(m, "", "MALLOC=libc") }
func (s server) Start(m *ice.Message, arg ...string) {
	password := m.OptionDefault(aaa.PASSWORD, kit.Hashs(mdb.UNIQ))
	s.Code.Start(m, "", "bin/redis-server", func(p string) []string {
		return append([]string{
			"--port", path.Base(p), "--logfile", "redis.log", "--requirepass", password,
		}, s.Code.Args(m, "cluster-enabled")...)
	})
	m.Cmd(client.Client{}, mdb.CREATE, aaa.SESS, kit.Hashs(tcp.LOCALHOST, m.Option(tcp.PORT)),
		tcp.HOST, "127.0.0.1", tcp.PORT, m.Option(tcp.PORT), aaa.PASSWORD, password)
}
func (s server) Cluster(m *ice.Message, arg ...string) {
	m.GoToast(func(toast func(string, int, int)) {
		count, port := kit.Int(m.Option(mdb.COUNT)), kit.Int(m.Option(tcp.PORT))
		cmd, args := "redis-cli", []string{}
		for i := 0; i < count; i++ {
			toast(kit.Format("start %d", i), i, count)
			m.SpawnSilent().Cmd(s, s.Start, tcp.PORT, port+i)
			cmd = path.Join(ice.USR_LOCAL_DAEMON, kit.Format(port), "bin/redis-cli")
			args = append(args, tcp.HostPort("127.0.0.1", kit.Format(port+i)))
		}
		m.SystemCmd(cmd, "-a", m.Option(aaa.PASSWORD), "--cluster", mdb.CREATE, args, "--cluster-replicas", m.Option("replicas"), "--cluster-yes")
		m.Sleep("3s").Cmd(client.Client{}, client.Client{}.Scan)
	})
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...).Action(s.Start, s.Cluster)
}

func init() { ice.CodeModCmd(server{}) }
