package kafka

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/web"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"
)

const KAFKA = "kafka"

var Index = &ice.Context{Name: "kafka", Help: "消息队列", Configs: map[string]*ice.Config{
	KAFKA: {Name: "kafka", Help: "消息队列", Value: kit.Data(
		cli.SOURCE, "https://mirror-hk.koddos.net/apache/kafka/2.8.0/kafka-2.8.0-src.tgz",
		cli.LINUX, "https://mirror-hk.koddos.net/apache/kafka/2.8.0/kafka_2.12-2.8.0.tgz",
	)},
}, Commands: map[string]*ice.Command{
	KAFKA: {Name: "kafka port path auto download", Help: "消息队列", Action: map[string]*ice.Action{
		web.DOWNLOAD: {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
			m.Cmdy(code.INSTALL, web.DOWNLOAD, m.Conf(KAFKA, kit.Keym(cli.LINUX)))
		}},
	}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {

	}},
}}

func init() { code.Index.Register(Index, nil) }
