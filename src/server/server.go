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
	source string `data:"http://mirrors.tencent.com/macports/distfiles/redis/redis-5.0.8.tar.gz"`
	start  string `name:"start port*=10001 password cluster-enabled=yes,no"`
}

func (s server) Init(m *ice.Message, arg ...string)  { m.PackageCreateSource("redis") }
func (s server) Build(m *ice.Message, arg ...string) { s.Code.Build(m, "", "MALLOC=libc") }
func (s server) Start(m *ice.Message, arg ...string) {
	password := m.OptionDefault(aaa.PASSWORD, kit.Hashs(mdb.UNIQ))
	s.Code.Start(m, "", "bin/redis-server", func(p string) []string {
		return append([]string{
			"--port", path.Base(p), "--requirepass", password,
			"--logfile", "redis.log",
		}, s.Code.Args(m, "cluster-enabled")...)
	})
	m.Cmd(client.Client{}, mdb.CREATE, aaa.SESS, kit.Hashs(tcp.LOCALHOST, m.Option(tcp.PORT)),
		tcp.HOST, "127.0.0.1", tcp.PORT, m.Option(tcp.PORT), aaa.PASSWORD, password)
}
func (s server) List(m *ice.Message, arg ...string) { s.Code.List(m, "", arg...) }

func init() { ice.CodeModCmd(server{}) }
