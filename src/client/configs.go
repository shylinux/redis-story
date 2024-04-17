package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/web/html"
	kit "shylinux.com/x/toolkits"
)

const (
	CONFIG = "config"
)

type configs struct{ Client }

func (s configs) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case nfs.SAVE:
		m.Push(arg[0], "900 1", "300 10", "60 10000")
	case "loglevel":
		m.Push(arg[0], "debug", "verbose", "notice", "warning")
	}
}
func (s configs) Modify(m *ice.Message, arg ...string) {
	s.Cmds(m.Spawn(), CONFIG, SET, arg[0], arg[1])
}
func (s configs) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		s.Client.List(m, arg...)
		return
	}
	m.FieldsSetDetail()
	m.Cmdy(s.Client, arg[0], CONFIG, GET, kit.Select("*", arg, 1), func(res ice.Any) {
		kit.For(res, func(k string, v ice.Any) { m.Push(k, v) })
	}).Action(html.FILTER).Sort(mdb.KEY)
}

func init() { ice.CodeModCmd(configs{}) }
