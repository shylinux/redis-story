package main

import (
	"shylinux.com/x/ice"
	_ "shylinux.com/x/redis-story/src/client"
	_ "shylinux.com/x/redis-story/src/server"
)

func main() { print(ice.Run()) }
