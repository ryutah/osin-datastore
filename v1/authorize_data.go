package datastore

import (
	"context"
	"strings"
	"time"

	"github.com/RangelReale/osin"
	"go.mercari.io/datastore"
)

// KindAuthorizeData is datastore kind name of OAuth2 authorize data stored
const KindAuthorizeData = "authorize_data"

type authorizeData struct {
	Code                string `datastore:"-"`
	ClientKey           string
	ExpiresIn           int64     `datastore:",noindex"`
	Scope               []string  `datastore:",noindex"`
	RedirectURI         string    `datastore:",noindex"`
	State               string    `datastore:",noindex"`
	CreatedAt           time.Time `datastore:",noindex"`
	UserData            string    `datastore:",noindex"`
	CodeChallenge       string    `datastore:",noindex"`
	CodeChallengeMethod string    `datastore:",noindex"`
}

func newAuthorizeDataFrom(a *osin.AuthorizeData) (*authorizeData, error) {
	var userData string
	if a.UserData != nil {
		ud, ok := a.UserData.(string)
		if !ok {
			return nil, ErrInvalidUserDataType
		}
		userData = ud
	}

	return &authorizeData{
		Code:                a.Code,
		ClientKey:           a.Client.GetId(),
		ExpiresIn:           int64(a.ExpiresIn),
		Scope:               strings.Split(a.Scope, " "),
		RedirectURI:         a.RedirectUri,
		State:               a.State,
		CreatedAt:           a.CreatedAt,
		CodeChallenge:       a.CodeChallenge,
		CodeChallengeMethod: a.CodeChallengeMethod,
		UserData:            userData,
	}, nil
}

type authorizeDataStorage struct {
	client datastore.Client
}

func newAuthorizeDataStorage(client datastore.Client) *authorizeDataStorage {
	return &authorizeDataStorage{client: client}
}

func (a *authorizeDataStorage) put(ctx context.Context, auth *authorizeData) error {
	key := a.client.NameKey(KindAuthorizeData, auth.Code, nil)
	_, err := a.client.Put(ctx, key, auth)
	return err
}

func (a *authorizeDataStorage) get(ctx context.Context, code string) (*authorizeData, error) {
	key := a.client.NameKey(KindAuthorizeData, code, nil)
	auth := new(authorizeData)
	if err := a.client.Get(ctx, key, auth); err != nil {
		return nil, err
	}
	auth.Code = code
	return auth, nil
}

func (a *authorizeDataStorage) delete(ctx context.Context, code string) error {
	key := a.client.NameKey(KindAuthorizeData, code, nil)
	return a.client.Delete(ctx, key)
}
