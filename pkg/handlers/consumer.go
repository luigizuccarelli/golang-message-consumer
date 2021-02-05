package handlers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/trackmate-message-consumer/pkg/connectors"
	"gitea-devops-shared-threefld-cicd.apps.c4.us-east-1.dev.aws.ocp.14west.io/cicd/trackmate-message-consumer/pkg/schema"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	gocb "github.com/couchbase/gocb/v2"
)

// Init : public function that connects to the kafka queue and redis cache
func Init(conn connectors.Clients) error {

	var run bool = true
	var err error
	var msg *kafka.Message
	cw := conn.KafkaConsumer()

	for run == true {
		msg, err = cw.Consumer.ReadMessage(-1)
		if err == nil {
			err = postToDB(conn, msg)
		} else {
			// The client will automatically try to recover from all errors.
			conn.Error("Consumer error: %v (%v)\n", err, msg)
		}
		test, _ := strconv.ParseBool(os.Getenv("TESTING"))
		if test == true {
			run = false
		}
	}

	cw.Consumer.Close()
	return err
}

// postToDB : private utility function that posts the json payload to couchbase
func postToDB(conn connectors.Clients, msg *kafka.Message) error {

	var analytics *schema.Trackmate

	// check if we have the updated detached json
	if msg != nil {
		payload, _ := url.PathUnescape(string(msg.Value))
		conn.Trace(fmt.Sprintf("Data from message queue %s", payload))

		// we have the new format
		errs := json.Unmarshal(msg.Value, &analytics)
		if errs != nil {
			conn.Error("postToDB unmarshalling new format %v", errs)
			return errs
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
