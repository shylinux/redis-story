package server

import (
	ice "github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/core/code"
	kit "github.com/shylinux/toolkits"

	"path"
	"runtime"
)

const (
	REDIS  = "redis"
	SERVER = "server"
)

var Index = &ice.Context{Name: REDIS, Help: "redis",
	Configs: map[string]*ice.Config{
		SERVER: {Name: SERVER, Help: "服务器", Value: kit.Data(
			"windows", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
			"darwin", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
			"linux", "http://download.redis.io/releases/redis-5.0.4.tar.gz",
		)},
	},
	Commands: map[string]*ice.Command{
		SERVER: {Name: "server port=auto path=auto auto 启动:button 构建:button 下载:button", Help: "服务器", Action: map[string]*ice.Action{
			"download": {Name: "download", Help: "下载", Hand: func(m *ice.Message, arg ...string) {
				m.Cmdy(code.INSTALL, "download", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"build": {Name: "build", Help: "构建", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv("prepare", func(p string) {})
				m.Cmdy(code.INSTALL, "build", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)))
			}},
			"start": {Name: "start", Help: "启动", Hand: func(m *ice.Message, arg ...string) {
				m.Optionv("prepare", func(p string) []string { return []string{"--port", path.Base(p)} })
				m.Cmdy(code.INSTALL, "start", m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS)), "bin/redis-server")
			}},
		}, Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			m.Cmdy(code.INSTALL, path.Base(m.Conf(SERVER, kit.Keys(kit.MDB_META, runtime.GOOS))), arg)
		}},
	},
}

func init() { code.Index.Register(Index, nil) }
