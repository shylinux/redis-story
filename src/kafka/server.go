package kafka

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code

	source string `data:"http://mirrors.tencent.com/apache/kafka/2.8.1/kafka-2.8.1-src.tgz"`
	linux  string `data:"http://mirrors.tencent.com/apache/kafka/2.8.1/kafka_2.12-2.8.1.tgz"`

	listTopic string `name:"listTopic" help:"主题列表"`
	addTopic  string `name:"addTopic topic=TASK_AGENT" help:"添加主题"`

	start string `name:"start port=10002" help:"启动"`
	list  string `name:"list port path auto start install download" help:"服务器"`
}

func (s server) zkport(port string) string {
	return kit.Format(kit.Int(port) + 10000)
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "bin/kafka-server-start.sh", "config/server.properties", func(p string) {
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
	})
}
func (s server) ListTopic(m *ice.Message, arg ...string) {
	s.Code.System(m, m.Option(nfs.DIR), kit.Format("bin/kafka-topics.sh --list --zookeeper localhost:%s", s.zkport(m.Option(tcp.PORT))))
}
func (s server) AddTopic(m *ice.Message, arg ...string) {
	s.Code.System(m, m.Option(nfs.DIR), kit.Format("bin/kafka-topics.sh --create --zookeeper localhost:%s --replication-factor 1 --partitions 1 --topic %s", s.zkport(m.Option(tcp.PORT)), m.Option(TOPIC)))
}
func (s server) List(m *ice.Message, arg ...string) {
	if s.Code.List(m, "kafka-server", arg...); len(arg) == 0 {
		m.PushAction(s.ListTopic, s.AddTopic)
	}
}

func init() { ice.CodeCtxCmd(server{}) }
