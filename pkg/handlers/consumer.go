package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-message-consumer/pkg/connectors"
	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-message-consumer/pkg/schema"
	"github.com/Shopify/sarama"
	gocb "github.com/couchbase/gocb/v2"
)

// Init : public function that connects to the kafka queue and redis cache
func Init(conn connectors.Clients) {

	mc := conn.KafkaConsumer()

	defer func() {
		if err := mc.Close(); err != nil {
			panic(err)
		}
	}()

	topics, _ := mc.Topics()

	consumer, errors := consume(conn, topics, mc)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// create a chanel for our consumer messages
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				msgCount++
				conn.Trace(fmt.Sprintf("Received messages %v ", msg))
				if os.Getenv("TESTING") == "true" && msgCount > 1 {
					conn.Info("Test flag set - auto interrupt")
					doneCh <- struct{}{}
				}
			case consumerError := <-errors:
				msgCount++
				conn.Error(fmt.Sprintf("Received consumerError  %v ", consumerError))
				doneCh <- struct{}{}
			case <-signals:
				conn.Info("Interrupt detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	conn.Info(fmt.Sprintf("Processed %d messages", msgCount))
}

// consume function - it iterates through each topic to find the specified topic, once found it then iterates through each partition
func consume(conn connectors.Clients, topics []string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {

	conn.Info(fmt.Sprintf("Function consume topics %v", topics))
	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)
	for _, topic := range topics {
		conn.Info(fmt.Sprintf("iterate topics %v", topics))
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
					conn.Error(fmt.Sprintf("Topic %v Partition: %v", topic, partitions[x]))
					break
				}
				conn.Info(fmt.Sprintf("Start consuming topic %v ", topic))
				go func(topic string, consumer sarama.PartitionConsumer) {
					for {
						select {
						case consumerError := <-consumer.Errors():
							errors <- consumerError
							conn.Error(fmt.Sprintf("consumerError: %v ", consumerError))

						case msg := <-consumer.Messages():
							consumers <- msg
							conn.Trace(fmt.Sprintf("Got message on topic %v : %v ", topic, msg))

							err := postToDB(conn, msg)
							if err != nil {
								conn.Error(fmt.Sprintf("Error : %v ", err))
							}
						}
					}
				}(topic, consumer)
			}
		}
	}
	return consumers, errors
}

// postToDB : private utility function that posts the json payload to couchbase
func postToDB(conn connectors.Clients, msg *sarama.ConsumerMessage) error {

	var analytics *schema.Trackmate

	// check if we have the updated detached json from segmentio
	if msg != nil {
		payload := string(msg.Value)
		conn.Trace(fmt.Sprintf("Data from message queue %s", payload))

		// we have the new format
		errs := json.Unmarshal(msg.Value, &analytics)
		if errs != nil {
			conn.Error("postToDB unmarshalling new format %v", errs)
			return errs
		}

		// use regex to extract pagename (for easier aggregation)
		//var validID = regexp.MustCompile(`pagename=[a-z]*`)
		//s := validID.FindAllString(analytics.Page.Referrer, -1)
		//conn.Debug(fmt.Sprintf("Regex referrer pagename  %s", s))
		//if len(s) > 0 {
		if strings.Contains(analytics.Page.Referrer, "pagename=") {
			validID := regexp.MustCompile(`pagename=[a-z]*`)
			s := validID.FindAllString(analytics.Page.Referrer, -1)
			conn.Debug(fmt.Sprintf("Regex referrer pagename  %s", s))
			analytics.Page.ReferrerName = strings.Split(s[0], "=")[1]
		} else {
			// this means we don't have the referrerName - we need to do some regex updates
			if strings.Contains(analytics.Page.Referrer, "pro.") {
				analytics.Page.ReferrerName = "promo"
			} else if strings.Contains(analytics.Page.Referrer, "orders.") {
				analytics.Page.ReferrerName = "order"
			} else if strings.Contains(analytics.Page.Referrer, "opium.") {
				analytics.Page.ReferrerName = "opium"
			} else {
				// take the domain name
				hld := strings.Split(analytics.Page.Referrer, "/")
				if len(hld) > 2 {
					analytics.Page.ReferrerName = hld[2]
				} else {
					analytics.Page.ReferrerName = "none"
				}
			}

		}

		if strings.Contains(analytics.Page.URL, "pagename=") {
			validID := regexp.MustCompile(`pagename=[a-z]*`)
			s := validID.FindAllString(analytics.Page.URL, -1)
			conn.Debug(fmt.Sprintf("Regex referrer pagename  %s", s))
			analytics.Page.URLName = strings.Split(s[0], "=")[1]
		} else {
			// this means we don't have the referrerName - we need to do some regex updates
			if strings.Contains(analytics.Page.URL, "pro.") {
				analytics.Page.URLName = "promo"
			} else if strings.Contains(analytics.Page.URL, "orders.") {
				analytics.Page.URLName = "order"
			} else if strings.Contains(analytics.Page.URL, "opium.") {
				analytics.Page.URLName = "opium"
			} else {
				// take the domain name
				hld := strings.Split(analytics.Page.URL, "/")
				if len(hld) > 2 {
					analytics.Page.URLName = hld[2]
				} else {
					analytics.Page.URLName = "none"
				}
			}
		}

		conn.Debug(fmt.Sprintf("Analytics struct  %v", analytics))
		res, err := conn.Upsert(analytics.MessageId, analytics, &gocb.UpsertOptions{})
		conn.Debug(fmt.Sprintf("Result from insert %v", res))
		if err != nil {
			conn.Error(fmt.Sprintf("Could not insert schema into couchbase %v", err))
			return err
		}

		// all good :)
		conn.Info("Analytics schema inserted into couchbase")
		return nil

	} else {
		conn.Info("Message data is nil")
		return nil
	}
}
