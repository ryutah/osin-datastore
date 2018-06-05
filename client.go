package datastore

import (
	"context"
	"strconv"

	"go.mercari.io/datastore"
)

// KindClient is datastore kind name of OAuth2 client stored
const KindClient = "client"

// Client is struct of OAuth2 client.
type Client struct {
	ID          int64  `json:"id" datastore:"-"`
	Secret      string `json:"secret" datastore:",noindex"`
	RedirectUri string `json:"redirect_uri" datastore:",noindex"`
}

func (c *Client) GetId() string {
	return strconv.FormatInt(c.ID, 10)
}

func (c *Client) GetSecret() string {
	return c.Secret
}

func (c *Client) GetRedirectUri() string {
	return c.RedirectUri
}

func (c *Client) GetUserData() interface{} {
	panic("not implemented")
}

type clientRepository struct {
	client datastore.Client
}

func (cl *clientRepository) Put(ctx context.Context, c *Client) (id string, err error) {
	key := cl.client.IncompleteKey(KindClient, nil)
	newKey, err := cl.client.Put(ctx, key, c)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(newKey.ID(), 10), nil
}

func (cl *clientRepository) Get(ctx context.Context, id string) (*Client, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	key := cl.client.IDKey(KindClient, intID, nil)
	dst := new(Client)
	if err := cl.client.Get(ctx, key, dst); err != nil {
		return nil, err
	}
	dst.ID = intID
	return dst, nil
}
