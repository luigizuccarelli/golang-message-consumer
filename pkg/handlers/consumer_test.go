package handlers

import (
	"bytes"
	"fmt"
	"github.com/microlib/simple"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	logger     simple.Logger
	connectors Clients
	// create a key value map (to fake redis)
	store map[string]string
)

type Clients interface {
	writeToKVStore(key string, b []byte) error
	postToAnalytics(b []byte) error
	Get(string) (string, error)
	Set(string, string, time.Duration) (string, error)
}

// FakeRedis
type FakeRedis struct {
}

type FakeQ struct {
}

// fake redis Get
func (r *Connectors) Get(key string) (string, error) {
	return store[key], nil
}

// fake redis Set
func (r *Connectors) Set(key string, value string, expr time.Duration) (string, error) {
	store[key] = value
	return string(expr), nil
}

type Connectors struct {
	http  *http.Client
	redis FakeRedis
	name  string
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// NewClientConnectors - inject our test connectors
func NewClientConnectors(filename string, code int) *Connectors {

	// initialise our store (cache)
	store = make(map[string]string)
	// in initialise the store
	store["latest"] = "test"

	// we first load the json payload to simulate response data
	// for now just ignore failures.
	file, _ := ioutil.ReadFile(filename)
	logger.Trace(fmt.Sprintf("File %s with data %s", filename, string(file)))

	// Execute our Init function to connect to Kafka and Redis
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(string(file))),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	redisclient := FakeRedis{}
	con := &Connectors{redis: redisclient, http: httpclient, name: "test"}
	return con
}

func TestAll(t *testing.T) {

	// create anonymous struct
	tests := []struct {
		Name     string
		Payload  string
		Handler  string
		FileName string
		Want     bool
		ErrorMsg string
	}{
		{
			"Test envars : should pass",
			"",
			"TestEnvarsFail",
			"tests/payload-example.json",
			true,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"Test envars : should pass",
			"",
			"TestEnvarsPass",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"Read message send to url(post) should pass",
			"",
			"ReadMessagePost",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"Read message send to redis should pass",
			"",
			"ReadMessageStore",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
	}
	var err error
	for _, tt := range tests {
		fmt.Println(fmt.Sprintf("Executing test : %s \n", tt.Name))
		switch tt.Handler {
		case "TestEnvarsFail":
			err = nil
			os.Setenv("LOG_LEVEL", "trace")
			os.Setenv("SERVER_PORT", "")
			os.Setenv("REDIS_HOST", "")
			err = CheckEnvar("LOG_LEVEL", false)
			err = CheckEnvar("REDIS_HOST", true)
		case "TestEnvarsPass":
			err = nil
			os.Setenv("LOG_LEVEL", "trace")
			os.Setenv("SERVER_PORT", "")
			os.Setenv("REDIS_HOST", "127.0.0.1")
			err = CheckEnvar("SERVER_PORT", false)
			err = CheckEnvar("REDIS_HOST", true)
		case "ReadMessagePost":
			connectors = NewClientConnectors(tt.FileName, 200)
			os.Setenv("LOG_LEVEL", "trace")
			os.Setenv("SERVER_PORT", "")
			os.Setenv("REDIS_HOST", "127.0.0.1")
			os.Setenv("REDIS_PORT", "6379")
			os.Setenv("REDIS_PASSWORD", "pt")
			os.Setenv("URL", "http://127.0.0.1:7001")
			os.Setenv("BROKERS", "localhost:9092")
			os.Setenv("TOPIC", "analytics")
			os.Setenv("TESTING", "true")
			Init()
		case "ReadMessageStore":
			// Execute our Init function to connect to Kafka and Redis
			os.Setenv("LOG_LEVEL", "trace")
			os.Setenv("SERVER_PORT", "")
			os.Setenv("REDIS_HOST", "127.0.0.1")
			os.Setenv("REDIS_PORT", "6379")
			os.Setenv("REDIS_PASSWORD", "pt")
			os.Setenv("CONNECTOR", "store")
			os.Setenv("URL", "http://127.0.0.1:7001")
			os.Setenv("BROKERS", "localhost:9092")
			os.Setenv("TOPIC", "analytics")
			os.Setenv("TESTING", "true")
			Init()
		}

		if !tt.Want {
			if err != nil {
				t.Errorf(fmt.Sprintf(tt.ErrorMsg, tt.Handler, err, nil))
			}
		} else {
			if err == nil {
				t.Errorf(fmt.Sprintf(tt.ErrorMsg, tt.Handler, "nil", "error"))
			}
		}

	}
}
