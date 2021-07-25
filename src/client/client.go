package client

import (
	"strings"
	"time"

	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/icebergs/base/mdb"
	"github.com/shylinux/icebergs/base/tcp"
	"github.com/shylinux/redis-story/src/server"
	kit "github.com/shylinux/toolkits"
)

const (
	REDIS_POOL = "redis_pool"
)
const CLIENT = "client"

var Index = &ice.Context{Name: CLIENT, Help: "客户端",
	Configs: map[string]*ice.Config{
		CLIENT: {Name: CLIENT, Help: "客户端", Value: kit.Data()},
	},
	Commands: map[string]*ice.Command{
		ice.CTX_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Watch(server.REDIS_SERVER_START, m.Prefix(CLIENT))
		}},
		ice.CTX_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
		}},

		CLIENT: {Name: "client hash 执行:button 返回 create cmd:textarea", Help: "客户端", Action: map[string]*ice.Action{
			server.REDIS_SERVER_START: {Name: "redis_server_start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},
			mdb.CREATE: {Name: "create host=localhost@key port=10000@key", Help: "连接", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.INSERT, m.Prefix(CLIENT), "", mdb.HASH, arg)
			}},
			mdb.MODIFY: {Name: "modify", Help: "编辑", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.MODIFY, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH), arg)
			}},
			mdb.REMOVE: {Name: "remove", Help: "删除", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(mdb.DELETE, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, m.Option(kit.MDB_HASH))
			}},
			mdb.INPUTS: {Name: "inputs", Help: "补全", Hand: func(m *ice.Message, arg ...string) {
				switch arg[0] {
				case tcp.PORT:
					m.Cmdy(server.SERVER).Appendv(ice.MSG_APPEND, kit.Split("port,status,pid,time"))
				case mdb.HASH:
					m.Option(mdb.FIELDS, "hash,time,host,port")
					m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
				}
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			if len(arg) == 0 || arg[0] == "" { // 连接列表
				m.Fields(len(kit.Slice(arg, 0, 1)), "time,hash,host,port")
				m.Cmdy(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH)
				m.PushAction(mdb.REMOVE)
				return
			}

			m.Cmd(mdb.SELECT, m.Prefix(CLIENT), "", mdb.HASH, kit.MDB_HASH, arg[0], func(fields []string, value map[string]interface{}) {
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
				for _, line := range kit.Split(strings.Join(arg[1:], " "), "\n", "\n", "\n") {
					m.Push(kit.MDB_TIME, kit.Format(time.Now()))
					m.Push(cli.CMD, line)
					if res, err := redis.Do(kit.Split(line)...); err == nil {
						m.Push(cli.ERR, "")
						m.Push(cli.RES, res)
						m.Echo("%v", res)
					} else {
						m.Push(cli.ERR, err)
						m.Push(cli.RES, "")
					}
				}
			})
		}},
	},
}

func init() { server.Index.Merge(Index) }
