package cluster

import (
	"github.com/shylinux/icebergs"
	"github.com/shylinux/icebergs/base/cli"
	"github.com/shylinux/toolkits"
	"os"
	"path"
)

var Index = &ice.Context{Name: "cluster", Help: "cluster",
	Caches: map[string]*ice.Cache{},
	Configs: map[string]*ice.Config{
		"cluster": {Name: "cluster", Help: "cluster", Value: kit.Data(
			kit.MDB_SHORT, "name",
			"template", "usr/cluster/bin",
			"display", "usr/cluster",
		)},
	},
	Commands: map[string]*ice.Command{
		ice.ICE_INIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},
		ice.ICE_EXIT: {Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {}},

		"cluster": {Name: "cluster begin count", Help: "cluster", Hand: func(m *ice.Message, c *ice.Context, cmd string, arg ...string) {
			for i := kit.Int(arg[0]); i < kit.Int(arg[1]); i++ {
				p := path.Join(m.Conf("cluster", "meta.display"), kit.Format("%d-master", i))
				if e := os.Mkdir(p, 0777); os.IsExist(e) {
					continue
				}
				m.Cmd("nfs.copy", path.Join(p, "cluster.conf"), path.Join(m.Conf("cluster", "meta.template"), "cluster.conf"))
			}
		}},
	},
}

func init() { cli.Index.Register(Index, nil) }
