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
	REDIS_POOL = "redis_pool"
)

type client struct {
	ice.Hash
	short string `data:"name"`
	field string `data:"time,name,host,port"`

	del    string `name:"del" help:"删除"`
	info   string `name:"info" help:"信息"`
	keys   string `name:"keys pattern" help:"列表"`
	create string `name:"create name=biz host=localhost port=10001 password=12345678" help:"连接"`
	prunes string `name:"prunes pattern" help:"清理"`
	list   string `name:"list name@key run info keys prunes create cmd:textarea" help:"缓存值"`
}

func (c client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case mdb.NAME:
		c.List(m).Cut(mdb.NAME)
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
	default:
		c.Hash.Inputs(m, arg...)
	}
}
func (c client) Prunes(m *ice.Message, arg ...string) {
	m.OptionCB(tcp.CLIENT, func(redis *redis) {
		res, err := redis.Do("keys", m.Option("pattern"))
		m.Assert(err)
		for _, k := range kit.Slice(kit.Simple(res), 0, 100) {
			m.Push(mdb.KEY, k)
			res, err := redis.Do("del", k)
			m.Push(ice.RES, kit.Format(res))
			m.Push(ice.ERR, kit.Format(err))
		}
	})
	c.List(m, m.Option(mdb.NAME), "keys")
}
func (c client) Del(m *ice.Message, arg ...string) {
	c.List(m, m.Option(mdb.NAME), "del", m.Option("key"))
}
func (c client) Info(m *ice.Message, arg ...string) {
	m.OptionCB(tcp.CLIENT, func(redis *redis) {
		res, _ := redis.Do("info")
		data, domain := kit.Dict(), ""
		for _, line := range strings.Split(kit.Format(res), "\r\n") {
			if strings.HasPrefix(line, "# ") {
				domain = strings.TrimPrefix(line, "# ")
				continue
			}
			ls := strings.SplitN(strings.TrimSpace(line), ice.DF, 2)
			if len(ls) > 1 {
				kit.Value(data, kit.Keys(domain, ls[0]), ls[1])
			}
		}
		m.PushDetail(data)
		m.StatusTimeCount()
	})
	c.List(m, arg[0], "keys")
}
func (c client) Keys(m *ice.Message, arg ...string) *ice.Message {
	m.OptionCB(tcp.CLIENT, func(redis *redis) {
		res, err := redis.Do("keys", kit.Select("*", m.Option("pattern")))
		m.Assert(err)
		for _, k := range kit.Slice(kit.Simple(res), 0, 100) {
			t := kit.Format(redis.Done("type", k))
			m.Push("type", t)
			m.Push("ttl", kit.Format(redis.Done("ttl", k)))
			m.Push("key", k)
			switch t {
			case "set":
				m.Push(mdb.VALUE, kit.Format(redis.Done("SMEMBERS", k)))
			case "zset":
				list := kit.Simple(redis.Done("ZRANGE", k, "0", "-1", "WITHSCORES"))
				data := kit.Dict()
				for i := 0; i < len(list)-1; i += 2 {
					data[list[i]] = list[i+1]
				}
				m.Push(mdb.VALUE, kit.Format(data))
			case "list":
				m.Push(mdb.VALUE, kit.Format(redis.Done("LRANGE", k, "0", "-1")))
			case "string":
				m.Push(mdb.VALUE, kit.Format(redis.Done("GET", k)))
			case "hash":
				m.Push(mdb.VALUE, kit.Formats(kit.Dict(redis.Done("HGETALL", k))))

			default:
				m.Push(mdb.VALUE, "")
			}
			m.PushAction(c.Del)
		}
		m.StatusTimeCount()
	})
	c.List(m, m.Option(mdb.NAME), "keys")
	return m
}
func (c client) List(m *ice.Message, arg ...string) *ice.Message {
	if len(arg) == 0 || arg[0] == "" { // 连接列表
		defer m.PushAction(c.Hash.Remove)
		defer m.Sort(mdb.NAME)
		c.Hash.List(m)
		return m
	}

	m.Cmd(mdb.SELECT, m.PrefixKey(), "", mdb.HASH, mdb.NAME, arg[0], func(fields []string, value map[string]interface{}) {
		if len(arg) == 1 {
			m.PushDetail(value)
			m.EchoScript(kit.Format("redis-cli -h %s -p %s -a '%s'", value[tcp.HOST], value[tcp.PORT], value[aaa.PASSWORD]))
			return
		}
		// 连接池
		var rp *RedisPool
		switch val := value[REDIS_POOL].(type) {
		case *RedisPool:
			rp = val
		default:
			rp = NewRedisPool(kit.Format("%s:%s", value[tcp.HOST], value[tcp.PORT]), kit.Format(value[aaa.PASSWORD]))
			value[REDIS_POOL] = rp
		}

		r := rp.Get()
		defer rp.Put(r)

		if cb := m.OptionCB(tcp.CLIENT); cb != nil {
			switch cb := cb.(type) {
			case func(*redis):
				cb(r)
				return
			}
		}

		// 命令行
		for _, line := range strings.Split(strings.Join(arg[1:], ice.SP), ice.NL) {
			m.Push(mdb.TIME, kit.Format(time.Now()))
			m.Push(ice.CMD, line)
			cmds := kit.Split(line)
			if res, err := r.Do(cmds[0], cmds[1:]...); err == nil {
				if cb := m.OptionCB(tcp.CLIENT); cb != nil {
					switch cb := cb.(type) {
					case func(interface{}):
						cb(res)
						continue
					}
				}
				m.Push(ice.ERR, "")
				m.Push(ice.RES, kit.Format(res))
				m.Echo("%v", res)
			} else {
				m.Push(ice.ERR, err)
				m.Push(ice.RES, "")
			}
		}
	})
	return m
}

func init() { ice.CodeModCmd(client{}) }
