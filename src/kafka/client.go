package kafka

import (
	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
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
	field string `data:"time,id,topic,text"`

	create string `name:"create cluster server topic group" help:"创建"`
	send   string `name:"send cluster topic=TASK_AGENT text=hello" help:"发送"`
	list   string `name:"list cluster@key id auto send" help:"消息队列"`
}

func (s client) Create(m *ice.Message, arg ...string) {
	s.Zone.Create(m, m.OptionSimple(CLUSTER, SERVER, TOPIC, GROUP)...)

	c, e := kafka.NewConsumer(&kafka.ConfigMap{"bootstrap.servers": m.Option(SERVER), "group.id": m.Option(GROUP)})
	m.Assert(e)

	cluster, topic := m.Option(CLUSTER), m.Option(TOPIC)
	m.Assert(c.SubscribeTopics([]string{topic}, nil))

	m.Go(func() {
		for {
			if msg, err := c.ReadMessage(-1); !m.Warn(err, msg) {
				s.Insert(m, CLUSTER, cluster, mdb.TIME, msg.Timestamp.Local().Format(ice.MOD_TIME), TOPIC, topic, mdb.TEXT, string(msg.Value))
			}
		}
	})
}

func (s client) Send(m *ice.Message, arg ...string) {
	msg := m.Cmd(mdb.SELECT, m.PrefixKey(), "", mdb.HASH, m.OptionSimple(CLUSTER), kit.Dict(ice.MSG_FIELDS, "time,cluster,server,topic,group"))

	p, e := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": msg.Append(SERVER)})
	m.Assert(e)
	defer p.Close()
	defer p.Flush(100)

	topic := m.Option(TOPIC)
	m.Assert(p.Produce(&kafka.Message{Value: []byte(m.Option(mdb.TEXT)),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	}, nil))
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
