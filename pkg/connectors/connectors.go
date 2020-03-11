package connectors

import (
	"os"
	"time"

	gocb "github.com/couchbase/gocb/v2"
	"github.com/go-redis/redis"
)

// Clients interface - the NewClientConnectors will return this struct
type Clients interface {
	writeToKVStore(key string, b []byte) error
	postToDB(b []byte) error
	Get(string) (string, error)
	Set(string, string, time.Duration) (string, error)
}

// The premise here is to use this as a reciever in the relevant functions
// this allows us then to mock/fake connections and calls
type Connectors struct {
	Bucket  *gocb.Bucket
	Cluster *gocb.Cluster
	redis   *redis.Client
	name    string
}

// NewClientConnectors : function that initialises connections to DB's, caches' queues etc
// Seperating this functionality here allows us to inject a fake or mock connection object for testing
func NewClientConnectors(filename string, code int) Clients {
	// connect to redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           0,
	})

	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: os.Getenv("COUCHBASE_USER"),
			Password: os.Getenv("COUCHBASE_PASSWORD"),
		},
	}

	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_HOST"), opts)
	if err != nil {
		panic(err)
	}

	// get a bucket reference
	bucket := cluster.Bucket(os.Getenv("COUCHBASE_BUCKET"), &gocb.BucketOptions{})
	return &Connectors{redis: redisClient, Bucket: bucket, Cluster: cluster, name: "RealConnectors"}
}

// Get : wrapper function for redis get
func (r *Connectors) Get(key string) (string, error) {
	val, err := r.redis.Get(key).Result()
	return val, err
}

// Set : wrapper function for redis set
func (r *Connectors) Set(key string, value string, expr time.Duration) (string, error) {
	val, err := r.redis.Set(key, value, expr).Result()
	return val, err
}

// Close : wrapper function for redis and couchbase
func (r *Connectors) Close() error {
	r.redis.Close()
	r.Cluster.Close(&gocb.ClusterCloseOptions{})
	return nil
}
