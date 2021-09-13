package mongodb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Count will return a count of the total number documents in the provided collection determinted by the bson filter.
//
// A value of -1 will be returned if any error occurs during this operation.
func (op Operator) Count(collection string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	if filter == nil {
		return -1, ErrNilFilter
	}

	if collection = strings.TrimSpace(collection); collection == "" {
		collection = op.config.DefaultCollection
	}

	c, err := op.getCollection(collection)
	if err != nil {
		return -1, fmt.Errorf("%s, %w", ErrCount.Error(), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.config.TimeoutMS)*time.Nanosecond)

	defer cancel()

	count, err := c.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return -1, fmt.Errorf("%s, %w", ErrCount.Error(), err)
	}

	return count, nil
}

// DeleteMany will delete multiple documents in the provided collection determined by the bson filter.
//
// A value of -1 will be returned if any error occurs during this operation.
func (op *Operator) DeleteMany(collection string, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	if filter == nil {
		return -1, ErrNilFilter
	}

	if collection = strings.TrimSpace(collection); collection == "" {
		collection = op.config.DefaultCollection
	}

	col, err := op.getCollection(collection)
	if err != nil {
		return -1, fmt.Errorf("%s; %w", ErrDeleteMany.Error(), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.config.TimeoutMS)*time.Nanosecond)

	defer cancel()

	res, err := col.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return -1, fmt.Errorf("%s, %w", ErrDeleteMany.Error(), err)
	}

	return res.DeletedCount, nil
}

// DeleteOne will delete at most one document from the provided collection using the provided bson filter.
//
// The deleted document will be decoded into the provided target, which should be a pointer.
func (op *Operator) DeleteOne(collection string, filter, target interface{}, opts ...*options.FindOneAndDeleteOptions) error {
	if filter == nil {
		return ErrNilFilter
	}

	if target == nil {
		return ErrNilTarget
	}

	if collection = strings.TrimSpace(collection); collection == "" {
		collection = op.config.DefaultCollection
	}

	col, err := op.getCollection(collection)
	if err != nil {
		return fmt.Errorf("%s; %w", ErrDeleteOne.Error(), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.config.TimeoutMS)*time.Nanosecond)

	defer cancel()

	res := col.FindOneAndDelete(ctx, filter, opts...)
	if res.Err() != nil && !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return fmt.Errorf("%s; %w", ErrDeleteOne, res.Err())
	}

	if err := res.Decode(target); err != nil {
		return fmt.Errorf("failed to decode deleted document to target interface; %w", err)
	}

	return nil
}

// FindMany will return a slice of documents from the provided collection determined by the bson filter.
//
// The target must be a non-nil pointer to a slice of desired type.
func (op Operator) FindMany(collection string, filter, target interface{}, opts ...*options.FindOptions) error {
	if filter == nil {
		return ErrNilFilter
	}

	if target == nil {
		return ErrNilTarget
	}

	if collection = strings.TrimSpace(collection); collection == "" {
		collection = op.config.DefaultCollection
	}

	c, err := op.getCollection(collection)
	if err != nil {
		return fmt.Errorf("%s, %w", ErrFindMany.Error(), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.config.TimeoutMS)*time.Nanosecond)

	defer cancel()

	cur, err := c.Find(ctx, filter, opts...)
	if err != nil {
		return fmt.Errorf("%s, %w", ErrFindMany.Error(), err)
	}

	if err := cur.All(context.Background(), target); err != nil {
		return fmt.Errorf("%s, %w", ErrFindMany.Error(), err)
	}

	return nil
}

// InsertMany will insert multiple documents into the provided collection.
func (op Operator) InsertMany(collection string, payload []interface{}, opts ...*options.InsertManyOptions) ([]interface{}, error) {
	if payload == nil {
		return nil, ErrNilPayload
	}

	if len(payload) <= 0 {
		return nil, fmt.Errorf("provided payload slice cannot be empty")
	}

	if collection = strings.TrimSpace(collection); collection == "" {
		collection = op.config.DefaultCollection
	}

	c, err := op.getCollection(collection)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection '%s'; %w", collection, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.config.TimeoutMS)*time.Nanosecond)
	defer cancel()

	res, err := c.InsertMany(ctx, payload, opts...)
	if err != nil {
		return nil, fmt.Errorf("%s; %w", ErrInsertMany.Error(), err)
	}

	return res.InsertedIDs, nil
}

// InsertOne will insert a single payload into the provided collection.
func (op Operator) InsertOne(collection string, payload interface{}, opts ...*options.InsertOneOptions) error {
	if payload == nil {
		return ErrNilPayload
	}

	if collection = strings.TrimSpace(collection); collection == "" {
		collection = op.config.DefaultCollection
	}

	c, err := op.getCollection(collection)
	if err != nil {
		return fmt.Errorf("%s, %w", ErrInsertOne.Error(), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.config.TimeoutMS)*time.Nanosecond)
	defer cancel()

	if _, err := c.InsertOne(ctx, payload, opts...); err != nil {
		return fmt.Errorf("%s, %w", ErrInsertOne.Error(), err)
	}

	return nil
}

// UpdateMany will perform an update on multiple documents in the provided collection using the update parameter in the payload.
//
// The count returned represents the total number of matched documents.
//
// A value of -1 will be returned if any error occurs during this operation.
func (op *Operator) UpdateMany(collection string, filter, payload interface{}, opts ...*options.UpdateOptions) (int64, error) {
	if filter == nil {
		return -1, ErrNilFilter
	}

	if payload == nil {
		return -1, ErrNilPayload
	}

	if collection = strings.TrimSpace(collection); collection == "" {
		collection = op.config.DefaultCollection
	}

	c, err := op.getCollection(collection)
	if err != nil {
		return -1, fmt.Errorf("%s, %w", ErrUpdateMany.Error(), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.config.TimeoutMS)*time.Nanosecond)

	defer cancel()

	res, err := c.UpdateMany(ctx, filter, bson.M{"$set": payload}, opts...)
	if err != nil {
		return -1, fmt.Errorf("%s, %w", ErrUpdateMany.Error(), err)

	}

	return res.MatchedCount, nil
}

// UpdateOne will update a single document in the provided collection using the update parameters defined in the payload.
func (op Operator) UpdateOne(collection string, filter, payload interface{}, opts ...*options.FindOneAndUpdateOptions) error {
	if filter == nil {
		return ErrNilFilter
	}

	if payload == nil {
		return ErrNilPayload
	}

	if collection = strings.TrimSpace(collection); collection == "" {
		collection = op.config.DefaultCollection
	}

	c, err := op.getCollection(collection)
	if err != nil {
		return fmt.Errorf("%s, %w", ErrUpdateOne.Error(), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(op.config.TimeoutMS)*time.Nanosecond)

	defer cancel()

	res := c.FindOneAndUpdate(ctx, filter, bson.M{"$set": payload}, opts...)
	if res.Err() != nil && errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return fmt.Errorf("%s; %w", ErrUpdateOne.Error(), err)
	}

	return nil
}

func (op Operator) getCollection(collection string) (*mongo.Collection, error) {
	if strings.TrimSpace(collection) == "" {
		return nil, ErrEmptyCollectionName
	}

	return op.client.Database(op.config.Database).Collection(collection), nil
}
