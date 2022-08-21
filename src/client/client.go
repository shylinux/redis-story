package client

import (
	"strings"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	KEYS = "KEYS"
)

type client struct {
	ice.Code
	ice.Hash
	short string `data:"name"`
	field string `data:"time,name,host,port"`

	del    string `name:"del" help:"删除"`
	info   string `name:"info" help:"信息"`
	keys   string `name:"keys limit=100 pattern" help:"列表"`
	prunes string `name:"prunes limit=100 pattern" help:"清理"`
	create string `name:"create name=biz host=localhost port=10001 password=root" help:"连接"`
	list   string `name:"list name@key run info keys prunes create cmd:textarea" help:"缓存"`
}

func (s client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case mdb.NAME:
		s.List(m).Cut(mdb.NAME)
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
	default:
		s.Hash.Inputs(m, arg...)
	}
}
func (s client) Del(m *ice.Message, arg ...string) {
	s.List(m, m.Option(mdb.NAME), "del", m.Option(mdb.KEY))
}
func (s client) Info(m *ice.Message, arg ...string) {
	s.List(m, arg[0], "info")
	data, domain := kit.Dict(), ""
	for _, line := range strings.Split(m.Result(), "\r\n") {
		if strings.HasPrefix(line, "# ") {
			domain = strings.TrimPrefix(line, "# ")
			continue
		}
		if ls := strings.SplitN(strings.TrimSpace(line), ice.DF, 2); len(ls) > 1 {
			kit.Value(data, kit.Keys(domain, ls[0]), ls[1])
		}
	}
	m.SetAppend().SetResult().PushDetail(data).StatusTimeCount()
}
func (s client) Keys(m *ice.Message, arg ...string) {
	m.OptionCB("", func(redis *redis) {
		res, err := redis.Do(KEYS, kit.Select("*", m.Option("pattern")))
		m.Assert(err)
		m.OptionFields("")
		for _, k := range kit.Slice(kit.Simple(res), 0, kit.Int(kit.Select("100", m.Option(mdb.LIMIT)))) {
			t := kit.Format(redis.Done("type", k))
			m.Push("type", t)
			m.Push("ttl", kit.Format(redis.Done("ttl", k)))
			m.Push("key", k)
			switch t {
			case "string":
				m.Push(mdb.VALUE, kit.Format(redis.Done("GET", k)))
			case "hash":
				m.Push(mdb.VALUE, kit.Formats(kit.Dict(redis.Done("HGETALL", k))))
			case "list":
				m.Push(mdb.VALUE, kit.Format(redis.Done("LRANGE", k, "0", "-1")))
			case "set":
				m.Push(mdb.VALUE, kit.Format(redis.Done("SMEMBERS", k)))
			case "zset":
				data, list := kit.Dict(), kit.Simple(redis.Done("ZRANGE", k, "0", "-1", "WITHSCORES"))
				for i := 0; i < len(list)-1; i += 2 {
					data[list[i]] = list[i+1]
				}
				m.Push(mdb.VALUE, kit.Format(data))
			default:
				m.Push(mdb.VALUE, "")
			}
			m.PushAction(s.Del)
		}
		m.Sort("type,key").StatusTimeCount()
	})
	s.List(m, m.Option(mdb.NAME), KEYS)
}
func (s client) Prunes(m *ice.Message, arg ...string) {
	m.OptionCB("", func(redis *redis) {
		res, err := redis.Do(KEYS, m.Option("pattern"))
		m.Assert(err)
		for _, k := range kit.Slice(kit.Simple(res), 0, 100) {
			m.Push(mdb.KEY, k)
			res, err := redis.Do("del", k)
			m.Push(ice.ERR, kit.Format(err))
			m.Push(ice.RES, kit.Format(res))
		}
	})
	s.List(m, m.Option(mdb.NAME), KEYS)
}
func (s client) List(m *ice.Message, arg ...string) *ice.Message {
	if s.Hash.List(m, arg...); len(arg) < 1 || arg[0] == "" {
		m.Sort(mdb.NAME).PushAction(s.Xterm, s.Hash.Remove)
		return m // 连接列表
	} else if len(arg) < 2 || arg[1] == "" {
		m.PushAction(s.Xterm, s.Hash.Remove)
		m.EchoScript(kit.Format("redis-cli -h %s -p %s -a '%s'", m.Append(tcp.HOST), m.Append(tcp.PORT), m.Append(aaa.PASSWORD)))
		return m // 连接详情
	}

	// 连接池
	rp := s.Hash.Target(m, arg[0], func() ice.Any {
		return NewRedisPool(kit.Format("%s:%s", m.Append(tcp.HOST), m.Append(tcp.PORT)), m.Append(aaa.PASSWORD))
	}).(*RedisPool)

	r := rp.Get()
	defer rp.Put(r)

	m.SetAppend()
	switch cb := m.OptionCB("").(type) {
	case func(*redis):
		cb(r)
		return m
	}

	// 命令行
	for _, line := range strings.Split(strings.Join(arg[1:], ice.SP), ice.NL) {
		m.Push(mdb.TIME, kit.Format(time.Now()))
		m.Push(ice.CMD, line)
		cmds := kit.Split(line)
		if res, err := r.Do(cmds[0], cmds[1:]...); err == nil {
			switch cb := m.OptionCB("").(type) {
			case func(ice.Any):
				cb(res)
				continue
			}
			m.Push(ice.ERR, "")
			m.Push(ice.RES, kit.Format(res))
			m.Echo("%v", res)
		} else {
			m.Push(ice.ERR, kit.Format(err))
			m.Push(ice.RES, "")
		}
	}
	if m.Length() == 1 && m.Append(ice.ERR) == "" {
		m.SetAppend()
	}
	return m
}
func (s client) Xterm(m *ice.Message, arg ...string) {
	m.OptionFields("host,port,password")
	msg := s.List(m.Spawn(), m.Option(mdb.NAME))
	s.Code.Xterm(m, kit.Format("%s -h %s -p %s -a %s", kit.Path(ice.USR_LOCAL_DAEMON, msg.Append(tcp.PORT), "bin/redis-cli"),
		msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(aaa.PASSWORD)), arg...)
}

func init() { ice.CodeModCmd(client{}) }
