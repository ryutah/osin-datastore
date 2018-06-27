package datastore

import (
	"context"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/aedatastore"
	"go.mercari.io/datastore/clouddatastore"
)

// StorageOption is option of Storage
type StorageOption func(s *Storage) error

// WithAppEngineClient is constructor option of Storage.This should be set if application is running on GAE/SE.If you want to set some client option, you can config those by cliOpts.cliOpts are other library option. See below link.
//  https://godoc.org/go.mercari.io/datastore#ClientOption
func WithAppEngineClient(ctx context.Context, cliOpts ...datastore.ClientOption) StorageOption {
	return func(s *Storage) error {
		cli, err := aedatastore.FromContext(ctx, cliOpts...)
		if err != nil {
			return err
		}
		s.setClient(cli)
		return nil
	}
}

// WithCloudDatastoreClient is constructor option of Storage.This should be set if application is running on other of GAE/SE.If you want to set some client option, you can config those by cliOpts.cliOpts are other library option. See below link.
//  https://godoc.org/go.mercari.io/datastore#ClientOption
func WithCloudDatastoreClient(ctx context.Context, cliOpts ...datastore.ClientOption) StorageOption {
	return func(s *Storage) error {
		cli, err := clouddatastore.FromContext(ctx, cliOpts...)
		if err != nil {
			return err
		}
		s.setClient(cli)
		return nil
	}
}
