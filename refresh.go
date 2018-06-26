package datastore

import (
	"context"

	"go.mercari.io/datastore"
)

// KindRefresh is datastore kind name of OAuth2 refresh token
const KindRefresh = "refresh"

type refresh struct {
	RefreshToken string `datastore:"-"`
	AccessToken  string `datastore:",noindex"`
}

type refreshRepository struct {
	client datastore.Client
}

func (r *refreshRepository) put(ctx context.Context, ref *refresh) error {
	key := r.client.NameKey(KindRefresh, ref.RefreshToken, nil)
	_, err := r.client.Put(ctx, key, ref)
	return err
}

func (r *refreshRepository) get(ctx context.Context, token string) (*refresh, error) {
	key := r.client.NameKey(KindRefresh, token, nil)
	ref := new(refresh)
	if err := r.client.Get(ctx, key, ref); err != nil {
		return nil, err
	}
	ref.RefreshToken = token
	return ref, nil
}

func (r *refreshRepository) delete(ctx context.Context, token string) error {
	key := r.client.NameKey(KindRefresh, token, nil)
	return r.client.Delete(ctx, key)
}
