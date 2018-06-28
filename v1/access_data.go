package datastore

import (
	"context"
	"strings"
	"time"

	"github.com/RangelReale/osin"
	"go.mercari.io/datastore"
)

// KindAccessData is datastore kind name of OAuth2 access token
const KindAccessData = "access_data"

type accessData struct {
	AccessToken       string `datastore:"-"`
	ParentAccessToken string
	ClientKey         string
	AuthorizeCode     string
	RefreshToken      string
	ExpiresIn         int64     `datastore:",noindex"`
	Scope             []string  `datastore:",noindex"`
	RedirectURI       string    `datastore:",noindex"`
	CreatedAt         time.Time `datastore:",noindex"`
	UserData          string    `datastore:",noindex"`
}

func newAccessDataFrom(a *osin.AccessData) (*accessData, error) {
	var (
		userData          string
		parentAccessToken string
	)
	if a.UserData != nil {
		ud, ok := a.UserData.(string)
		if !ok {
			return nil, ErrInvalidUserDataType
		}
		userData = ud
	}
	if a.AccessData != nil {
		parentAccessToken = a.AccessData.AccessToken
	}

	return &accessData{
		AccessToken:       a.AccessToken,
		ParentAccessToken: parentAccessToken,
		ClientKey:         a.Client.GetId(),
		AuthorizeCode:     a.AuthorizeData.Code,
		RefreshToken:      a.RefreshToken,
		ExpiresIn:         int64(a.ExpiresIn),
		Scope:             strings.Split(a.Scope, " "),
		RedirectURI:       a.RedirectUri,
		CreatedAt:         a.CreatedAt,
		UserData:          userData,
	}, nil
}

type accessDataRepository struct {
	client datastore.Client
}

func newAccessDataRepository(client datastore.Client) *accessDataRepository {
	return &accessDataRepository{client: client}
}

func (a *accessDataRepository) put(ctx context.Context, ac *accessData) error {
	key := a.client.NameKey(KindAccessData, ac.AccessToken, nil)
	_, err := a.client.Put(ctx, key, ac)
	return err
}

func (a *accessDataRepository) get(ctx context.Context, token string) (*accessData, error) {
	key := a.client.NameKey(KindAccessData, token, nil)
	access := new(accessData)
	if err := a.client.Get(ctx, key, access); err != nil {
		return nil, err
	}
	access.AccessToken = token
	return access, nil
}

func (a *accessDataRepository) delete(ctx context.Context, token string) error {
	key := a.client.NameKey(KindAccessData, token, nil)
	return a.client.Delete(ctx, key)
}
