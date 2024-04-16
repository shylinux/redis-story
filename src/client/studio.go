package client

import "shylinux.com/x/ice"

type studio struct {
	Client
	online string `data:"true"`
	create string `name:"create sess*=biz host*=localhost port*=10001 password*=root"`
	list   string `name:"list list" icon:"redis.png"`
}

func (s studio) Create(m *ice.Message, arg ...string) { m.Cmdy(s.Client, m.ActionKey(), arg) }
func (s studio) Remove(m *ice.Message, arg ...string) { m.Cmdy(s.Client, m.ActionKey(), arg) }
func (s studio) Scan(m *ice.Message, arg ...string)   { m.Cmdy(s.Client, m.ActionKey(), arg) }
func (s studio) List(m *ice.Message, arg ...string) {
	m.Cmdy(s.Client, arg).PushAction(s.Remove).Display("")
}

func init() { ice.CodeModCmd(studio{}) }
