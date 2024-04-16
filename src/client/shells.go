package client

import (
	"path"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/ctx"
	"shylinux.com/x/icebergs/base/tcp"
	"shylinux.com/x/icebergs/base/web/html"
	"shylinux.com/x/icebergs/core/code"
	kit "shylinux.com/x/toolkits"
)

type shells struct {
	Client
	create string `name:"create sess*"`
}

func (s shells) Inputs(m *ice.Message, arg ...string) { m.Cmdy(s.Client, m.ActionKey(), arg) }
func (s shells) Create(m *ice.Message, arg ...string) *ice.Message {
	msg := m.Cmd(s.Client, m.Option(aaa.SESS))
	m.ProcessXterm("", kit.Format("%s -h %s -p %s -a %s",
		path.Join(ice.USR_LOCAL_DAEMON, msg.Append(tcp.PORT), "bin/redis-cli"),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.PASSWORD)), arg...)
	m.Push(ctx.STYLE, html.OUTPUT)
	m.Option(ice.FIELD_PREFIX, ctx.ACTION, ctx.RUN, code.XTERM)
	m.Option("field.tabs", m.Option(aaa.SESS))
	return m
}
func (s shells) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		s.Client.List(m, arg...).Action(s.Create)
	} else {
		s.Create(m.Options(aaa.SESS, arg[0])).Action(s.Create)
	}
}

func init() { ice.CodeModCmd(shells{}) }
