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
				t.Errorf("put key\n want %#v\n got %#v", wantPutKey, key)
			}
			if !reflect.DeepEqual(wantClient, src) {
				t.Errorf("put client\n want %#v\n got %#v", wantClient, src)
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
		t.Errorf("id want %v, got %v", wantID, id)
	}
}
