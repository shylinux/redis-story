package kafka

import (
	"context"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"

	"github.com/segmentio/kafka-go"
)

const (
	CLUSTER = "cluster"
	SERVER  = "server"
	TOPIC   = "topic"
	GROUP   = "group"
)

type client struct {
	ice.Zone
	short string `data:"cluster"`
	field string `data:"time,id,topic,keys,value"`

	create string `name:"create cluster server topic group" help:"创建"`
	send   string `name:"send cluster topic=TASK_AGENT keys=hi value=hello" help:"发送"`
	list   string `name:"list cluster@key id auto send" help:"消息队列"`
}

func (s client) Create(m *ice.Message, arg ...string) {
	s.Zone.Create(m, m.OptionSimple(CLUSTER, SERVER, TOPIC, GROUP)...)

	cluster, topic := m.Option(CLUSTER), m.Option(TOPIC)
	r := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{m.Option(SERVER)}, Topic: topic})
	r.SetOffset(-1)

	m.Go(func() {
		for {
			if msg, err := r.ReadMessage(context.Background()); !m.Warn(err, msg) {
				s.Insert(m, CLUSTER, cluster, mdb.TIME, msg.Time.Local().Format(ice.MOD_TIME), TOPIC, topic, "keys", string(msg.Key), mdb.VALUE, string(msg.Value))
			}
		}
	})
}

func (s client) Send(m *ice.Message, arg ...string) {
	msg := m.Cmd(mdb.SELECT, m.PrefixKey(), "", mdb.HASH, m.OptionSimple(CLUSTER), kit.Dict(ice.MSG_FIELDS, "time,cluster,server,topic,group"))

	w := &kafka.Writer{Addr: kafka.TCP(msg.Append(SERVER)), Topic: m.Option(TOPIC)}
	defer w.Close()

	m.Assert(w.WriteMessages(context.Background(), kafka.Message{Key: []byte(m.Option("keys")), Value: []byte(m.Option(mdb.VALUE))}))
}

func (s client) List(m *ice.Message, arg ...string) {
	switch len(kit.Slice(arg, 0, 2)) {
	case 0:
		m.OptionFields("time,cluster,server,topic,group")
	case 1:
		m.OptionPage(kit.Slice(arg, 2)...)
		m.Action(mdb.PAGE)
	}
	s.Zone.List(m, kit.Slice(arg, 0, 2)...)
}

func init() { ice.CodeCtxCmd(client{}) }
