package main

import (
	"github.com/shylinux/icebergs"
	_ "github.com/shylinux/icebergs/base"
	_ "github.com/shylinux/icebergs/core"
	_ "github.com/shylinux/icebergs/misc"
	"github.com/shylinux/toolkits"

	"github.com/shylinux/icebergs/core/wiki"

	_ "github.com/shylinux/linux-story/cli/make"
	_ "github.com/shylinux/linux-story/cli/text"
	_ "github.com/shylinux/redis-story/src/client"
	_ "github.com/shylinux/redis-story/src/data"
	_ "github.com/shylinux/redis-story/src/raft"
	_ "github.com/shylinux/redis-story/src/server"
)

func init() {
	wiki.Index.Register(&ice.Context{Name: "redis", Help: "redis",
		Caches: map[string]*ice.Cache{},
		Configs: map[string]*ice.Config{
			"hi": {Name: "hi", Help: "hi", Value: kit.Data()},
		},
		Commands: map[string]*ice.Command{
			ice.ICE_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				m.Load()
			}},
			ice.ICE_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
				m.Save("hi")
			}},
		},
	}, nil)
}
func main() { println(ice.Run()) }
