package datastore

import (
	"context"
	"reflect"
	"testing"

	"go.mercari.io/datastore"
)

func TestRefreshRepository_Put_ValidClient(t *testing.T) {
	type want struct {
		key datastore.Key
		ref *refresh
		id  string
	}
	tests := []struct {
		testName string
		in       *refresh
		want     want
	}{
		{
			testName: "test1",
			in: &refresh{
				RefreshToken: "refresh",
				AccessToken:  "access",
			},
			want: want{
				key: &mockKey{kind: "test", name: "sample"},
				ref: &refresh{
					RefreshToken: "refresh",
					AccessToken:  "access",
				},
				id: "refresh",
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
					if !reflect.DeepEqual(tt.want.ref, src) {
						t.Errorf("put refresh\nwant %#v\n got %#v", tt.want.ref, src)
					}
					return &mockKey{kind: "samplekind", name: "newsample"}, nil
				},
				nameKey: func(kind, name string, parent datastore.Key) datastore.Key {
					if kind != KindRefresh {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindRefresh, kind)
					}
					if name != tt.want.id {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.id, name)
					}
					return tt.want.key
				},
			}

			rr := &refreshRepository{client: mockDSClient}
			err := rr.put(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRefreshRepository_Get(t *testing.T) {
	type want struct {
		ref     *refresh
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
			in:       "refresh",
			want: want{
				keyName: "refresh",
				key:     &mockKey{kind: "kind", name: "token"},
				ref: &refresh{
					RefreshToken: "refresh",
					AccessToken:  "access",
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
					wVal := reflect.ValueOf(tt.want.ref)
					val.Elem().Set(wVal.Elem())
					return nil
				},
				nameKey: func(kind, name string, parent datastore.Key) datastore.Key {
					if kind != KindRefresh {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindRefresh, kind)
					}
					if name != tt.want.keyName {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.keyName, name)
					}
					return tt.want.key
				},
			}

			rr := &refreshRepository{client: mockDatastoreClient}
			ref, err := rr.get(context.Background(), tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.want.ref, ref) {
				t.Errorf("authorize data\nwant %#v\n got %#v", tt.want.ref, ref)
			}
		})
	}
}

func TestRefreshRepository_Delete(t *testing.T) {
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
					if kind != KindRefresh {
						t.Errorf("NameKey kind\nwant %v\n got %v", KindRefresh, kind)
					}
					if name != tt.want.keyName {
						t.Errorf("NameKey name\nwant %v\n got %v", tt.want.keyName, name)
					}
					return tt.want.key
				},
			}

			rr := &refreshRepository{client: mockDatastoreClient}
			if err := rr.delete(context.Background(), tt.in); err != nil {
				t.Error(err)
			}
		})
	}
}
