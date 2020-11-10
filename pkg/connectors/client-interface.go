package connectors

import (
	gocb "github.com/couchbase/gocb/v2"
)

// Clients interface - the NewClientConnectors will return this struct
type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	Upsert(col string, value interface{}, opts *gocb.UpsertOptions) (*gocb.MutationResult, error)
	KafkaConsumer() *KafkaConsumerWrapper
	Close()
}
