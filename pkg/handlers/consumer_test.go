package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-message-consumer/pkg/connectors"
	"github.com/Shopify/sarama"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

type Connectors struct {
	Bucket *FakeCouchbase
	Logger *simple.Logger
	Kafka  sarama.Consumer
	Name   string
}

type FakeCouchbase struct {
}

func (conn Connectors) Close() {
}

func (r *Connectors) Error(msg string, val ...interface{}) {
	r.Logger.Error(fmt.Sprintf(msg, val...))
}

func (r *Connectors) Info(msg string, val ...interface{}) {
	r.Logger.Info(fmt.Sprintf(msg, val...))
}

func (r *Connectors) Debug(msg string, val ...interface{}) {
	r.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (r *Connectors) Trace(msg string, val ...interface{}) {
	r.Logger.Trace(fmt.Sprintf(msg, val...))
}

// Upsert : wrapper function for couchbase update
func (r *Connectors) Upsert(col string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error) {
	return &gocb.MutationResult{}, nil
}

func (c *Connectors) KafkaConsumer() sarama.Consumer {
	return c.Kafka
}

// NewTestClientConnectors - inject our test connectors
func NewTestClientConnectors(filename string, code int, err string, logger *simple.Logger) connectors.Clients {

	// we first load the json payload to simulate response data
	// for now just ignore failures.
	file, _ := ioutil.ReadFile(filename)
	logger.Trace(fmt.Sprintf("File %s with data %s", filename, string(file)))

	// reference our mock consumer (see concumer-mock.go)
	consumer := NewConsumer(logger, nil)

	// we use this flag to inject/force errors
	//if err != "error" {
	consumer.SetTopicMetadata(map[string][]int32{
		"test":  {0, 1, 2, 3},
		"test1": {0, 1, 2, 3, 4, 5, 6, 7},
	})
	//}

	consumer.ExpectConsumePartition("test", 0, sarama.OffsetOldest).YieldMessage(&sarama.ConsumerMessage{Value: file})
	consumer.ExpectConsumePartition("test", 0, sarama.OffsetOldest).YieldError(sarama.ErrOutOfBrokers)
	consumer.ExpectConsumePartition("test", 0, sarama.OffsetOldest).YieldMessage(&sarama.ConsumerMessage{Value: []byte("{hello world again}")})
	consumer.ExpectConsumePartition("test", 0, AnyOffset).YieldMessage(&sarama.ConsumerMessage{Value: file})

	logger.Trace(fmt.Sprintf("Mock consumer details %v", consumer))
	conn := &Connectors{Kafka: consumer, Bucket: &FakeCouchbase{}, Logger: logger, Name: "test"}
	return conn
}

func TestAll(t *testing.T) {

	logger := &simple.Logger{Level: "trace"}

	t.Run("Message consumer : should pass", func(t *testing.T) {
		conn := NewTestClientConnectors("../../tests/new-format.json", 200, "normal", logger)
		os.Setenv("REDIS_HOST", "redis.myportfolio.svc.cluster.local")
		os.Setenv("KAFKA_BROKERS", "my-cluster-kafka-brokers.apache-kafka.svc.cluster.local:9092")
		os.Setenv("LOG_LEVEL", "trace")
		os.Setenv("SERVER_PORT", "")
		os.Setenv("REDIS_PORT", "6379")
		os.Setenv("REDIS_PASSWORD", "pt")
		os.Setenv("URL", "http://127.0.0.1:7001/")
		os.Setenv("TOPIC", "test")
		os.Setenv("TESTING", "false")
		os.Setenv("CONNECTOR", "none")
		// call and test our consumer.go
		Init(conn)
	})

	t.Run("Message consumer : should pass", func(t *testing.T) {
		conn := NewTestClientConnectors("../../tests/new-format-no-utm.json", 200, "normal", logger)
		os.Setenv("REDIS_HOST", "redis.myportfolio.svc.cluster.local")
		os.Setenv("KAFKA_BROKERS", "my-cluster-kafka-brokers.apache-kafka.svc.cluster.local:9092")
		os.Setenv("LOG_LEVEL", "trace")
		os.Setenv("SERVER_PORT", "")
		os.Setenv("REDIS_PORT", "6379")
		os.Setenv("REDIS_PASSWORD", "pt")
		os.Setenv("URL", "http://127.0.0.1:7001/")
		os.Setenv("TOPIC", "test")
		os.Setenv("TESTING", "false")
		os.Setenv("CONNECTOR", "none")
		// call and test our consumer.go
		Init(conn)
	})

	t.Run("Message consumer : should pass", func(t *testing.T) {
		conn := NewTestClientConnectors("../../tests/new-format.json", 200, "normal", logger)
		os.Setenv("REDIS_HOST", "redis.myportfolio.svc.cluster.local")
		os.Setenv("KAFKA_BROKERS", "my-cluster-kafka-brokers.apache-kafka.svc.cluster.local:9092")
		os.Setenv("LOG_LEVEL", "trace")
		os.Setenv("SERVER_PORT", "")
		os.Setenv("REDIS_PORT", "6379")
		os.Setenv("REDIS_PASSWORD", "pt")
		os.Setenv("URL", "http://127.0.0.1:7001/")
		os.Setenv("TOPIC", "test")
		os.Setenv("TESTING", "true")
		os.Setenv("CONNECTOR", "none")
		// call and test our consumer.go
		Init(conn)
	})

	t.Run("Message consumer : should fail", func(t *testing.T) {
		conn := NewTestClientConnectors("../../tests/new-format.json", 200, "error", logger)
		os.Setenv("REDIS_HOST", "redis.myportfolio.svc.cluster.local")
		os.Setenv("KAFKA_BROKERS", "my-cluster-kafka-brokers.apache-kafka.svc.cluster.local:9092")
		os.Setenv("LOG_LEVEL", "trace")
		os.Setenv("SERVER_PORT", "")
		os.Setenv("REDIS_PORT", "6379")
		os.Setenv("REDIS_PASSWORD", "pt")
		os.Setenv("URL", "http://127.0.0.1:7001/")
		os.Setenv("TOPIC", "test")
		os.Setenv("TESTING", "false")
		os.Setenv("CONNECTOR", "none")
		// call and test our consumer.go
		Init(conn)
	})
}
