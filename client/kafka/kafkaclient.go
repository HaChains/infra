package kafkaclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
)

type Message struct {
	Meta any
	K    string
	V    proto.Message
}

type Config struct {
	// producer, consumer

	Addr         string
	DefaultTopic string
	UseSsl       bool
	CaCert       string
	ClientCert   string
	ClientKey    string

	// consumer

	GroupID       string
	ReadBatchSize int
}

type KafkaClient struct {
	ctx context.Context

	client *kgo.Client
	config *Config
}

func New(ctx context.Context, config *Config) (*KafkaClient, error) {
	if config.Addr == "" {
		return nil, fmt.Errorf("kafka addr is empty")
	}
	if config.DefaultTopic == "" {
		return nil, fmt.Errorf("kafka topic is empty")
	}

	brokers := strings.Split(config.Addr, ",")
	var tlsConfig *tls.Config
	if config.UseSsl {
		if config.CaCert == "" || config.ClientCert == "" || config.ClientKey == "" {
			return nil, fmt.Errorf("caCert, clientCert or clientKey must be provided when useSsl is true")
		}

		ca, err := os.ReadFile(config.CaCert)
		if err != nil {
			return nil, fmt.Errorf("read kafka ca cert failed: [%w]", err)
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(ca) {
			return nil, fmt.Errorf("append kafka ca cert failed")
		}

		cert, err := tls.LoadX509KeyPair(config.ClientCert, config.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("load kafka client cert and key failed: [%w]", err)
		}

		tlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: false,
			RootCAs:            pool,
			MinVersion:         tls.VersionTLS12,
		}
	}
	ops := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(config.DefaultTopic),
		kgo.FetchMinBytes(1),                         // 1B
		kgo.FetchMaxBytes(50 * 1024 * 1024),          // 50MB
		kgo.BrokerMaxReadBytes(100 * 1024 * 1024),    // 100MB
		kgo.BrokerMaxWriteBytes(100 * 1024 * 1024),   // 100MB
		kgo.ProducerBatchMaxBytes(100 * 1024 * 1024), // 100MB
		kgo.AllowAutoTopicCreation(),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchCompression(kgo.Lz4Compression()),
		//kgo.InstanceID(""),
	}
	if tlsConfig != nil {
		ops = append(ops, kgo.DialTLSConfig(tlsConfig))
	}
	if config.GroupID != "" {
		ops = append(
			ops,
			kgo.ConsumerGroup(config.GroupID),
			kgo.DisableAutoCommit(),
		)
	}

	client, err := kgo.NewClient(ops...)
	if err != nil {
		return nil, err
	}

	// read batch size must be at least 1
	config.ReadBatchSize = max(config.ReadBatchSize, 1)
	//bufferSize := max(2*config.ReadBatchSize, 50)
	return &KafkaClient{
		ctx:    ctx,
		client: client,
		config: config,
	}, nil
}
