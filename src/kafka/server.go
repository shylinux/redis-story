package kafka

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code

	source string `data:"http://mirrors.tencent.com/apache/kafka/2.8.1/kafka-2.8.1-src.tgz"`
	linux  string `data:"http://mirrors.tencent.com/apache/kafka/2.8.1/kafka_2.12-2.8.1.tgz"`

	install string `name:"install" help:"安装"`
	list    string `name:"list path auto start install download" help:"服务器"`
}

func (s server) Download(m *ice.Message, arg ...string) {
	s.Code.Download(m, m.Config(nfs.SOURCE), arg...)
}
func (s server) Install(m *ice.Message, arg ...string) {
	s.Code.Download(m, m.Config(cli.LINUX), arg...)
}
func (s server) Start(m *ice.Message, arg ...string) {
	p := path.Join(m.Conf(code.INSTALL, kit.Keym(nfs.PATH)), kit.TrimExt(m.Config(cli.LINUX)))
	s.Code.Daemon(m, p, "bin/zookeeper-server-start.sh", "config/zookeeper.properties")
	m.Sleep3s()
	s.Code.Daemon(m, p, "bin/kafka-server-start.sh", "config/server.properties")
}
func (s server) List(m *ice.Message, arg ...string) {
	m.OptionFields("time,status,pid,cmd,dir")
	m.Debug("what %v", 123)
	m.Cmd(cli.DAEMON).Table(func(index int, value map[string]string, head []string) {
		m.Debug("what %v", 123)
		if strings.Contains(value[ice.CMD], "bin/kafka") || strings.Contains(value[ice.CMD], "bin/zookeeper") {
			m.Push("", value, head)
		}
	})
}

func init() { ice.CodeCtxCmd(server{}) }
