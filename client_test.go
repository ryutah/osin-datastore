package datastore

import (
	"context"
	"reflect"
	"testing"

	"go.mercari.io/datastore"
)

func TestClientRepository_Put_ValidClient(t *testing.T) {
	var (
		wantPutKey = &mockKey{kind: "test", id: 111}
		wantClient = &Client{Secret: "secret", RedirectUri: "redirect"}
		wantID     = "123"
	)

	mockDSClient := &mockDatastoreClient{
		put: func(ctx context.Context, key datastore.Key, src interface{}) (datastore.Key, error) {
			if !reflect.DeepEqual(wantPutKey, key) {
				t.Errorf("put key\nwant %#v\n got %#v", wantPutKey, key)
			}
			if !reflect.DeepEqual(wantClient, src) {
				t.Errorf("put client\nwant %#v\n got %#v", wantClient, src)
			}
			return &mockKey{kind: "samplekind", id: 123}, nil
		},
		incompleteKey: func(kind string, parent datastore.Key) datastore.Key {
			return wantPutKey
		},
	}

	cr := &clientRepository{client: mockDSClient}
	id, err := cr.Put(context.Background(), &Client{
		Secret:      "secret",
		RedirectUri: "redirect",
	})
	if err != nil {
		t.Fatal(err)
	}
	if id != wantID {
		t.Errorf("id\nwant %v, got %v", wantID, id)
	}
}

func TestClientRepository_Get(t *testing.T) {
	var (
		wantClient = &Client{
			ID:          1,
			Secret:      "secret",
			RedirectUri: "redirect",
		}
		wantKeyID  int64 = 1
		wantGetKey       = &mockKey{kind: "kind", id: 1}
	)

	mockDatastoreClient := &mockDatastoreClient{
		get: func(ctx context.Context, key datastore.Key, dst interface{}) error {
			if !reflect.DeepEqual(wantGetKey, key) {
				t.Errorf("get key\nwant %#v\n get %#v", wantGetKey, key)
			}

			val := reflect.ValueOf(dst)
			wVal := reflect.ValueOf(&Client{Secret: "secret", RedirectUri: "redirect"})
			val.Elem().Set(wVal.Elem())
			return nil
		},
		idKey: func(kind string, id int64, parent datastore.Key) datastore.Key {
			if id != wantKeyID {
				t.Errorf("key id\nwant %v, got %v", wantKeyID, id)
			}
			return wantGetKey
		},
	}

	cr := &clientRepository{client: mockDatastoreClient}
	client, err := cr.Get(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(wantClient, client) {
		t.Errorf("client\nwant %#v\n got %#v", wantClient, client)
	}
}

func TestClientRepository_Delete(t *testing.T) {
	var (
		wantKeyID     int64 = 123
		wantDeleteKey       = &mockKey{kind: "kind", id: 123}
	)

	mockDatastoreClient := &mockDatastoreClient{
		delete: func(ctx context.Context, key datastore.Key) error {
			return nil
		},
		idKey: func(kind string, id int64, parent datastore.Key) datastore.Key {
			if id != wantKeyID {
				t.Errorf("key id\nwant %v, got %v", wantKeyID, id)
			}
			return wantDeleteKey
		},
	}

	cr := &clientRepository{client: mockDatastoreClient}
	if err := cr.Delete(context.Background(), "123"); err != nil {
		t.Error(err)
	}
}
