package pulsar

import (
	"path"
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/nfs"
	kit "shylinux.com/x/toolkits"
)

type server struct {
	ice.Code
	linux     string `data:"https://mirrors.tencent.com/apache/pulsar/pulsar-2.10.2/apache-pulsar-2.10.2-bin.tar.gz"`
	darwin    string `data:"https://mirrors.tencent.com/apache/pulsar/pulsar-2.10.2/apache-pulsar-2.10.2-bin.tar.gz"`
	source    string `data:"https://mirrors.tencent.com/apache/pulsar/pulsar-2.10.2/apache-pulsar-2.10.2-src.tar.gz"`
	action    string `data:"listTopic,addTopic"`
	listTopic string `name:"listTopic" help:"主题列表"`
	addTopic  string `name:"addTopic topic=TASK_AGENT" help:"添加主题"`
	start     string `name:"start port=10003" help:"启动"`
	list      string `name:"list port path auto start install" help:"消息队列"`
}

func (s server) ListTopic(m *ice.Message, arg ...string) {
	s.Code.System(m, m.Option(nfs.DIR), "bin/pulsar-admin", "topics", "list", "public/default")
}
func (s server) AddTopic(m *ice.Message, arg ...string) {
}
func (s server) Start(m *ice.Message, arg ...string) {
	s.Code.Start(m, "", "bin/pulsar", "standalone", func(p string, port int) {
		nfs.Rewrite(m.Message, path.Join(p, "conf/standalone.conf"), func(text string) string {
			switch text {
			case "webServicePort=8080":
				text = strings.Replace(text, "8080", kit.Format(port+10000), 1)
			case "brokerServicePort=6650":
				text = strings.Replace(text, "6650", kit.Format(port), 1)
			}
			return text
		})
		nfs.Rewrite(m.Message, path.Join(p, "conf/client.conf"), func(text string) string {
			switch text {
			case "webServiceUrl=http://localhost:8080/":
				text = strings.Replace(text, "8080", kit.Format(port+10000), 1)
			case "brokerServiceUrl=pulsar://localhost:6650/":
				text = strings.Replace(text, "6650", kit.Format(port), 1)
			}
			return text
		})
		nfs.Rewrite(m.Message, path.Join(p, "conf/zookeeper.conf"), func(text string) string {
			switch text {
			case "clientPort=2181":
				text = strings.Replace(text, "2181", kit.Format(port+20000), 1)
			case "admin.serverPort=9990":
				text = strings.Replace(text, "9990", kit.Format(port+30000), 1)
			}
			return text
		})
	})
}
func (s server) List(m *ice.Message, arg ...string) { s.Code.List(m, "pulsar", arg...) }

func init() { ice.CodeCtxCmd(server{}) }
