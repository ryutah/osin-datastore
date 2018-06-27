package datastore

import (
	"context"
	"reflect"
	"testing"
	"time"

	"go.mercari.io/datastore"
)

func TestAccessDataRepository_Put_ValidClient(t *testing.T) {
	createAt := time.Now()
	type want struct {
		key  datastore.Key
		auth *accessData
		id   string
	}
	tests := []struct {
		testName string
		in       *accessData
		want     want
	}{
		{
			testName: "test1",
			in: &accessData{
				AccessToken:       "token",
				ParentAccessToken: "parentToken",
				ClientKey:         "client",
				AuthorizeCode:     "authCode",
				RefreshToken:      "refreshToken",
				ExpiresIn:         123,
				Scope:             []string{"scope1", "scope2"},
				RedirectURI:       "redirect",
				CreatedAt:         createAt,
				UserData:          "userData",
			},
			want: want{
				key: &mockKey{kind: "test", name: "sample"},
				auth: &accessData{
					AccessToken:       "token",
					ParentAccessToken: "parentToken",
					ClientKey:         "client",
					AuthorizeCode:     "authCode",
					RefreshToken:      "refreshToken",
					ExpiresIn:         123,
					Scope:             []string{"scope1", "scope2"},
					RedirectURI:       "redirect",
					CreatedAt:         createAt,
					UserData:          "userData",
				},
				id: "token",
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
					if kind != KindAccessData {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindAccessData, kind)
					}
					if name != tt.want.id {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.id, name)
					}
					return tt.want.key
				},
			}

			ar := &accessDataRepository{client: mockDSClient}
			err := ar.put(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAccessDataRepository_Get(t *testing.T) {
	type want struct {
		auth    *accessData
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
			in:       "token",
			want: want{
				keyName: "token",
				key:     &mockKey{kind: "kind", name: "token"},
				auth: &accessData{
					AccessToken:       "token",
					ParentAccessToken: "parentToken",
					ClientKey:         "client",
					AuthorizeCode:     "authCode",
					RefreshToken:      "refreshToken",
					ExpiresIn:         123,
					Scope:             []string{"scope1", "scope2"},
					RedirectURI:       "redirect",
					CreatedAt:         time.Now(),
					UserData:          "userData",
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
					if kind != KindAccessData {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindAccessData, kind)
					}
					if name != tt.want.keyName {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.keyName, name)
					}
					return tt.want.key
				},
			}

			ar := &accessDataRepository{client: mockDatastoreClient}
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

func TestAccessDataRepository_Get_Error(t *testing.T) {
	tests := []struct {
		testName string
		in       string
		want     error
	}{
		{
			testName: "test1",
			in:       "token",
			want:     datastore.ErrNoSuchEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			mockDatastoreClient := &mockDatastoreClient{
				get: func(context.Context, datastore.Key, interface{}) error {
					return datastore.ErrNoSuchEntity
				},
				nameKey: func(kind, name string, parent datastore.Key) datastore.Key {
					return new(mockKey)
				},
			}

			ar := &accessDataRepository{client: mockDatastoreClient}
			_, err := ar.get(context.Background(), tt.in)
			if err == nil {
				t.Fatalf("should be return error %T, but nil", tt.want)
			}
			if err != tt.want {
				t.Errorf("\nwant %#v\n got %#v", tt.want, err)
			}
		})
	}
}

func TestAccessDataRepository_Delete(t *testing.T) {
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
					if kind != KindAccessData {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindAccessData, kind)
					}
					if name != tt.want.keyName {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.keyName, name)
					}
					return tt.want.key
				},
			}

			ar := &accessDataRepository{client: mockDatastoreClient}
			if err := ar.delete(context.Background(), tt.in); err != nil {
				t.Error(err)
			}
		})
	}
}
