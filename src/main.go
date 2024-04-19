package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/ice/devops"

	_ "shylinux.com/x/redis-story/src/client"
	_ "shylinux.com/x/redis-story/src/demo"
	_ "shylinux.com/x/redis-story/src/server"
)

func main() { print(ice.Run()) }

func init() { ice.Info.NodeIcon = "src/server/redis.png" }
