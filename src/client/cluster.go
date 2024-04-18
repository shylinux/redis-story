package client

import (
	"strings"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	NODEID = "nodeid"
	MASTER = "master"
	SLOTS  = "slots"
)

type cluster struct {
	Client
	create  string `name:"create cluster-replicas*=1"`
	add     string `name:"add from* to"`
	del     string `name:"del from*"`
	reshard string `name:"reshard from* to* slots"`
	slaveOf string `name:"slaveOf to*"`
}

func (s cluster) Create(m *ice.Message, arg ...string) {
	if m.Option(aaa.SESS) == "" {
		m.ToastFailure("please select node")
		return
	}
	list := m.CmdMap(s.Client, aaa.SESS)
	cmd, password, args := "redis-cli", "", []string{}
	kit.For(kit.Split(m.Option(aaa.SESS)), func(p string) {
		cmd, password = s.findCmds(m, list[p][tcp.PORT]), list[p][aaa.PASSWORD]
		args = append(args, tcp.HostPort(list[p][tcp.HOST], list[p][tcp.PORT]))
	})
	if len(args) < 6 {
		m.ToastFailure("please select more than 6 node")
		return
	}
	m.SystemCmd(cmd, "-a", password, "--cluster", mdb.CREATE, args, s.Code.Args(m), "--cluster-yes")
	m.Sleep("3s").Cmd(s.Client, s.Scan)
}
func (s cluster) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		s.Client.List(m, arg...).Table(func(value ice.Maps) {
			button := []ice.Any{}
			if value[aaa.ROLE] == MASTER && value[MASTER] == "" {
				button = append(button, s.SlaveOf)
			} else if value[aaa.ROLE] == MASTER {
				button = append(button, s.Rebalance, s.Reshard, s.Add, s.Del)
			} else {
				button = append(button, s.Del)
			}
			m.PushButton(button...)
		}).Action(s.Create, s.Scan)
	} else {
		kit.For(kit.SplitLine(m.Cmd(s.Client, arg[0], "cluster", "nodes").Append(ice.RES)), func(text string) {
			ls := kit.SplitWord(text)
			ls[2] = strings.Replace(ls[2], "myself,", "", 1)
			m.Push(mdb.ID, ls[0]).Push(tcp.HOSTPORT, ls[1])
			m.Push(aaa.ROLE, ls[2]).Push(MASTER, ls[3]).Push(SLOTS, kit.Select("", ls, 8))
		})
	}
}
func (s cluster) SlaveOf(m *ice.Message, arg ...string) {
	msg := m.Cmd(s.Client, m.Option(nfs.TO))
	s.Client.Cmds(m, "", msg.Append(tcp.HOST), msg.Append(tcp.PORT))
	s.Client.Cmds(m, CONFIG, SET, "masterauth", msg.Append(aaa.PASSWORD))
	m.ProcessRefresh().ToastSuccess()
	m.Cmd(s.Client, s.Client.Scan)
}
func (s cluster) Rebalance(m *ice.Message, arg ...string) {
	s.Cmds(m, "", func(node string, from *ice.Message) []string {
		return []string{node}
	})
}
func (s cluster) Reshard(m *ice.Message, arg ...string) {
	s.Cmds(m, "", func(node string, from *ice.Message) []string {
		to := m.Cmd(s.Client, m.Option(nfs.TO))
		return []string{node, "--cluster-from", from.Append(NODEID), "--cluster-to", to.Append(NODEID), "--cluster-slots", m.Option(SLOTS), "--cluster-yes"}
	})
}
func (s cluster) Add(m *ice.Message, arg ...string) {
	s.Cmds(m, "add-node", func(node string, from *ice.Message) []string {
		args := []string{tcp.HostPort(from.Append(tcp.HOST), from.Append(tcp.PORT)), node}
		kit.If(m.Option(nfs.TO), func(p string) {
			args = append(args, "--cluster-slave", "--cluster-master-id", m.Cmd(s.Client, p).Append(NODEID))
		})
		return args
	})
}
func (s cluster) Del(m *ice.Message, arg ...string) {
	s.Cmds(m, "del-node", func(node string, from *ice.Message) []string {
		return []string{node, from.Append(NODEID)}
	})
}

func init() { ice.CodeModCmd(cluster{}) }

func (s cluster) CmdsXterm(m *ice.Message, cmd string, cb func(string, *ice.Message) []string, arg ...string) {
	m.ProcessXterm(kit.Keys(cmd), func() string { return kit.JoinCmds(s.cmds(m, cmd, cb)...) }, arg...)
}
func (s cluster) Cmds(m *ice.Message, cmd string, cb func(string, *ice.Message) []string) {
	defer m.ToastProcess()()
	m.SystemCmd(s.cmds(m, cmd, cb))
	m.Cmd(s.Client, s.Client.Scan)
}
func (s cluster) cmds(m *ice.Message, cmd string, cb func(string, *ice.Message) []string) []string {
	msg, from := m.Cmd(s.Client, m.Option(aaa.SESS)), m.Cmd(s.Client, m.Option(nfs.FROM))
	return kit.Simple(s.findCmds(m, msg.Append(tcp.PORT)), "-a", msg.Append(aaa.PASSWORD), "--cluster",
		kit.Select(m.ActionKey(), cmd), cb(tcp.HostPort(msg.Append(tcp.HOST), msg.Append(tcp.PORT)), from))
}
