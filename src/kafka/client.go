package kafka

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"

	"shylinux.com/x/ice"
	"shylinux.com/x/icebergs/base/aaa"
	"shylinux.com/x/icebergs/base/mdb"
	"shylinux.com/x/icebergs/base/tcp"
	kit "shylinux.com/x/toolkits"
)

const (
	TOPIC = "topic"
	KEYS  = "keys"
)

type client struct {
	ice.Zone
	short string `data:"sess"`
	field string `data:"time,sess,host,port,topic,count"`

	create string `name:"create sess=biz host=localhost port=10003 topic=TASK_AGENT" help:"创建"`
	send   string `name:"send sess=biz keys=hi value=hello" help:"发送"`
	list   string `name:"list sess@key id auto" help:"消息队列"`
}

func (s client) Inputs(m *ice.Message, arg ...string) {
	switch arg[0] {
	case tcp.PORT:
		m.Cmdy(tcp.SERVER).Cut("port,status,time")
	default:
		if strings.Contains(m.Config(mdb.FIELD), arg[0]) {
			s.Hash.Inputs(m, arg...)
		} else {
			s.Zone.Inputs(m, arg...)
		}
	}
}
func (s client) Create(m *ice.Message, arg ...string) {
	s.Hash.Create(m)
	sess, topic := m.Option(aaa.SESS), m.Option(TOPIC)
	s.Hash.Target(m, sess, func() ice.Any {
		r := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{kit.Format("%s:%s", kit.Select(tcp.LOCALHOST, m.Option(tcp.HOST)), m.Option(tcp.PORT))}, Topic: topic})
		r.SetOffset(-1)
		m.Go(func() {
			for {
				if msg, err := r.ReadMessage(context.Background()); !m.Warn(err, msg) {
					s.Insert(m, aaa.SESS, sess, mdb.TIME, msg.Time.Local().Format(ice.MOD_TIME), KEYS, string(msg.Key), mdb.VALUE, string(msg.Value))
				} else {
					break
				}
			}
		})
		return r
	})
}
func (s client) Send(m *ice.Message, arg ...string) {
	s.Hash.List(m, m.Option(mdb.ZONE))
	w := &kafka.Writer{Addr: kafka.TCP(kit.Format("%s:%s", kit.Select(tcp.LOCALHOST, m.Append(tcp.HOST)), m.Append(tcp.PORT))), Topic: m.Append(TOPIC)}
	defer w.Close()
	m.SetAppend()
	m.Warn(w.WriteMessages(context.Background(), kafka.Message{Key: []byte(m.Option(KEYS)), Value: []byte(m.Option(mdb.VALUE))}))
}
func (s client) List(m *ice.Message, arg ...string) {
	switch len(kit.Slice(arg, 0, 2)) {
	case 0:
		s.Hash.List(m).Action(s.Create)
	case 1:
		m.OptionFields("time,id,keys,value")
		fallthrough
	default:
		mdb.OptionPage(m.Message, kit.Slice(arg, 2)...)
		s.Zone.List(m, kit.Slice(arg, 0, 2)...).Action(s.Send, mdb.PAGE)
	}
}

func init() { ice.CodeCtxCmd(client{}) }
