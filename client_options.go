package datastore

import (
	"context"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/aedatastore"
)

type datastoreClient interface {
	setClient(cli datastore.Client)
}

type datastoreClientOptionHandler func(datastoreClient) error

// WithAppEngineClient set datastore handler client to use appengine SDK datastore library.
func WithAppEngineClient(ctx context.Context) datastoreClientOptionHandler {
	return func(d datastoreClient) error {
		c, err := aedatastore.FromContext(ctx)
		if err != nil {
			return err
		}
		d.setClient(c)
		return nil
	}
}
