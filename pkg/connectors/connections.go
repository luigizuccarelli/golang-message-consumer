package connectors

import (
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

// Clients interface - the NewClientConnectors will return this struct
type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	Upsert(col string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error)
	//Upsert(col string, value interface{}, expiry uint32) (gocb.Cas, error)
	KafkaConsumer() sarama.Consumer
	Close()
}

// The premise here is to use this as a reciever in the relevant functions
// this allows us then to mock/fake connections and calls
type Connectors struct {
	Bucket *gocb.Bucket
	Logger *simple.Logger
	Kafka  sarama.Consumer
	Name   string
}

// NewClientConnectors : function that initialises connections to DB's, caches' queues etc
// Seperating this functionality here allows us to inject a fake or mock connection object for testing
func NewClientConnectors(logger *simple.Logger) Clients {

	logger.Info(fmt.Sprintf("Couchbase envars: %s %s %s", os.Getenv("COUCHBASE_HOST"), os.Getenv("COUCHBASE_USER"), os.Getenv("COUCHBASE_PASSWORD")))

	opts := gocb.ClusterOptions{
		Username: os.Getenv("COUCHBASE_USER"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	}
	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_HOST"), opts)
	if err != nil {
		logger.Error(fmt.Sprintf("Couchbase connection: %v", err))
		panic(err)
	}

	// get a bucket reference
	// bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"), &gocb.BucketOptions{}) v.2.0.0-beta-1
	bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"))
	logger.Info(fmt.Sprintf("Couchbase connection: %v", bucket))

	// kafka connector
	cfg := sarama.NewConfig()
	cfg.ClientID = "go-kafka-consumer"
	cfg.Consumer.Return.Errors = true

	// check by way of logging the kafka brokers in an HA setup
	brokerList := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	logger.Info(fmt.Sprintf("Kafka brokers: %s", strings.Join(brokerList, ", ")))

	// Create new consumer
	mc, err := sarama.NewConsumer(brokerList, cfg)
	if err != nil {
		panic(err)
	}
	return &Connectors{Bucket: bucket, Kafka: mc, Logger: logger, Name: "RealConnectors"}
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
	collection := r.Bucket.DefaultCollection()
	return collection.Upsert(col, value, opts)
}

func (c *Connectors) KafkaConsumer() sarama.Consumer {
	return c.Kafka
}

func (c *Connectors) Close() {
	c.Kafka.Close()
}
