package kafka

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	source string `data:"http://mirrors.tencent.com/apache/kafka/2.8.1/kafka-2.8.1-src.tgz"`
	linux  string `data:"http://mirrors.tencent.com/apache/kafka/2.8.1/kafka_2.12-2.8.1.tgz"`

	download string `name:"download" help:"下载"`
	install  string `name:"install" help:"安装"`
	start    string `name:"start" help:"启动"`
	list     string `name:"list path auto start install download" help:"服务器"`
}

func (s server) Download(m *ice.Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(tcp.SERVER, kit.Keym(cli.SOURCE)))
}
func (s server) Install(m *ice.Message, arg ...string) {
	m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(tcp.SERVER, kit.Keym(cli.LINUX)))
}
func (s server) Start(m *ice.Message, arg ...string) {
	m.Option(cli.CMD_DIR, path.Join(m.Conf(code.INSTALL, kit.META_PATH), kit.TrimExt(m.Conf(tcp.SERVER, kit.Keym(cli.LINUX)))))
	m.Cmdy(cli.DAEMON, "bin/zookeeper-server-start.sh", "config/zookeeper.properties")
	m.Sleep("1s")
	m.Cmdy(cli.DAEMON, "bin/kafka-server-start.sh", "config/server.properties")
}
func (s server) List(m *ice.Message, arg ...string) {
	m.OptionFields("time,status,pid,cmd,dir")
	m.Cmd(cli.DAEMON).Table(func(index int, value map[string]string, head []string) {
		if strings.Contains(value["cmd"], "bin/kafka") || strings.Contains(value["cmd"], "bin/zookeeper") {
			m.Push("", value, head)
		}
	})
}

func init() { ice.Cmd("web.code.kafka.server", server{}) }
