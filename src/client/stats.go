package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

type stats struct {
	Client
	ice.Zone
	short  string `data:"sess"`
	field  string `data:"time,sess,host,port,keys,used_memory_human"`
	fields string `data:"time,id,keys,used_memory,used_cpu_user,used_cpu_sys,connected_clients"`
	list   string `name:"list sess id auto"`
}

func (s stats) Scan(m *ice.Message, arg ...string) {
	m.Cmd(s.Client).Table(func(value ice.Maps) {
		info := s.cmdInfo(m, value[aaa.SESS])
		get := func(z, k string) []string { return []string{k, kit.Format(kit.Value(info, kit.Keys(z, k)))} }
		keys := kit.Select("", kit.Split(get("Keyspace", "db0")[1], "=,"), 1)
		s.Zone.Insert(m.Spawn(), kit.Simple(aaa.SESS, value[aaa.SESS], KEYS, keys, get("Memory", "used_memory"),
			get("CPU", "used_cpu_user"), get("CPU", "used_cpu_sys"), get("Clients", "connected_clients"),
		)...)
		s.Hash.Modify(m.Spawn(), kit.Simple(aaa.SESS, value[aaa.SESS],
			tcp.HOST, value[tcp.HOST], tcp.PORT, value[tcp.PORT],
			KEYS, keys, get("Memory", "used_memory_human"))...)
	})
}
func (s stats) List(m *ice.Message, arg ...string) {
	if s.Zone.List(m, arg...); len(arg) == 0 {
		m.Action(s.Scan).SortIntR(tcp.PORT)
	} else {
		m.Action()
	}
}

func init() { ice.CodeModCmd(stats{}) }
