package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/web"
	kit "shylinux.com/x/toolkits"
)

type studio struct {
	Client
	online string `data:"true"`
	tools  string `data:"web.code.redis.server,web.code.redis.client,web.code.redis.cluster"`
	create string `name:"create sess*=biz host*=localhost port*=10001 password*=root"`
	list   string `name:"list list" icon:"redis.png"`
}

func (s studio) Init(m *ice.Message, arg ...string) {
	web.AddPortalProduct(m.Message, "Redis Studio", `
一款网页版的 Redis 工作台，用来下载安装 Redis、组建集群、操作数据等。
`, 10.0)
}
func (s studio) Create(m *ice.Message, arg ...string) { m.Cmdy(s.Client, m.ActionKey(), arg) }
func (s studio) Remove(m *ice.Message, arg ...string) { m.Cmdy(s.Client, m.ActionKey(), arg) }
func (s studio) List(m *ice.Message, arg ...string) {
	m.Cmdy(s.Client, arg).PushAction(s.Remove).Action(s.Create)
	kit.If(len(m.Appendv(aaa.SESS)) > 0, func() { m.Display("") })
}

func init() { ice.CodeModCmd(studio{}) }
