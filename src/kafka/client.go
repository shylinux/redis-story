package kafka

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type client struct {
	list string `name:"list path auto" help:"客户端"`
}

func (c client) List(m *ice.Message, arg ...string) {
	m.Option(cli.CMD_DIR, path.Join(m.Conf(code.INSTALL, kit.META_PATH), kit.TrimExt(m.Conf(tcp.SERVER, kit.Keym(cli.LINUX)))))
	m.Cmdy(cli.SYSTEM, "bin/kafka-topics.sh", "--list", "--zookeeper", "localhost:2181")
}

func init() { ice.Cmd("web.code.kafka.client", client{}) }
