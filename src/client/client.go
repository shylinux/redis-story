package client

import (
	"path"
	"strings"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/nfs"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	INFO = "INFO"
)
const (
	REDIS  = "redis"
	STRING = "string"
	HASH   = "hash"
	LIST   = "list"
	ZSET   = "zset"
)

type client struct {
	ice.Hash
	export  string `data:"true"`
	short   string `data:"sess"`
	field   string `data:"time,sess,host,port,role,master,save,appendonly,maxclients,maxmemory"`
	create  string `name:"create sess*=biz host*=localhost port*=10001 password*=root"`
	list    string `name:"list sess auto"`
	prunes  string `name:"prunes limit*=100 pattern"`
	keys    string `name:"keys limit*=100 pattern"`
	slaveOf string `name:"slaveOf to*" icon:"bi bi-copy"`
}

func (s client) Inputs(m *ice.Message, arg ...string) {
	kit.If(arg[0] == nfs.TO, func() { arg[0] = aaa.SESS })
	if arg[0] == aaa.SESS {
		m.Cmdy("").Cut(aaa.SESS, tcp.HOST, tcp.PORT)
		return
	}
	switch s.Hash.Inputs(m, arg...); arg[0] {
	case tcp.PORT:
		m.Push(arg[0], "6379")
	}
}
func (s client) Scan(m *ice.Message, arg ...string) {
	list := map[string]string{}
	m.Cmd("").Table(func(value ice.Maps) {
		list[tcp.HostPort(value[tcp.HOST], value[tcp.PORT])] = value[aaa.SESS]
	}).GoToastTable(aaa.SESS, func(value ice.Maps) {
		msg := m.Cmd(configs{}, value[aaa.SESS])
		data := s.cmdInfo(m, value[aaa.SESS])
		get := func(key string) string { return kit.Format(kit.Value(data, "Replication."+key)) }
		s.Hash.Modify(m, kit.Simple(aaa.SESS, value[aaa.SESS],
			aaa.ROLE, get(aaa.ROLE), "master", list[tcp.HostPort(get("master_host"), get("master_port"))],
			msg.AppendSimple("save,appendonly,maxclients,maxmemory"),
		)...)
	})
}
func (s client) List(m *ice.Message, arg ...string) {
	if s.Hash.List(m, arg...); len(arg) < 1 || arg[0] == "" {
		m.Table(func(value ice.Maps) {
			switch value["role"] {
			case "master", "":
				m.PushButton(s.SlaveOf, s.Save, s.Info, s.Xterm, s.Remove)
			default:
				m.PushButton(s.Save, s.Info, s.Xterm, s.Remove)
			}
		}).Action(s.Create, s.Scan).Sort(aaa.SESS)
	} else if len(arg) < 2 || arg[1] == "" {
		m.PushAction(s.Xterm, s.Remove).Action(s.Info)
		m.EchoScript(kit.Format("redis-cli -h %s -p %s -a '%s'", m.Append(tcp.HOST), m.Append(tcp.PORT), m.Append(aaa.PASSWORD)))
	} else {
		m.SetAppend()
		s.cmds(m, arg...)
	}
}
func (s client) SlaveOf(m *ice.Message, arg ...string) {
	msg := m.Cmd("", m.Option(nfs.TO))
	s.cmds(m, m.Option(aaa.SESS), "slaveof", msg.Append(tcp.HOST), msg.Append(tcp.PORT))
	s.cmds(m, m.Option(aaa.SESS), CONFIG, SET, "masterauth", msg.Append(aaa.PASSWORD))
	m.ProcessRefresh().ToastSuccess()
}
func (s client) Xterm(m *ice.Message, arg ...string) {
	msg := s.Hash.List(m.Spawn(), m.Option(aaa.SESS))
	m.ProcessXterm(kit.Keys(REDIS, msg.Append(aaa.SESS)), kit.Format("%s -h %s -p %s -a %s",
		path.Join(ice.USR_LOCAL_DAEMON, msg.Append(tcp.PORT), "bin/redis-cli"),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.PASSWORD)), arg...)
}
func (s client) Save(m *ice.Message, arg ...string) { s.Cmds(m, "") }
func (s client) Info(m *ice.Message, arg ...string) {
	m.Echo(kit.Format(s.cmdInfo(m, kit.Select(m.Option(aaa.SESS), arg, 0)))).DisplayStoryJSON()
}

func init() { ice.CodeModCmd(client{}) }

func (s client) cmdInfo(m *ice.Message, sess string) ice.Any {
	data, domain := kit.Dict(), ""
	for _, line := range strings.Split(s.cmds(m.Spawn(), sess, INFO).Append(ice.RES), "\r\n") {
		if strings.HasPrefix(line, "# ") {
			domain = strings.TrimPrefix(line, "# ")
			continue
		}
		if ls := strings.SplitN(strings.TrimSpace(line), ice.DF, 2); len(ls) > 1 {
			kit.Value(data, kit.Keys(domain, ls[0]), ls[1])
		}
	}
	return data
}
func (s client) cmds(m *ice.Message, arg ...string) *ice.Message {
	msg := s.Hash.List(m.Spawn(), arg[0])
	rp := s.Hash.Target(m, arg[0], func() ice.Any {
		return NewRedisPool(tcp.HostPort(msg.Append(tcp.HOST), msg.Append(tcp.PORT)), msg.Append(aaa.PASSWORD))
	}).(*RedisPool)
	r := rp.Get()
	defer rp.Put(r)
	switch cb := m.OptionCB("").(type) {
	case func(*redis):
		cb(r)
		return m
	case func(ice.Any):
		if res, err := r.Do(arg[1], arg[2:]...); !m.Warn(err) {
			cb(res)
		}
		return m
	}
	for _, line := range strings.Split(strings.TrimSpace(strings.Join(arg[1:], ice.SP)), ice.NL) {
		m.Push(mdb.TIME, kit.Format(time.Now())).Push(ice.CMD, line)
		cmds := kit.Split(line)
		if res, err := r.Do(cmds[0], cmds[1:]...); err == nil {
			m.Push(ice.ERR, "").Push(ice.RES, kit.Format(res))
		} else {
			m.Push(ice.ERR, kit.Format(err)).Push(ice.RES, "")
		}
	}
	kit.If(m.Append(ice.ERR), func(p string) { m.ToastFailure(p) })
	return m
}
func (s client) Cmds(m *ice.Message, cmd string, arg ...string) *ice.Message {
	return m.Cmd(s, m.Option(aaa.SESS), kit.Select(m.ActionKey(), cmd), arg)
}

type Client struct {
	client
	list string `name:"list sess auto"`
}

func (s Client) List(m *ice.Message, arg ...string) *ice.Message {
	if len(arg) == 0 {
		m.Cmdy(s.client).Cut("time,sess,host,port,role").PushAction().Action()
	}
	return m
}
func init() { ice.CodeModCmd(Client{}) }
