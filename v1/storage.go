// Package datastore is osin storage plugin for Google Cloud Datastore
package datastore

import (
	"context"
	"strings"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/aedatastore"
	"go.mercari.io/datastore/clouddatastore"

	"github.com/RangelReale/osin"
)

type (
	clientGetter interface {
		Get(ctx context.Context, id string) (*Client, error)
	}

	authDataHandler interface {
		put(ctx context.Context, auth *authorizeData) error
		get(ctx context.Context, code string) (*authorizeData, error)
		delete(ctx context.Context, code string) error
	}

	accessDataHandler interface {
		put(ctx context.Context, ac *accessData) error
		get(ctx context.Context, token string) (*accessData, error)
		delete(ctx context.Context, token string) error
	}

	refreshHandler interface {
		put(ctx context.Context, ref *refresh) error
		get(ctx context.Context, token string) (*refresh, error)
		delete(ctx context.Context, token string) error
	}
)

// Storage is handler to store OAuth2 tokens at GCP Datastore.
//  This struct must be live in only request scope.
//  So, instance should be created at like each "handler.ServeHTTP" methods.
type Storage struct {
	ctx               context.Context
	client            datastore.Client
	clientGetter      clientGetter
	authDataHandler   authDataHandler
	accessDataHandler accessDataHandler
	refreshHandler    refreshHandler
}

// NewStorage is constructor for storage of Google Cloud Datastore.
// The object created by this constructor uses Google Cloud Client Library for Go.
// If you want to use on Google App Engine Standard Edition, it should be recommanded to create object by NewStorageForGAE rather than use this.
func NewStorage(ctx context.Context, opts ...datastore.ClientOption) (*Storage, error) {
	client, err := clouddatastore.FromContext(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Storage{
		ctx:               ctx,
		client:            client,
		clientGetter:      newClientStorage(client),
		authDataHandler:   newAuthorizeDataRepository(client),
		accessDataHandler: newAccessDataRepository(client),
		refreshHandler:    newRefreshRepository(client),
	}, nil
}

// NewStorageForGAE is constructor for storage of Google Cloud Datastore.
// The object created by this constructor uses Google App Engine SDK for Go.
// If you want to use on other of Google App Engine Standard Edition, you must create object by NewStorage rather than use this.
func NewStorageForGAE(ctx context.Context, opts ...datastore.ClientOption) (*Storage, error) {
	client, err := aedatastore.FromContext(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &Storage{
		ctx:               ctx,
		client:            client,
		clientGetter:      newClientStorage(client),
		authDataHandler:   newAuthorizeDataRepository(client),
		accessDataHandler: newAccessDataRepository(client),
		refreshHandler:    newRefreshRepository(client),
	}, nil
}

// Clone is clonning storage instance
func (d *Storage) Clone() osin.Storage {
	// TODO investigate some error could be happen.
	return d
}

// Close releases resources used as datastore connections.
// This method must be call to finish use storage instance.
func (d *Storage) Close() {
	d.client.Close()
}

// GetClient loads client entity from datastore.
// If there is no match entity for the id, GetClient returns osin.ErrNotFound.
func (d *Storage) GetClient(id string) (osin.Client, error) {
	client, err := d.clientGetter.Get(d.ctx, id)
	if err != nil {
		return nil, errNoEntityOrDefault(err)
	}

	return client, nil
}

// SaveAuthorize stores authorize data entity to datastore.
func (d *Storage) SaveAuthorize(auth *osin.AuthorizeData) error {
	dauth, err := newAuthorizeDataFrom(auth)
	if err != nil {
		return err
	}

	return d.authDataHandler.put(d.ctx, dauth)
}

// LoadAuthorize loads authorize data entity with client entity from datastore.
// If there is no match entity for the id, LoadAuthorize returns osin.ErrNotFound.
func (d *Storage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	auth, err := d.authDataHandler.get(d.ctx, code)
	if err != nil {
		return nil, errNoEntityOrDefault(err)
	}

	client, err := d.GetClient(auth.ClientKey)
	if err != nil {
		return nil, err
	}

	return &osin.AuthorizeData{
		Code:                auth.Code,
		Client:              client,
		ExpiresIn:           int32(auth.ExpiresIn),
		Scope:               strings.Join(auth.Scope, " "),
		RedirectUri:         auth.RedirectURI,
		State:               auth.State,
		CreatedAt:           auth.CreatedAt,
		CodeChallenge:       auth.CodeChallenge,
		CodeChallengeMethod: auth.CodeChallengeMethod,
		UserData:            auth.UserData,
	}, nil
}

// RemoveAuthorize delete authorize data from datastore.
func (d *Storage) RemoveAuthorize(code string) error {
	return d.authDataHandler.delete(d.ctx, code)
}

// SaveAccess stores accesstoken entity to datastore.
func (d *Storage) SaveAccess(a *osin.AccessData) error {
	ad, err := newAccessDataFrom(a)
	if err != nil {
		return err
	}
	if err := d.accessDataHandler.put(d.ctx, ad); err != nil {
		return err
	}

	if a.RefreshToken != "" {
		return d.refreshHandler.put(d.ctx, newRefresh(a.RefreshToken, a.AccessToken))
	}

	return nil
}

// LoadAccess loads accesstoken data entity for access token with authorize data entity and client entity from datastore.
// If there is no match entity for the access token, LoadAuthorize returns osin.ErrNotFound.
func (d *Storage) LoadAccess(token string) (*osin.AccessData, error) {
	ad, err := d.accessDataHandler.get(d.ctx, token)
	if err != nil {
		return nil, errNoEntityOrDefault(err)
	}

	client, err := d.GetClient(ad.ClientKey)
	if err != nil {
		return nil, err
	}

	auth, err := d.LoadAuthorize(ad.AuthorizeCode)
	if err != nil {
		return nil, err
	}

	return &osin.AccessData{
		AccessToken:   ad.AccessToken,
		AuthorizeData: auth,
		Client:        client,
		RefreshToken:  ad.RefreshToken,
		ExpiresIn:     int32(ad.ExpiresIn),
		Scope:         strings.Join(ad.Scope, " "),
		RedirectUri:   ad.RedirectURI,
		CreatedAt:     ad.CreatedAt,
		UserData:      ad.UserData,
	}, nil
}

// RemoveAccess delete accesstoken data from datastore.
func (d *Storage) RemoveAccess(token string) error {
	return d.accessDataHandler.delete(d.ctx, token)
}

// LoadRefresh loads accesstoken data entity for refresh token with authorize data entity and client entity from datastore.
// If there is no match entity for the refresh token, LoadAuthorize returns osin.ErrNotFound.
func (d *Storage) LoadRefresh(token string) (*osin.AccessData, error) {
	ref, err := d.refreshHandler.get(d.ctx, token)
	if err != nil {
		return nil, errNoEntityOrDefault(err)
	}

	return d.LoadAccess(ref.AccessToken)
}

// RemoveRefresh delete refreshtoken data from datastore.
func (d *Storage) RemoveRefresh(token string) error {
	return d.refreshHandler.delete(d.ctx, token)
}

func errNoEntityOrDefault(err error) error {
	if err == datastore.ErrNoSuchEntity {
		return osin.ErrNotFound
	}
	return err
}
