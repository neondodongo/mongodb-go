package mongodb

import (
	"errors"
	"fmt"
)

var (
	// ErrMongo is a base level error defining the operator that threw it.
	ErrMongo error = errors.New("mongodb operator error")
	// ErrEmptyCollectionName is thrown when an empty or whitespace collection name is provided.
	ErrEmptyCollectionName error = fmt.Errorf("[%w] empty or whitespace collection name provided", ErrMongo)
	//ErrEmptyConnURI is thrown when an empty or whitespace connection URI is provided.
	ErrEmptyConnURI error = fmt.Errorf("[%w] cannot initialize MongoDB driver client with empty or whitespace connection URI", ErrMongo)
	// ErrCount is thrown when a count operation fails.
	ErrCount error = fmt.Errorf("[%w] failed to count documents", ErrMongo)
	// ErrCursorDecode is thrown when the MongoDB driver cursor fails to decode its contents into a target interface.
	ErrCursorDecode error = fmt.Errorf("[%w] failed to decode cursor contents to target interface", ErrMongo)
	// ErrDeleteMany is thrown when a delete operation fails.
	ErrDeleteMany error = fmt.Errorf("[%w] failed to delete multiple documents", ErrMongo)
	// ErrDeleteOne is thrown when a delete operation fails.
	ErrDeleteOne error = fmt.Errorf("[%w] failed to delete document", ErrMongo)
	// ErrFailedToConnect is thrown when the Agent fails to establish a connection with MongoDB.
	ErrFailedToConnect error = fmt.Errorf("[%w] failed to connect to MongoDB instance", ErrMongo)
	// ErrFindOne is thrown when a find operation fails.
	ErrFindOne error = fmt.Errorf("[%w] failed to find document", ErrMongo)
	// ErrFindMany is thrown when a find operation fails.
	ErrFindMany error = fmt.Errorf("[%w] failed to find many documents", ErrMongo)
	// ErrInitClient is thrown when the Agent fails to initialize the underlying MongoDB driver client.
	ErrInitClient error = fmt.Errorf("[%w] failed to initialize MongoDB driver client", ErrMongo)
	// ErrInsertOne is thrown when an insert operation fails.
	ErrInsertOne error = fmt.Errorf("[%w] failed to insert document", ErrMongo)
	// ErrInsertMany is thrown when an insert many operation fails.
	ErrInsertMany error = fmt.Errorf("[%w] failed to insert many documents", ErrMongo)
	// ErrNilFilter is thrown when an operation receives a nil filter.
	ErrNilFilter error = fmt.Errorf("[%w] cannot perform operation with a nil filter", ErrMongo)
	// ErrNilTarget is thrown when an operation receives a nil target.
	ErrNilTarget error = fmt.Errorf("[%w] cannot perform operation with nil target", ErrMongo)
	// ErrNilPayload is thrown when an operation receives a nil payload.
	ErrNilPayload error = fmt.Errorf("[%w] cannot perform operation with a nil payload", ErrMongo)
	// ErrPing is thrown when a database ping fails.
	ErrPing error = fmt.Errorf("[%w] failed to ping MongoDB instance", ErrMongo)
	// ErrUpdateMany is thrown when an update one operation fails.
	ErrUpdateMany error = fmt.Errorf("[%w] failed to update multiple documents", ErrMongo)
	// ErrUpdateOne is thrown when an update one operation fails.
	ErrUpdateOne error = fmt.Errorf("[%w] failed to update document", ErrMongo)
)
