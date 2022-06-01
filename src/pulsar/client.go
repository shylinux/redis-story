package pulsar

import (
	"context"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/mdb"
	kit "shylinux.com/x/toolkits"

	"github.com/apache/pulsar-client-go/pulsar"
)

const (
	CLUSTER = "cluster"
	TOPIC   = "topic"
	GROUP   = "group"
	SERVER  = "server"
	TOKEN   = "token"
)

const (
	KEYS       = "keys"
	PROPERTIES = "properties"
	PREFIX     = "persistent://"
)

type client struct {
	ice.Zone
	short string `data:"cluster"`
	field string `data:"time,id,keys,value"`

	create string `name:"create cluster topic group server token" help:"创建"`
	send   string `name:"send cluster keys=hi value:textarea=hello" help:"发送"`
	list   string `name:"list cluster id auto send" help:"消息队列"`
}

func (s client) Create(m *ice.Message, arg ...string) {
	s.Zone.Create(m, m.OptionSimple(CLUSTER, TOPIC, GROUP, SERVER, TOKEN)...)

	client, e := pulsar.NewClient(pulsar.ClientOptions{URL: m.Option(SERVER), Authentication: pulsar.NewAuthenticationToken(m.Option(TOKEN))})
	m.Assert(e)

	c, e := client.Subscribe(pulsar.ConsumerOptions{Topic: PREFIX + m.Option(TOPIC), SubscriptionName: m.Option(GROUP), Type: pulsar.Shared})
	m.Assert(e)

	cluster := m.Option(CLUSTER)
	m.Go(func() {
		for {
			if msg, err := c.Receive(context.Background()); !m.Warn(err) {
				s.Zone.Insert(m, CLUSTER, cluster,
					KEYS, msg.Key(), mdb.VALUE, string(msg.Payload()), PROPERTIES, kit.Format(msg.Properties()))
				c.Ack(msg)
			}
		}
	})
}

func (s client) Send(m *ice.Message, arg ...string) {
	msg := m.Cmd(mdb.SELECT, m.PrefixKey(), "", mdb.HASH, m.OptionSimple(CLUSTER), kit.Dict(ice.MSG_FIELDS, kit.Fields(TOPIC, SERVER, TOKEN)))

	client, e := pulsar.NewClient(pulsar.ClientOptions{URL: msg.Append(SERVER), Authentication: pulsar.NewAuthenticationToken(msg.Append(TOKEN))})
	m.Assert(e)

	p, e := client.CreateProducer(pulsar.ProducerOptions{Topic: PREFIX + msg.Append(TOPIC)})
	m.Assert(e)

	_, e = p.Send(context.Background(), &pulsar.ProducerMessage{Key: m.Option(KEYS), Payload: []byte(m.Option(mdb.VALUE)), Properties: map[string]string{}})
	m.Assert(e)
}

func (s client) List(m *ice.Message, arg ...string) {
	if len(arg) == 0 {
		m.OptionFields(CLUSTER, TOPIC, GROUP, SERVER, TOKEN)
	}
	if s.Zone.List(m, arg...); m.FieldsIsDetail() {
		m.Append(mdb.VALUE, kit.Formats(kit.UnMarshal(m.Append(mdb.VALUE))))
	}
}

func init() { ice.CodeCtxCmd(client{}) }
