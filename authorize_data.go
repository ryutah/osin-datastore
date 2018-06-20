package datastore

import (
	"context"
	"time"

	"go.mercari.io/datastore"
)

// KindAuthorizeData is datastore kind name of OAuth2 authorize data stored
const KindAuthorizeData = "authorize_data"

type authorizeData struct {
	Code                string `datastore:"-"`
	ClientKey           string
	ExpiresIn           int64     `datastore:",noindex"`
	Scope               string    `datastore:",noindex"`
	RedirectURI         string    `datastore:",noindex"`
	State               string    `datastore:",noindex"`
	CreatedAt           time.Time `datastore:",noindex"`
	UserData            string    `datastore:",noindex"`
	CodeChallenge       string    `datastore:",noindex"`
	CodeChallengeMethod string    `datastore:",noindex"`
}

type authorizeDataRepository struct {
	client datastore.Client
}

func (a *authorizeDataRepository) put(ctx context.Context, auth *authorizeData) error {
	key := a.client.NameKey(KindAuthorizeData, auth.Code, nil)
	_, err := a.client.Put(ctx, key, auth)
	return err
}

func (a *authorizeDataRepository) get(ctx context.Context, code string) (*authorizeData, error) {
	key := a.client.NameKey(KindAuthorizeData, code, nil)
	auth := new(authorizeData)
	if err := a.client.Get(ctx, key, auth); err != nil {
		return nil, err
	}
	auth.Code = code
	return auth, nil
}

func (a *authorizeDataRepository) delete(ctx context.Context, code string) error {
	key := a.client.NameKey(KindAuthorizeData, code, nil)
	return a.client.Delete(ctx, key)
}
