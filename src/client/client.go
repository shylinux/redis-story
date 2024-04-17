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
	REDIS = "redis"
	INFO  = "info"
	SAVE  = "save"
)

type client struct {
	ice.Hash
	ice.Code
	checkbox      string `data:"true"`
	export        string `data:"true"`
	short         string `data:"sess"`
	field         string `data:"time,sess,host,port,password,role,master,nodeid,slots,save,appendonly,maxclients,maxmemory"`
	create        string `name:"create sess*=biz host*=localhost port*=10001 password*=root"`
	list          string `name:"list sess auto"`
	slaveOf       string `name:"slaveOf to*"`
	clusterCreate string `name:"clusterCreate cluster-replicas*=1"`
}

func (s client) Inputs(m *ice.Message, arg ...string) {
	kit.If(kit.IsIn(arg[0], nfs.FROM, nfs.TO), func() { arg[0] = aaa.SESS })
	if arg[0] == aaa.SESS {
		m.Cmdy("").Cut(aaa.SESS, tcp.HOST, tcp.PORT, aaa.ROLE)
	} else if s.Hash.Inputs(m, arg...); arg[0] == tcp.PORT {
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
			aaa.ROLE, get(aaa.ROLE), MASTER, list[tcp.HostPort(get("master_host"), get("master_port"))],
			msg.AppendSimple("save,appendonly,maxclients,maxmemory"),
		)...)
		if kit.Format(kit.Value(data, "Cluster.cluster_enabled")) == "1" {
			m.Cmd(cluster{}, value[aaa.SESS]).Table(func(value ice.Maps) {
				ls := kit.Split(value[tcp.HOSTPORT], ":@")
				s.Hash.Modify(m, kit.Simple(aaa.SESS, list[tcp.HostPort(ls[0], ls[1])],
					aaa.ROLE, value[aaa.ROLE], MASTER, value[MASTER], NODEID, value[mdb.ID], SLOTS, value[SLOTS],
				)...)
			})
		}
	})
}
func (s client) List(m *ice.Message, arg ...string) {
	if s.Hash.List(m, arg...); len(arg) == 0 || arg[0] == "" {
		m.Table(func(value ice.Maps) {
			button := []ice.Any{}
			kit.If(value[aaa.ROLE] == MASTER && value[MASTER] == "", func() { button = append(button, s.SlaveOf) })
			button = append(button, s.Save, s.Info, s.Xterm, s.Remove)
			m.PushButton(button...)
		}).Action(s.Create, s.ClusterCreate, s.Scan).Sort(tcp.PORT, ice.INT_R)
	} else if len(arg) == 1 || arg[1] == "" {
		m.PushAction(s.Info, s.Xterm, s.Remove).Action()
	} else {
		m.SetAppend()
		s.cmds(m, arg...)
	}
}
func (s client) ClusterCreate(m *ice.Message, arg ...string) {
	list := m.CmdMap(s, aaa.SESS)
	cmd, password, args := "redis-cli", "", []string{}
	kit.For(kit.Split(m.Option(aaa.SESS)), func(p string) {
		cmd, password = s.findCmds(m, list[p][tcp.PORT]), list[p][aaa.PASSWORD]
		args = append(args, tcp.HostPort(list[p][tcp.HOST], list[p][tcp.PORT]))
	})
	m.SystemCmd(cmd, "-a", password, "--cluster", mdb.CREATE, args, s.Code.Args(m), "--cluster-yes")
}
func (s client) SlaveOf(m *ice.Message, arg ...string) {
	msg := m.Cmd("", m.Option(nfs.TO))
	s.cmds(m, m.Option(aaa.SESS), "slaveof", msg.Append(tcp.HOST), msg.Append(tcp.PORT))
	s.cmds(m, m.Option(aaa.SESS), CONFIG, SET, "masterauth", msg.Append(aaa.PASSWORD))
	m.ProcessRefresh().ToastSuccess()
}
func (s client) Xterm(m *ice.Message, arg ...string) {
	m.ProcessXterm(kit.Format("%s(%s:%s)", m.Option(aaa.SESS), m.Option(tcp.HOST), m.Option(tcp.PORT)), s.findCmdArgs(m, m.Option(aaa.SESS)), arg...)
}
func (s client) Save(m *ice.Message, arg ...string) {
	s.Cmds(m, "")
}
func (s client) Info(m *ice.Message, arg ...string) {
	m.Echo(kit.Format(s.cmdInfo(m, m.Option(aaa.SESS)))).DisplayStoryJSON()
}

func init() { ice.CodeModCmd(client{}) }

func (s client) findCmdArgs(m *ice.Message, sess string) string {
	msg := m.Cmd(s, sess)
	return kit.Format("%s -h %s -p %s -a %s", s.findCmds(m, msg.Append(tcp.PORT)), msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.PASSWORD))
}
func (s client) findCmds(m *ice.Message, port string) string {
	cmd := "redis-cli"
	nfs.Exists(m.Message, path.Join(ice.USR_LOCAL_DAEMON, port, "bin/redis-cli"), func(p string) { cmd = p })
	return cmd
}
func (s client) cmdInfo(m *ice.Message, sess string) ice.Any {
	data, domain := kit.Dict(), ""
	s.Cmds(m.Options(aaa.SESS, sess), INFO, func(res ice.Any) {
		kit.For(kit.SplitLine(kit.Format(res), "\r\n"), func(line string) {
			if strings.HasPrefix(line, "# ") {
				domain = strings.TrimPrefix(line, "# ")
			} else if ls := strings.SplitN(strings.TrimSpace(line), ice.DF, 2); len(ls) > 1 {
				kit.Value(data, kit.Keys(domain, ls[0]), ls[1])
			}
		})
	})
	return data
}
func (s client) cmds(m *ice.Message, arg ...string) *ice.Message {
	msg := s.Hash.List(m.Spawn(), arg[0])
	rp := s.Hash.Target(m, arg[0], func() ice.Any {
		return NewRedisPool(tcp.HostPort(msg.Append(tcp.HOST), msg.Append(tcp.PORT)), msg.Append(aaa.PASSWORD))
	}).(*RedisPool)
	r := rp.Get()
	if r == nil {
		return m
	}
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
func (s client) Cmds(m *ice.Message, cmd string, arg ...ice.Any) *ice.Message {
	return m.Cmd(s, m.Option(aaa.SESS), kit.Select(m.ActionKey(), cmd), arg)
}

type Client struct {
	client
	list string `name:"list sess auto"`
}

func (s Client) List(m *ice.Message, arg ...string) *ice.Message {
	if len(arg) == 0 {
		m.Cmdy(s.client).Cut("time,sess,host,port,role,master,nodeid,slots").PushAction().Action()
	}
	return m
}
func init() { ice.CodeModCmd(Client{}) }
