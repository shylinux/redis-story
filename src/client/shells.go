package client

import (
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

func (s shells) Create(m *ice.Message, arg ...string) *ice.Message {
	msg := m.Cmd(s.Client, m.Option(aaa.SESS))
	m.ProcessXterm("", kit.Format("%s -h %s -p %s -a %s", s.findCmds(m, msg.Append(tcp.PORT)), msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.PASSWORD)), arg...)
	m.Push(ctx.STYLE, html.OUTPUT).Options("field.tabs", kit.Hashs("uniq")).Option(ice.FIELD_PREFIX, ctx.ACTION, ctx.RUN, code.XTERM)
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
