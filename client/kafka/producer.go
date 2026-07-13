package kafkaclient

import (
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	*KafkaClient
	topic  string
	chSend chan *Message
	commit func(m *Message) error
}

func (kc *KafkaClient) Producer(topic string, commit func(m *Message) error) *Producer {
	if topic == "" {
		topic = kc.config.DefaultTopic
	}
	return &Producer{
		KafkaClient: kc,
		topic:       topic,
		chSend:      make(chan *Message, 10),
		commit:      commit,
	}
}

func (p *Producer) Produce(m *Message) error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	default:
	}

	// TODO: Signal when Run exits so Produce does not block forever on a full
	// queue if Run returns an error before the parent context is canceled.
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case p.chSend <- m:
		return nil
	}
}

func (p *Producer) round() error {
	select {
	case <-p.ctx.Done():
		return p.ctx.Err()
	case msg := <-p.chSend:
		v, err := proto.Marshal(msg.V)
		if err != nil {
			return err
		}
		record := &kgo.Record{
			Topic: p.topic,
			Key:   []byte(msg.K),
			Value: v,
		}
		for {
			err = p.client.ProduceSync(p.ctx, record).FirstErr()
			if err == nil {
				break
			}
			select {
			case <-p.ctx.Done():
				return p.ctx.Err()
			case <-time.After(time.Second):
			}
		}
		return p.commit(msg)
	}
}

func (p *Producer) Run() (err error) {
	for {
		err = p.round()
		if err != nil {
			return err
		}
	}
}
