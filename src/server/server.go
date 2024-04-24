package server

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"

	"shylinux.com/x/redis-story/src/client"
)

const (
	BIN_REDIS_CLI    = "bin/redis-cli"
	BIN_REDIS_SERVER = "bin/redis-server"
	CLUSTER_ENABLED  = "cluster-enabled"
)

type server struct {
	ice.Code
	client  client.Client
	source  string `data:"http://mirrors.tencent.com/macports/distfiles/redis/redis-5.0.8.tar.gz"`
	start   string `name:"start port*=10001 host password cluster-enabled=no,yes"`
	cluster string `name:"cluster port*=10001 host password count*=6 replicas*=1"`
}

func (s server) Init(m *ice.Message, arg ...string)  { m.PackageCreateSource("redis") }
func (s server) Build(m *ice.Message, arg ...string) { s.Code.Build(m, "", "MALLOC=libc") }
func (s server) Start(m *ice.Message, arg ...string) {
	bind := m.OptionDefault(tcp.HOST, "127.0.0.1")
	password := m.OptionDefault(aaa.PASSWORD, kit.HashsUniq())
	s.Code.Start(m, "", BIN_REDIS_SERVER, func(p string) []string {
		m.Cmd(nfs.SAVE, path.Join(p, "redis.conf"), kit.Format(`
port %s
bind %s
requirepass %s
cluster-enabled %s
logfile redis.log
`, path.Base(p), bind, password, m.Option(CLUSTER_ENABLED)))
		return kit.Simple("redis.conf")
	})
	m.Cmd(s.client, s.client.Create, aaa.SESS, kit.Hashs(bind, m.Option(tcp.PORT)), tcp.HOST, bind, tcp.PORT, m.Option(tcp.PORT), aaa.PASSWORD, password)
	kit.If(m.Option(CLUSTER_ENABLED) == "no", func() { m.Sleep("3s").Cmd(s.client, s.client.Scan) })
}
func (s server) Cluster(m *ice.Message, arg ...string) {
	m.GoToast(func(toast func(string, int, int)) {
		bind := m.OptionDefault(tcp.HOST, "127.0.0.1")
		password := m.OptionDefault(aaa.PASSWORD, kit.HashsUniq())
		port, count := kit.Int(m.Option(tcp.PORT)), kit.Int(m.Option(mdb.COUNT))
		cmd, args := path.Base(BIN_REDIS_CLI), []string{}
		for i := 0; i < count; i++ {
			toast(kit.Format("start %d", port+i), i, count)
			m.SpawnSilent().Cmd(s, s.Start, tcp.PORT, port+i, CLUSTER_ENABLED, "yes")
			cmd = path.Join(ice.USR_LOCAL_DAEMON, kit.Format(port), BIN_REDIS_CLI)
			args = append(args, tcp.HostPort(bind, kit.Format(port+i)))
		}
		m.SystemCmd(cmd, "-a", password, "--cluster", mdb.CREATE, args, "--cluster-replicas", m.Option("replicas"), "--cluster-yes")
		m.Sleep("5s").Cmd(s.client, s.client.Scan)
	})
}
func (s server) List(m *ice.Message, arg ...string) {
	s.Code.List(m, "", arg...).Action(s.Start, s.Cluster)
}

func init() { ice.CodeModCmd(server{}) }
