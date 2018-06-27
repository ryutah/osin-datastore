package datastore

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.mercari.io/datastore"
)

func TestAuthorizeDataRepository_Put_ValidClient(t *testing.T) {
	createAt := time.Now()
	type want struct {
		key  datastore.Key
		auth *authorizeData
		id   string
	}
	tests := []struct {
		testName string
		in       *authorizeData
		want     want
	}{
		{
			testName: "test1",
			in: &authorizeData{
				Code:                "code",
				ClientKey:           "clientKey",
				ExpiresIn:           123,
				Scope:               []string{"scope", "scope2"},
				RedirectURI:         "redirect",
				State:               "state",
				CreatedAt:           createAt,
				UserData:            "userdata",
				CodeChallenge:       "code_challenge",
				CodeChallengeMethod: "code_challenge_method",
			},
			want: want{
				key: &mockKey{kind: "test", name: "sample"},
				auth: &authorizeData{
					Code:                "code",
					ClientKey:           "clientKey",
					ExpiresIn:           123,
					Scope:               []string{"scope", "scope2"},
					RedirectURI:         "redirect",
					State:               "state",
					CreatedAt:           createAt,
					UserData:            "userdata",
					CodeChallenge:       "code_challenge",
					CodeChallengeMethod: "code_challenge_method",
				},
				id: "code",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockDSClient := &mockDatastoreClient{
				put: func(ctx context.Context, key datastore.Key, src interface{}) (datastore.Key, error) {
					if !reflect.DeepEqual(tt.want.key, key) {
						t.Errorf("put key\nwant %#v\n got %#v", tt.want.key, key)
					}
					if !reflect.DeepEqual(tt.want.auth, src) {
						t.Errorf("put auth\nwant %#v\n got %#v", tt.want.auth, src)
					}
					return &mockKey{kind: "samplekind", name: "newsample"}, nil
				},
				nameKey: func(kind, name string, parent datastore.Key) datastore.Key {
					if kind != KindAuthorizeData {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindAuthorizeData, kind)
					}
					if name != tt.want.id {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.id, name)
					}
					return tt.want.key
				},
			}

			ar := &authorizeDataRepository{client: mockDSClient}
			err := ar.put(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAuthorizeDataRepository_Get(t *testing.T) {
	type want struct {
		auth    *authorizeData
		keyName string
		key     datastore.Key
	}
	tests := []struct {
		testName string
		in       string
		want     want
	}{
		{
			testName: "test1",
			in:       "code",
			want: want{
				keyName: "code",
				key:     &mockKey{kind: "kind", name: "code"},
				auth: &authorizeData{
					Code:                "code",
					ClientKey:           "clientKey",
					ExpiresIn:           123,
					Scope:               []string{"scope", "scope2"},
					RedirectURI:         "redirect",
					State:               "state",
					CreatedAt:           time.Now(),
					UserData:            "userdata",
					CodeChallenge:       "code_challenge",
					CodeChallengeMethod: "code_challenge_method",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockDatastoreClient := &mockDatastoreClient{
				get: func(ctx context.Context, key datastore.Key, dst interface{}) error {
					if !reflect.DeepEqual(tt.want.key, key) {
						t.Errorf("get key\nwant %#v\n get %#v", tt.want.key, key)
					}

					val := reflect.ValueOf(dst)
					wVal := reflect.ValueOf(tt.want.auth)
					val.Elem().Set(wVal.Elem())
					return nil
				},
				nameKey: func(kind, name string, parent datastore.Key) datastore.Key {
					if kind != KindAuthorizeData {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindAuthorizeData, kind)
					}
					if name != tt.want.keyName {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.keyName, name)
					}
					return tt.want.key
				},
			}

			ar := &authorizeDataRepository{client: mockDatastoreClient}
			auth, err := ar.get(context.Background(), tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.want.auth, auth) {
				t.Errorf("authorize data\nwant %#v\n got %#v", tt.want.auth, auth)
			}
		})
	}
}

func TestAuthorizeDataRepository_Delete(t *testing.T) {
	type want struct {
		keyName string
		key     datastore.Key
	}
	tests := []struct {
		testName string
		in       string
		want     want
	}{
		{
			testName: "test1",
			in:       "sample",
			want: want{
				keyName: "sample",
				key:     &mockKey{kind: "kind", name: "sample"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			mockDatastoreClient := &mockDatastoreClient{
				delete: func(ctx context.Context, key datastore.Key) error {
					if !reflect.DeepEqual(key, tt.want.key) {
						t.Errorf("delete key\nwant %#v\n get %#v", tt.want.key, key)
					}
					return nil
				},
				nameKey: func(kind, name string, parent datastore.Key) datastore.Key {
					if kind != KindAuthorizeData {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindAuthorizeData, kind)
					}
					if name != tt.want.keyName {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.keyName, name)
					}
					return tt.want.key
				},
			}

			ar := &authorizeDataRepository{client: mockDatastoreClient}
			if err := ar.delete(context.Background(), tt.in); err != nil {
				t.Error(err)
			}
		})
	}
}
