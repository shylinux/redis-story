package client

import (
	"strings"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	REDIS_POOL = "redis_pool"
)

type client struct {
	ice.Hash

	short string `data:""`
	field string `data:"time,hash,host,port"`

	create string `name:"create host=localhost port=10001" help:"连接"`
	list   string `name:"list hash run:button back create cmd:textarea" help:"客户端"`
}

func (c client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
	default:
		c.Hash.Inputs(m, arg...)
	}
}
func (c client) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || arg[0] == "" { // 连接列表
		defer m.PushAction(c.Hash.Remove)
		c.Hash.List(m)
		return
	}

	m.Cmd(mdb.SELECT, m.PrefixKey(), "", mdb.HASH, mdb.HASH, arg[0], func(fields []string, value map[string]interface{}) {
		// 连接池
		var rp *RedisPool
		switch val := value[REDIS_POOL].(type) {
		case *RedisPool:
			rp = val
		default:
			rp = NewRedisPool(kit.Format("%s:%s", value[tcp.HOST], value[tcp.PORT]))
			value[REDIS_POOL] = rp
		}

		redis := rp.Get()
		defer rp.Put(redis)

		// 命令行
		for _, line := range strings.Split(strings.Join(arg[1:], ice.SP), ice.NL) {
			m.Push(mdb.TIME, kit.Format(time.Now()))
			m.Push(ice.CMD, line)
			cmds := kit.Split(line)
			if res, err := redis.Do(cmds[0], cmds[1:]...); err == nil {
				m.Push(ice.ERR, "")
				m.Push(ice.RES, res)
				m.Echo("%v", res)
			} else {
				m.Push(ice.ERR, err)
				m.Push(ice.RES, "")
			}
		}
	})
}

func init() { ice.CodeModCmd(client{}) }
