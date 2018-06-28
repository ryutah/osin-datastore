package datastore

import (
	"context"
	"fmt"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/aedatastore"
	"go.mercari.io/datastore/clouddatastore"
)

// KindClient is datastore kind name of OAuth2 client stored
const KindClient = "client"

// Client is struct of OAuth2 client.
type Client struct {
	ID          string `json:"id,omitempty" datastore:"-"`
	Secret      string `json:"secret,omitempty" datastore:",noindex"`
	RedirectUri string `json:"redirect_uri,omitempty" datastore:",noindex"`
	UserData    string `json:"user_data,omitempty" datastore:",noindex"`
}

// GetId return client id.
func (c *Client) GetId() string {
	return c.ID
}

// GetSecret return client secret.
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetRedirectUri return redirect uri
func (c *Client) GetRedirectUri() string {
	return c.RedirectUri
}

// GetUserData return user data of client
func (c *Client) GetUserData() interface{} {
	return c.UserData
}

func (c Client) String() string {
	return fmt.Sprintf(
		"ID: %q, Secret: %q, RedirectURI: %q, UserData: %q",
		c.ID, c.Secret, c.RedirectUri, c.UserData,
	)
}

// ClientStorage is datastore handler for client.
type ClientStorage struct {
	client datastore.Client
}

func newClientStorage(client datastore.Client) *ClientStorage {
	return &ClientStorage{client: client}
}

// NewClientStorage create ClientStorage object.
// The object created by this constructor uses Google Cloud Client Library for Go.
// If you want to use on Google App Engine Standard Edition, it should be recommanded to create object by NewClientStorageForGAE rather than use this.
func NewClientStorage(ctx context.Context, opt ...datastore.ClientOption) (*ClientStorage, error) {
	client, err := clouddatastore.FromContext(ctx, opt...)
	if err != nil {
		return nil, err
	}
	return &ClientStorage{client: client}, nil
}

// NewClientStorageForGAE create ClientStorage object.
// The object created by this constructor uses Google App Engine SDK for Go.
// If you want to use on other of Google App Engine Standard Edition, you must create object by NewClientStorage rather than use this.
func NewClientStorageForGAE(ctx context.Context, opt ...datastore.ClientOption) (*ClientStorage, error) {
	client, err := aedatastore.FromContext(ctx, opt...)
	if err != nil {
		return nil, err
	}
	return &ClientStorage{client: client}, nil
}

// Put create or update client entity.
// The ID field of Client uses as Datastore's key.
func (cl *ClientStorage) Put(ctx context.Context, c *Client) error {
	if c.GetId() == "" {
		return ErrEmptyClientID
	}
	key := cl.client.NameKey(KindClient, c.GetId(), nil)
	_, err := cl.client.Put(ctx, key, c)
	return err
}

// PutMulti create or update multiple client entities.
// The ID field of Client uses as Datastore's key.
func (cl *ClientStorage) PutMulti(ctx context.Context, cs []*Client) error {
	keys := make([]datastore.Key, len(cs))
	for i, c := range cs {
		if c.GetId() == "" {
			return ErrEmptyClientID
		}
		keys[i] = cl.client.NameKey(KindClient, c.GetId(), nil)
	}
	_, err := cl.client.PutMulti(ctx, keys, cs)
	return err
}

// Get search client for given id.
func (cl *ClientStorage) Get(ctx context.Context, id string) (*Client, error) {
	key := cl.client.NameKey(KindClient, id, nil)
	dst := new(Client)
	if err := cl.client.Get(ctx, key, dst); err != nil {
		return nil, err
	}
	dst.ID = id
	return dst, nil
}

// GetMulti search multiple client for given ids.
func (cl *ClientStorage) GetMulti(ctx context.Context, ids []string) ([]*Client, error) {
	keys := make([]datastore.Key, len(ids))
	for i, id := range ids {
		keys[i] = cl.client.NameKey(KindClient, id, nil)
	}

	clients := make([]*Client, len(keys))
	if err := cl.client.GetMulti(ctx, keys, clients); err != nil {
		return nil, err
	}
	for i, id := range ids {
		clients[i].ID = id
	}
	return clients, nil
}

// Delete removes client entitye for id from Datastore.
func (cl *ClientStorage) Delete(ctx context.Context, id string) error {
	key := cl.client.NameKey(KindClient, id, nil)
	return cl.client.Delete(ctx, key)
}

// DeleteMulti removes multiple clients entitye for ids from Datastore.
func (cl *ClientStorage) DeleteMulti(ctx context.Context, ids []string) error {
	keys := make([]datastore.Key, len(ids))
	for i, id := range ids {
		keys[i] = cl.client.NameKey(KindClient, id, nil)
	}
	return cl.client.DeleteMulti(ctx, keys)
}
