package client

import (
	"strings"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/cli"
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

	create string `name:"create host=localhost@key port=10001@key" help:"连接"`
	list   string `name:"list hash run:button back create cmd:textarea" help:"客户端"`
}

func (c client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Append("append", "port", "status", "time")
	case mdb.HASH:
		c.Hash.List(m)
	}
}
func (c client) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 || arg[0] == "" { // 连接列表
		defer m.PushAction(c.Hash.Remove)
		c.Hash.List(m)
		return
	}

	m.Cmd(mdb.SELECT, m.PrefixKey(), "", mdb.HASH, kit.MDB_HASH, arg[0], func(fields []string, value map[string]interface{}) {
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
		for _, line := range strings.Split(strings.Join(arg[1:], " "), "\n") {
			m.Push(kit.MDB_TIME, kit.Format(time.Now()))
			m.Push(cli.CMD, line)
			cmds := kit.Split(line)
			if res, err := redis.Do(cmds[0], cmds[1:]...); err == nil {
				m.Push(cli.ERR, "")
				m.Push(cli.RES, res)
				m.Echo("%v", res)
			} else {
				m.Push(cli.ERR, err)
				m.Push(cli.RES, "")
			}
		}
	})
}

func init() { ice.Cmd("web.code.redis.client", client{}) }
