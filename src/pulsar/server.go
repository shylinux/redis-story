package pulsar

import (
	"shylinux.com/x/ice"
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
func (s server) Start(m *ice.Message, arg ...string) {
}
func (s server) List(m *ice.Message, arg ...string) {
	if s.Code.List(m, "kafka-server", arg...); len(arg) == 0 {
	}
}

func init() { ice.CodeCtxCmd(server{}) }
