package pulsar

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code

	source string `data:"https://mirrors.tencent.com/apache/pulsar/pulsar-2.10.0/apache-pulsar-2.10.0-src.tar.gz"`
	linux  string `data:"https://mirrors.tencent.com/apache/pulsar/pulsar-2.10.0/apache-pulsar-2.10.0-bin.tar.gz"`

	listTopic string `name:"listTopic" help:"主题列表"`
	addTopic  string `name:"addTopic topic=TASK_AGENT" help:"添加主题"`

	start string `name:"start port=10002" help:"启动"`
	list  string `name:"list port path auto start install download" help:"服务器"`
}

func (s server) zkport(port string) string {
	return kit.Format(kit.Int(port) + 10000)
}
func (s server) Download(m *ice.Message, arg ...string) {
	s.Code.Download(m, m.Config(nfs.SOURCE), arg...)
}
func (s server) Install(m *ice.Message, arg ...string) {
	s.Code.Download(m, m.Config(cli.LINUX), arg...)
}
func (s server) Start(m *ice.Message, arg ...string) {
	m.Option(code.INSTALL, ice.PT)
	m.Option(code.PREPARE, func(p string) []string {
		port := path.Base(p)

		m.Cmd(nfs.SAVE, path.Join(p, "config/zookeeper.properties"), kit.Format(`
clientPort=%s
dataDir=%s
maxClientCnxns=0
admin.enableServer=false
`, s.zkport(port), kit.Path(p, "var/zookeeper")))

		m.Cmd(nfs.SAVE, path.Join(p, "config/server.properties"), kit.Format(`
zookeeper.connect=localhost:%s
listeners=PLAINTEXT://:%s
log.dirs=%s
broker.id=1
num.partitions=1

num.network.threads=3
num.io.threads=8

socket.send.buffer.bytes=102400
socket.receive.buffer.bytes=102400
socket.request.max.bytes=104857600
num.recovery.threads.per.data.dir=1
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1

log.retention.hours=168
log.segment.bytes=1073741824

log.retention.check.interval.ms=300000
zookeeper.connection.timeout.ms=18000
group.initial.rebalance.delay.ms=0
`, s.zkport(port), port, kit.Path(p, "var/kafka-logs")))

		s.Code.Daemon(m, p, "bin/zookeeper-server-start.sh", "config/zookeeper.properties")
		return []string{}
	})
	s.Code.Start(m, m.Config(cli.LINUX), "bin/kafka-server-start.sh", "config/server.properties")
}
func (s server) List(m *ice.Message, arg ...string) {
	if s.Code.List(m, "kafka-server", arg...); len(arg) == 0 {
	}
}

func init() { ice.CodeCtxCmd(server{}) }
