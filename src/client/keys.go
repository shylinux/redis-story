package client

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/lex"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/web/html"
	kit "shylinux.com/x/toolkits"
)

const (
	KEYS   = "keys"
	STRING = "string"
	HASH   = "hash"
	LIST   = "list"
	ZSET   = "zset"
	SET    = "set"
	GET    = "get"
	DEL    = "del"
	TTL    = "TTL"
)

type keys struct {
	Client
	prunes string `name:"prunes pattern"`
	set    string `name:"set value"`
	expire string `name:"expire TTL*"`
	rename string `name:"rename newname"`
	hget   string `name:"hget field"`
	hset   string `name:"hset field value"`
	hdel   string `name:"hdel field"`
	lpush  string `name:"lpush value"`
	sadd   string `name:"sadd member"`
	srem   string `name:"srem member"`
	zadd   string `name:"zadd score member"`
	zrem   string `name:"zrem member"`
}

func (s keys) Modify(m *ice.Message, arg ...string) {
	switch arg[0] {
	case mdb.KEY:
		m.Cmd("", s.Rename, "newname", arg[1])
	case mdb.VALUE:
		switch m.Option(mdb.TYPE) {
		case STRING:
			m.Cmd("", s.Set, mdb.VALUE, arg[1])
		}
	case TTL:
		m.Cmd("", s.Expire, arg)
	}
}
func (s keys) Prunes(m *ice.Message, arg ...string) {
	m.Cmdy(s.Client, m.Option(aaa.SESS), KEYS, m.Option(lex.PATTERN), func(redis *redis) {
		res, err := redis.Do(KEYS, m.Option(lex.PATTERN))
		if m.Warn(err) {
			return
		}
		for _, k := range kit.Slice(kit.Simple(res), 0, 100) {
			res, err := redis.Do(DEL, k)
			m.Push(mdb.KEY, k).Push(ice.ERR, kit.Format(err)).Push(ice.RES, kit.Format(res))
		}

	})
}
func (s keys) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		s.Client.List(m, arg...)
		return
	}
	m.Cmdy(s.Client, arg[0], KEYS, kit.Select("*", arg, 1), func(redis *redis) {
		res, err := redis.Do(KEYS, kit.Select(mdb.FOREACH, kit.Select("*", arg, 1)))
		if m.Warn(err) {
			return
		}
		for _, k := range kit.Slice(kit.Simple(res), 0, kit.Int(kit.Select("100", m.Option(mdb.LIMIT)))) {
			t := kit.Format(redis.Done(mdb.TYPE, k))
			ttl := kit.Format(redis.Done(TTL, k))
			button := []ice.Any{}
			var value ice.Any
			switch t {
			case STRING:
				value, button = redis.Done(GET, k), append(button, s.Get, s.Set)
			case HASH:
				value, button = kit.Formats(kit.Dict(redis.Done("HGETALL", k))), append(button, s.Hget, s.Hset, s.Hdel)
			case LIST:
				value, button = redis.Done("LRANGE", k, "0", "-1"), append(button, s.Lpush, s.Lpop)
			case ZSET:
				data, list := kit.Dict(), kit.Simple(redis.Done("ZRANGE", k, "0", "-1", "WITHSCORES"))
				kit.For(list, func(k, v string) { data[k] = v })
				value, button = data, append(button, s.Zadd, s.Zrem)
			case SET:
				value, button = redis.Done("SMEMBERS", k), append(button, s.Sadd, s.Srem)
			}
			button = append(button, s.Expire, s.Rename, s.Del)
			m.Push(mdb.TYPE, t).Push(TTL, ttl).Push(mdb.KEY, k)
			m.Push(mdb.VALUE, kit.Format(value)).PushButton(button...)
		}
	}).Action(html.FILTER, s.Prunes).StatusTimeCount().Sort(mdb.KEY)
}
func (s keys) Get(m *ice.Message, arg ...string)    { m.Echo(s.Cmds(m).Append(ice.RES)) }
func (s keys) Set(m *ice.Message, arg ...string)    { s.Cmds(m, m.Option(mdb.VALUE)) }
func (s keys) Del(m *ice.Message, arg ...string)    { s.Cmds(m) }
func (s keys) Expire(m *ice.Message, arg ...string) { s.Cmds(m, m.Option(TTL)) }
func (s keys) Rename(m *ice.Message, arg ...string) { s.Cmds(m, m.Option("newname")) }
func (s keys) Exists(m *ice.Message, arg ...string) { m.Echo(s.Cmds(m).Append(ice.RES)) }
func (s keys) Hget(m *ice.Message, arg ...string) {
	m.Echo(s.Cmds(m, m.Option(mdb.FIELD)).Append(ice.RES))
}
func (s keys) Hset(m *ice.Message, arg ...string) {
	s.Cmds(m, m.Option(mdb.FIELD), m.Option(mdb.VALUE))
}
func (s keys) Hdel(m *ice.Message, arg ...string)  { s.Cmds(m, m.Option(mdb.FIELD)) }
func (s keys) Lpush(m *ice.Message, arg ...string) { s.Cmds(m, m.Option(mdb.VALUE)) }
func (s keys) Lpop(m *ice.Message, arg ...string)  { s.Cmds(m) }
func (s keys) Sadd(m *ice.Message, arg ...string)  { s.Cmds(m, m.Option(mdb.MEMBER)) }
func (s keys) Srem(m *ice.Message, arg ...string)  { s.Cmds(m, m.Option(mdb.MEMBER)) }
func (s keys) Zadd(m *ice.Message, arg ...string) {
	s.Cmds(m, m.Option(mdb.SCORE), m.Option(mdb.MEMBER))
}
func (s keys) Zrem(m *ice.Message, arg ...string) { s.Cmds(m, m.Option(mdb.MEMBER)) }

func init() { ice.CodeModCmd(keys{}) }

func (s keys) Cmds(m *ice.Message, arg ...string) *ice.Message {
	return s.Client.Cmds(m, "", m.Option(mdb.KEY), arg)
}
