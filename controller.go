package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_defaultTimeoutMS = 30000
)

// Controller acts as a wrapper of the MongoDB driver exposing the MongoDB operator functionality
type Controller interface {
	// Mongo Operations

	Count(collection string, filter interface{}, opts ...*options.CountOptions) (int64, error)
	DeleteMany(collection string, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	DeleteOne(collection string, filter, target interface{}, opts ...*options.FindOneAndDeleteOptions) error
	FindMany(collection string, filter, target interface{}, opts ...*options.FindOptions) error
	FindOne(colletion string, filter, target interface{}, opts ...*options.FindOneOptions) error
	InsertMany(collection string, payload []interface{}, opts ...*options.InsertManyOptions) ([]interface{}, error)
	InsertOne(collection string, payload interface{}, opts ...*options.InsertOneOptions) error
	UpdateMany(collection string, filter, payload interface{}, opts ...*options.UpdateOptions) (int64, error)
	UpdateOne(collection string, filter, payload interface{}, opts ...*options.FindOneAndUpdateOptions) error

	// Network Operations

	Ping() error

	// Getters

	Database() string
	DefaultCollection() string
}

// Operator implements the Controller interface to perform MongoDB operations
type Operator struct {
	client *mongo.Client
	config Config
}

// Config holds the connection info required to interface with MongoDB.
type Config struct {
	// Database represents the Mongo database to connect to; required.
	Database string `mapstructure:"database"`

	// DefaultCollection represents a default collection to perform operations with.
	DefaultCollection string `mapstructure:"defaultCollection"`

	// Password represents the database user's password.
	Password string `mapstructure:"password"`

	// TimeoutMS represents a contextual timeout for operations in milliseconds.
	TimeoutMS int64 `mapstructure:"timeoutMS"`

	// URI represents the database URI used to establish a connection.
	URI string `mapstructure:"uri"`

	// Username represents the database username.
	Username string `mapstructure:"username"`
}

// New creates an instance of a MongoDB Connection with the provided connection information
func New(c Config) (Operator, error) {
	// TODO: implement a config sanitizer/validator
	if c.URI == "" {
		return Operator{}, ErrEmptyConnURI
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(c.URI))
	if err != nil {
		return Operator{}, fmt.Errorf("%s, %w", ErrInitClient.Error(), err)
	}

	if c.TimeoutMS <= 0 {
		c.TimeoutMS = _defaultTimeoutMS
	}

	ctx, done := context.WithTimeout(context.TODO(), time.Duration(c.TimeoutMS)*time.Millisecond)
	defer done()

	if err := client.Connect(ctx); err != nil {
		return Operator{}, fmt.Errorf("%s, %w", ErrFailedToConnect.Error(), err)
	}

	return Operator{
		client: client,
		config: c,
	}, nil
}

// Database is a getter for the database value set in the Operator's config
func (op *Operator) Database() string {
	return op.config.Database
}

// DefaultCollection is a getter for the default collection value set in the Operator's config
func (op *Operator) DefaultCollection() string {
	return op.config.DefaultCollection
}
