package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Controller acts as a wrapper of the MongoDB driver exposing the MongoDB operator functionality.
type Controller interface {
	Count(collection string, filter interface{}, opts ...*options.CountOptions) (int64, error)
	DeleteMany(collection string, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	DeleteOne(collection string, filter, target interface{}, opts ...*options.FindOneAndDeleteOptions) error
	FindMany(collection string, filter, target interface{}, opts ...*options.FindOptions) error
	FindOne(colletion string, filter, target interface{}, opts ...*options.FindOneOptions) error
	InsertMany(collection string, payload []interface{}, opts ...*options.InsertManyOptions) ([]interface{}, error)
	InsertOne(collection string, payload interface{}, opts ...*options.InsertOneOptions) error
	UpdateMany(collection string, filter, payload interface{}, opts ...*options.UpdateOptions) (int64, error)
	UpdateOne(collection string, filter, payload interface{}, opts ...*options.FindOneAndUpdateOptions) error
	Ping() error
}

// Operator implements the Controller interface to perform MongoDB operations.
type Operator struct {
	client *mongo.Client
	config Config
}

// New creates an instance of a MongoDB Connection with the provided connection information.
func New(c Config) (*Operator, error) {
	if err := c.sanitizeAndValidate(); err != nil {
		return nil, fmt.Errorf("config failed validation; %w", err)
	}

	// TODO: Configure additional auth mechanisms; current default is SCRAM
	client, err := mongo.NewClient(
		options.Client().ApplyURI(c.URI),
		options.Client().SetAuth(options.Credential{
			Username: c.Username,
			Password: c.Password,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", ErrInitClient.Error(), err)
	}

	ctx, done := context.WithTimeout(context.TODO(), time.Duration(c.TimeoutMS)*time.Millisecond)
	defer done()

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("%s, %w", ErrFailedToConnect.Error(), err)
	}

	return &Operator{
		client: client,
		config: c,
	}, nil
}
