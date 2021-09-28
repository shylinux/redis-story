package kafka

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
)

type client struct {
	list string `name:"list path auto" help:"客户端"`
}

func (c client) List(m *ice.Message, arg ...string) {
	_dir(m)
	m.Cmdy(cli.SYSTEM, "bin/kafka-topics.sh", "--list", "--zookeeper", "localhost:2181")
}

func init() { ice.Cmd("web.code.kafka.client", client{}) }
