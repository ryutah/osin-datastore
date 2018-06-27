// Package datastore is osin storage plugin for Google Cloud Datastore
package datastore

import (
	"context"
	"strings"

	"go.mercari.io/datastore"

	"github.com/RangelReale/osin"
)

type (
	clientHandler interface {
		put(ctx context.Context, c *Client) error
		get(ctx context.Context, id string) (*Client, error)
		delete(ctx context.Context, id string) error
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
	clientHandler     clientHandler
	authDataHandler   authDataHandler
	accessDataHandler accessDataHandler
	refreshHandler    refreshHandler
}

// NewStorage is constructor for DatastoreStorage
func NewStorage(ctx context.Context, opts ...StorageOption) (*Storage, error) {
	s := new(Storage)
	for _, o := range opts {
		if err := o(s); err != nil {
			return nil, err
		}
	}

	if s.client == nil {
		if err := WithCloudDatastoreClient(ctx)(s); err != nil {
			return nil, err
		}
	}

	s.clientHandler = newClientRepository(s.client)
	s.authDataHandler = newAuthorizeDataRepository(s.client)
	s.accessDataHandler = newAccessDataRepository(s.client)
	s.refreshHandler = newRefreshRepository(s.client)

	return s, nil
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
	client, err := d.clientHandler.get(d.ctx, id)
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
	_, err := d.refreshHandler.get(d.ctx, token)
	if err != nil {
		return nil, errNoEntityOrDefault(err)
	}

	return d.LoadAccess(token)
}

// RemoveRefresh delete refreshtoken data from datastore.
func (d *Storage) RemoveRefresh(token string) error {
	return d.refreshHandler.delete(d.ctx, token)
}

func (d *Storage) setClient(c datastore.Client) {
	d.client = c
}

func errNoEntityOrDefault(err error) error {
	if err == datastore.ErrNoSuchEntity {
		return osin.ErrNotFound
	}
	return err
}
