package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/microlib/simple"
)

// Init : public function that connects to the kafka queue and redis cache / couchbase DB
func Init(logger *simple.Logger) {

	cfg := sarama.NewConfig()
	cfg.ClientID = "go-kafka-consumer"
	cfg.Consumer.Return.Errors = true

	// set the logger level
	logger.Level = os.Getenv("LOG_LEVEL")

	connectors = NewClientConnectors("NA", 0)

	// check by way of logging the kafka brokers in an HA setup
	brokerList := strings.Split(os.Getenv("BROKERS"), ",")
	logger.Info(fmt.Sprintf("Kafka brokers: %s", strings.Join(brokerList, ", ")))

	// Create new consumer
	mc, err := sarama.NewConsumer(brokerList, cfg)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mc.Close(); err != nil {
			panic(err)
		}
	}()

	topics, _ := mc.Topics()

	consumer, errors := consume(topics, mc)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// crerate a chanel for our consumer messages
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				msgCount++
				logger.Debug(fmt.Sprintf("Received messages %s %s\n", string(msg.Key), string(msg.Value)))
				if os.Getenv("TESTING") == "true" && msgCount > 1 {
					logger.Info("Test flag set - auto interrupt")
					doneCh <- struct{}{}
				}
			case consumerError := <-errors:
				msgCount++
				logger.Error(fmt.Sprintf("Received consumerError %s %s %v \n", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err))
				doneCh <- struct{}{}
			case <-signals:
				logger.Info("Interrupt detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	logger.Info(fmt.Sprintf("Processed %d messages", msgCount))
}

// consume function - it iterates through each topic to find the specified topic, once found it then iterates through each partition
func consume(topics []string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {

	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)
	for _, topic := range topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}
		// We are only interested in the specified topic
		if topic == os.Getenv("TOPIC") {
			partitions, _ := master.Partitions(topic)
			for x, _ := range partitions {
				// consider using sarama.OffsetNewest
				consumer, err := master.ConsumePartition(topic, partitions[x], sarama.OffsetOldest)
				if nil != err {
					logger.Error(fmt.Sprintf("Topic %v Partition: %v", topic, partitions[x]))
					break
				}
				logger.Info(fmt.Sprintf("Start consuming topic %v ", topic))
				go func(topic string, consumer sarama.PartitionConsumer) {
					for {
						select {
						case consumerError := <-consumer.Errors():
							errors <- consumerError
							logger.Error(fmt.Sprintf("consumerError: %v ", consumerError.Err))

						case msg := <-consumer.Messages():
							consumers <- msg
							logger.Debug(fmt.Sprintf("Got message on topic %v : %v ", topic, msg.Value))

							if os.Getenv("CONNECTOR") == "cache" {
								// write to cache
								logger.Debug(fmt.Sprintf("Writing data to redis cache"))
								err := connectors.writeToKVStore("latest", msg.Value)
								if err != nil {
									logger.Error(fmt.Sprintf("Error : %v ", err))
								}
							} else {
								// write to nosql db
								logger.Debug(fmt.Sprintf("Writing data to couchbase "))
								err := connectors.postToDB(msg.Value)
								if err != nil {
									logger.Error(fmt.Sprintf("Error : %v ", err))
								}
							}
						}
					}
				}(topic, consumer)
			}
		}
	}
	return consumers, errors
}

// writeToKVStore : private utility function that writes the payload to a KV store (redis)
func (c *Connectors) writeToKVStore(key string, b []byte) error {
	_, err := c.Set(key, string(b), 0)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not write to kv store key = %s", key))
		return err
	} else {
		logger.Debug(fmt.Sprintf("Data written to kv store key = %s", key))
	}
	return nil
}

// postToDB : private utility function that posts the json payload to couchbase

func (c *Connectors) postToDB(b []byte) error {

	var analytics Analytics

	// we first unmarshal the payload and add needed values before posting to couchbase
	errs := json.Unmarshal(b, &analytics)
	if errs != nil {
		logger.Error(fmt.Sprintf("Could not unmarshal analytics data to json %v", errs))
		return errs
	}

	analytics.ProductName = "Trackmate"

	// ensure uniqueness
	collection := r.Bucket.DefaultCollection()
	_, err := collection.Upsert(analytics.TrackingId+"-"+analytics.AffiliateId, analytics, &gocb.CollectionOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("Could not upsert schema into couchbase %v", err))
		return err
	}

	// all good :)
	logger.Debug(fmt.Sprintf("Analytics schema inserted into couchbase  %v \n", analytics))
	return nil
}
