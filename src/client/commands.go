package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/web/html"
	kit "shylinux.com/x/toolkits"
)

const (
	COMMAND = "command"
)

type commands struct {
	Client
	export string `data:"true"`
	short  string `data:"command"`
	field  string `data:"time,command,type,name,text"`
}

var types = []string{KEYS, STRING, HASH, LIST, ZSET, SET, "admin"}

func (s commands) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case mdb.TYPE:
		m.Push(arg[0], types)
	default:
		s.Hash.Inputs(m, arg...)
	}
}
func (s commands) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		s.Client.List(m, arg...)
		return
	}
	list := map[string]ice.Maps{}
	s.Hash.List(m.Spawn()).Table(func(value ice.Maps) { list[value[COMMAND]] = value })
	m.Cmdy(s.Client, arg[0], COMMAND, func(res ice.Any) {
		kit.For(res, func(value ice.Any) {
			command := kit.Format(kit.Value(value, "0"))
			value, ok := list[command]
			button := []ice.Any{}
			kit.If(!ok, func() {
				value, button = map[string]string{mdb.TIME: m.Time(), COMMAND: command}, append(button, s.Create)
			})
			m.PushRecord(value, kit.Split(m.Config(mdb.FIELD))...).PushButton(button...)
		})
	}).Action(html.FILTER).StatusTimeCountStats(mdb.TYPE).Sort("type,command", types, ice.STR)
}

func init() { ice.CodeModCmd(commands{}) }
