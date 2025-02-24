package pulsar

import (
	"context"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"

	"github.com/apache/pulsar-client-go/pulsar"
)

const (
	TOKEN = "token"
	TOPIC = "topic"
	GROUP = "group"
)

const (
	KEYS       = "keys"
	MSGID      = "msgid"
	PROPERTIES = "properties"
	PERSISTENT = "persistent://"
)

type client struct {
	ice.Hash
	ice.Zone
	short  string `data:"sess"`
	field  string `data:"time,id,msgid,keys,value"`
	fields string `data:"time,sess,host,port,token,topic,count"`
	create string `name:"create sess*=biz host*=localhost port*=10003 token topic='public/default/my-topic' group" help:"创建"`
	send   string `name:"send sess* keys*=hi value:textarea*=hello" help:"发送"`
	list   string `name:"list sess@key id auto" help:"消息队列"`
}

func (s client) Client(m *ice.Message, host, port, token string) pulsar.Client {
	options := pulsar.ClientOptions{URL: kit.Format("pulsar://%s:%s", kit.Select(tcp.LOCALHOST, host), port)}
	kit.If(token != "", func() { options.Authentication = pulsar.NewAuthenticationToken(token) })
	client, e := pulsar.NewClient(options)
	m.Assert(e)
	return client
}
func (s client) Create(m *ice.Message, arg ...string) {
	s.Zone.Create(m, arg...)
	client := s.Client(m, m.Option(tcp.HOST), m.Option(tcp.PORT), m.Option(TOKEN))
	c, e := client.Subscribe(pulsar.ConsumerOptions{Topic: PERSISTENT + m.Option(TOPIC), SubscriptionName: kit.Select(ice.Info.NodeName, m.Option(GROUP)), Type: pulsar.Shared})
	m.Assert(e)
	sess := m.Option(aaa.SESS)
	m.Go(func() {
		for {
			if msg, err := c.Receive(context.Background()); !m.Warn(err) {
				s.Zone.Insert(m, aaa.SESS, sess, MSGID, kit.Format("%v", msg.ID()), KEYS, msg.Key(), mdb.VALUE, string(msg.Payload()), PROPERTIES, kit.Format(msg.Properties()))
				c.Ack(msg)
			} else {
				break
			}
		}
	})
}
func (s client) Send(m *ice.Message, arg ...string) {
	msg := s.Hash.List(m.Spawn(), m.Option(aaa.SESS))
	p, e := s.Client(m, msg.Append(tcp.HOST), msg.Append(tcp.PORT), msg.Append(TOKEN)).CreateProducer(pulsar.ProducerOptions{Topic: PERSISTENT + msg.Append(TOPIC)})
	m.Assert(e)
	if msgid, e := p.Send(context.Background(), &pulsar.ProducerMessage{Key: m.Option(KEYS), Payload: []byte(m.Option(mdb.VALUE))}); !m.Warn(e) {
		m.Echo(msgid.String())
	}
}

func (s client) List(m *ice.Message, arg ...string) {
	switch len(kit.Slice(arg, 0, 2)) {
	case 0:
		s.Zone.List(m).Action(s.Create)
	case 1:
		m.OptionFields("time,id,msgid,keys,value")
		fallthrough
	default:
		if s.Zone.ListPage(m, arg...); m.FieldsIsDetail() {
			m.Append(mdb.VALUE, kit.Formats(kit.UnMarshal(m.Append(mdb.VALUE))))
		} else {
			m.Action(s.Send, mdb.PAGE)
		}
	}
}

func init() { ice.CodeCtxCmd(client{}) }
