package datastore

import (
	"context"

	"go.mercari.io/datastore"
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

func (c *Client) GetId() string {
	return c.ID
}

func (c *Client) GetSecret() string {
	return c.Secret
}

func (c *Client) GetRedirectUri() string {
	return c.RedirectUri
}

func (c *Client) GetUserData() interface{} {
	return c.UserData
}

type clientRepository struct {
	client datastore.Client
}

func newClientRepository(client datastore.Client) *clientRepository {
	return &clientRepository{client: client}
}

func (cl *clientRepository) Put(ctx context.Context, c *Client) error {
	key := cl.client.NameKey(KindClient, c.GetId(), nil)
	_, err := cl.client.Put(ctx, key, c)
	return err
}

func (cl *clientRepository) Get(ctx context.Context, id string) (*Client, error) {
	key := cl.client.NameKey(KindClient, id, nil)
	dst := new(Client)
	if err := cl.client.Get(ctx, key, dst); err != nil {
		return nil, err
	}
	dst.ID = id
	return dst, nil
}

func (cl *clientRepository) Delete(ctx context.Context, id string) error {
	key := cl.client.NameKey(KindClient, id, nil)
	return cl.client.Delete(ctx, key)
}
