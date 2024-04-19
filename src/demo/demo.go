package client

import (
	"github.com/redis/go-redis/v9"
	"shylinux.com/x/ice"
)

type demo struct {
	ice.Hash

	list string `name:"list hash auto" help:"示例"`
}

func (s demo) List(m *ice.Message, arg ...string) {
	// rdb := redis.NewClusterClient(&redis.ClusterOptions{Password: "123", Addrs: []string{":10001", ":10002", ":10003", ":10004", ":10005", ":10006"}})
	rdb := redis.NewClusterClient(&redis.ClusterOptions{Password: "123", Addrs: []string{":10001"}})
	m.Echo("%v", rdb.Get(m, "hi"))
}

func init() { ice.Cmd("web.code.redis.client.demo", demo{}) }
