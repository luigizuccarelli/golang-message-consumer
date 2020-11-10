// +build fake

package connectors

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	gocb "github.com/couchbase/gocb/v2"
	"github.com/microlib/simple"
)

type Connectors struct {
	Bucket *FakeCouchbase
	Logger *simple.Logger
	Kafka  *KafkaConsumerWrapper
	Flag   string
	Name   string
}

type FakeCouchbase struct {
}

type FakeKafkaConsumer struct {
	Flag string
	File string
}

type KafkaConsumerWrapper struct {
	Consumer *FakeKafkaConsumer
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
	if r.Flag == "errorDB" {
		return &gocb.MutationResult{}, errors.New("forced upsert error")
	}
	return &gocb.MutationResult{}, nil
}

func (c *Connectors) KafkaConsumer() *KafkaConsumerWrapper {
	return c.Kafka
}

func (c *Connectors) Close() {
	// do nothing
}

func (f *FakeKafkaConsumer) Close() {
}

func (f *FakeKafkaConsumer) ReadMessage(int) (*kafka.Message, error) {
	if f.Flag == "error" {
		return &kafka.Message{}, errors.New("forced read message error")
	}

	b, _ := ioutil.ReadFile(f.File)
	return &kafka.Message{Value: b}, nil
}

// NewTestClientConnectors - inject our test connectors
func NewTestClientConnectors(filename string, err string, logger *simple.Logger) Clients {

	// we first load the json payload to simulate response data
	// for now just ignore failures.

	// reference our fake consumer
	fc := &FakeKafkaConsumer{Flag: err, File: filename}
	cw := &KafkaConsumerWrapper{Consumer: fc}
	// we use this flag to inject/force errors
	//if err != "error" {

	conn := &Connectors{Kafka: cw, Bucket: &FakeCouchbase{}, Logger: logger, Name: "test", Flag: err}
	return conn
}
