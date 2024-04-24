package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
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
	limit  string `data:"300"`
	vendor string `data:"https://redis.io/docs/latest/commands/get/"`
	field  string `data:"time,command,order,type,name,text"`
}

var types = []string{
	"generic", STRING, HASH, LIST, ZSET, SET,
	"sorted_set", "geo", "hyperloglog",
	"pubsub", "stream", "transactions",
	"server", "cluster", "connection",
	"command", "scripting", "internal",
}

func (s commands) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case COMMAND:
		m.Cmdy("", m.Option(aaa.SESS)).Cut(arg[0], mdb.TYPE)
	case mdb.TYPE:
		m.Push(arg[0], types)
	default:
		s.Hash.Inputs(m, arg...)
	}
}
func (s commands) HelpCmd(m *ice.Message, arg ...string) {
}
func (s commands) Pie(m *ice.Message, arg ...string) {
	list := map[string]int{}
	m.Cmd("", m.Option(aaa.SESS)).Table(func(value ice.Maps) { list[value[mdb.TYPE]]++ })
	kit.For(list, func(k string, n int) { m.Push(mdb.KEY, k).Push(mdb.COUNT, n) })
	m.Display("/plugin/story/pie.js")
	m.SortIntR(mdb.COUNT)
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
				value, button = map[string]string{mdb.TIME: m.Time(), mdb.TYPE: "unknown", COMMAND: command}, append(button, s.HelpCmd, s.Create)
			}, func() {
				button = append(button, s.Remove)
			})
			m.PushRecord(value, kit.Split(m.Config(mdb.FIELD))...).PushButton(button...)
		})
	})
	m.Action(html.FILTER, s.Create, s.Pie, s.Vendor).StatusTimeCountStats(mdb.TYPE).Sort("type,order,command", types, ice.STR, ice.STR)
}

func init() { ice.CodeModCmd(commands{}) }
