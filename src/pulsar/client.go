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
	ice.Zone
	short string `data:"sess"`
	field string `data:"time,sess,host,port,token,topic,count"`

	create string `name:"create sess=biz host=localhost port=10004 token topic='public/default/my-topic' group" help:"创建"`
	send   string `name:"send sess keys=hi value:textarea=hello" help:"发送"`
	list   string `name:"list sess@key id auto" help:"消息队列"`
}

func (s client) Client(m *ice.Message, host, port, token string) pulsar.Client {
	options := pulsar.ClientOptions{URL: kit.Format("pulsar://%s:%s", kit.Select(tcp.LOCALHOST, host), port)}
	if token != "" {
		options.Authentication = pulsar.NewAuthenticationToken(token)
	}
	client, e := pulsar.NewClient(options)
	m.Assert(e)
	return client
}
func (s client) Create(m *ice.Message, arg ...string) {
	s.Hash.Create(m)
	client := s.Client(m, m.Option(tcp.HOST), m.Option(tcp.PORT), m.Option(TOKEN))
	c, e := client.Subscribe(pulsar.ConsumerOptions{Topic: PERSISTENT + m.Option(TOPIC), SubscriptionName: kit.Select(ice.Info.HostName, m.Option(GROUP)), Type: pulsar.Shared})
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
	msg := &pulsar.ProducerMessage{Key: m.Option(KEYS), Payload: []byte(m.Option(mdb.VALUE))}

	s.Hash.List(m, m.Option(aaa.SESS))
	client := s.Client(m, m.Append(tcp.HOST), m.Append(tcp.PORT), m.Append(TOKEN))
	p, e := client.CreateProducer(pulsar.ProducerOptions{Topic: PERSISTENT + m.Append(TOPIC)})
	m.Assert(e)

	msgid, e := p.Send(context.Background(), msg)
	m.Push(MSGID, msgid)
	m.SetAppend()
	m.Assert(e)
}

func (s client) List(m *ice.Message, arg ...string) {
	switch len(kit.Slice(arg, 0, 2)) {
	case 0:
		s.Hash.List(m).Action(s.Create)
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
