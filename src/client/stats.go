package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	kit "shylinux.com/x/toolkits"
)

type stats struct {
	Client
	ice.Zone
	short  string `data:"sess"`
	field  string `data:"time,sess,count,keys,used_memory_human"`
	fields string `data:"time,id,connected_clients,keys,used_memory,used_cpu_user,used_cpu_sys"`
	list   string `name:"list sess id auto"`
}

func (s stats) Scan(m *ice.Message, arg ...string) {
	m.Cmd(s.Client).Table(func(value ice.Maps) {
		info := s.cmdInfo(m, value[aaa.SESS])
		get := func(z, k string) []string { return []string{k, kit.Format(kit.Value(info, kit.Keys(z, k)))} }
		keys := kit.Select("", kit.Split(get("Keyspace", "db0")[1], "=,"), 1)
		s.Zone.Insert(m.Spawn(), kit.Simple(aaa.SESS, value[aaa.SESS],
			get("Clients", "connected_clients"),
			get("Memory", "used_memory"), KEYS, keys,
			get("CPU", "used_cpu_user"),
			get("CPU", "used_cpu_sys"),
		)...)
		s.Hash.Modify(m.Spawn(), kit.Simple(aaa.SESS, value[aaa.SESS],
			get("Memory", "used_memory_human"), KEYS, keys,
		)...)
	})
}
func (s stats) List(m *ice.Message, arg ...string) {
	s.Zone.List(m, arg...).Action(s.Scan)
}

func init() { ice.CodeModCmd(stats{}) }
