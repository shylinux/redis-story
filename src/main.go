package main

import (
	"github.com/shylinux/icebergs"
	_ "github.com/shylinux/icebergs/base"
	_ "github.com/shylinux/icebergs/core"
	_ "github.com/shylinux/icebergs/misc"

	_ "github.com/shylinux/icebergs/misc/alpha"

	_ "github.com/shylinux/linux-story/cli/make"
	_ "github.com/shylinux/linux-story/cli/text"
	_ "github.com/shylinux/redis-story/cluster"
	_ "github.com/shylinux/redis-story/server"
)

func main() {
	println(ice.Run())
}
