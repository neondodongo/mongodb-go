package mongodb

import (
	"context"
	"fmt"
	"time"
)

// Ping will send a request to the MongoDB instance returning an error if the connection is not active
func (op Operator) Ping() error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(op.config.TimeoutMS)*time.Millisecond)

	defer cancel()

	if err := op.client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("%s, %w", ErrPing.Error(), err)
	}

	return nil
}
