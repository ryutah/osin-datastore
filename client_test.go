package datastore

import (
	"context"
	"reflect"
	"testing"

	"go.mercari.io/datastore"
)

func TestClientRepository_Put_ValidClient(t *testing.T) {
	type want struct {
		key    datastore.Key
		client *Client
		id     string
	}
	tests := []struct {
		testName string
		in       *Client
		want     want
	}{
		{
			testName: "test1",
			in:       &Client{ID: "sample", Secret: "secret", RedirectUri: "redirect"},
			want: want{
				key:    &mockKey{kind: "test", name: "sample"},
				client: &Client{ID: "sample", Secret: "secret", RedirectUri: "redirect"},
				id:     "sample",
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
					if !reflect.DeepEqual(tt.want.client, src) {
						t.Errorf("put client\nwant %#v\n got %#v", tt.want.client, src)
					}
					return &mockKey{kind: "samplekind", name: "newsample"}, nil
				},
				nameKey: func(kind, name string, parent datastore.Key) datastore.Key {
					if kind != KindClient {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindClient, kind)
					}
					if name != tt.want.id {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.id, name)
					}
					return tt.want.key
				},
			}

			cr := &clientRepository{client: mockDSClient}
			err := cr.Put(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClientRepository_Get(t *testing.T) {
	type want struct {
		client  *Client
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
				client:  &Client{ID: "sample", Secret: "secret", RedirectUri: "redirect"},
				keyName: "sample",
				key:     &mockKey{kind: "kind", name: "sample"},
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
					wVal := reflect.ValueOf(&Client{Secret: "secret", RedirectUri: "redirect"})
					val.Elem().Set(wVal.Elem())
					return nil
				},
				nameKey: func(kind, name string, parent datastore.Key) datastore.Key {
					if kind != KindClient {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindClient, kind)
					}
					if name != tt.want.keyName {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.keyName, name)
					}
					return tt.want.key
				},
			}

			cr := &clientRepository{client: mockDatastoreClient}
			client, err := cr.Get(context.Background(), tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.want.client, client) {
				t.Errorf("client\nwant %#v\n got %#v", tt.want.client, client)
			}
		})
	}
}

func TestClientRepository_Delete(t *testing.T) {
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
					if kind != KindClient {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindClient, kind)
					}
					if name != tt.want.keyName {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.keyName, name)
					}
					return tt.want.key
				},
			}

			cr := &clientRepository{client: mockDatastoreClient}
			if err := cr.Delete(context.Background(), tt.in); err != nil {
				t.Error(err)
			}
		})
	}
}
