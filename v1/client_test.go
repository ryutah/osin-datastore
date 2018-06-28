package datastore

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"go.mercari.io/datastore"
)

func TestClientRepository_Put_ValidClient(t *testing.T) {
	type want struct {
		id string
	}

	type returns struct {
		key    datastore.Key
		newKey datastore.Key
	}

	tests := []struct {
		testName string
		in       *Client
		want     want
		returns  returns
	}{
		{
			testName: "test1",
			in:       &Client{ID: "sample", Secret: "secret", RedirectUri: "redirect"},
			want: want{
				id: "sample",
			},
			returns: returns{
				key:    &mockKey{kind: "test", name: "sample"},
				newKey: &mockKey{kind: "samplekind", name: "newsample"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDSClient := NewMockClient(ctrl)
			mockDSClient.EXPECT().NameKey(KindClient, tt.want.id, gomock.Nil()).Return(tt.returns.key)
			mockDSClient.EXPECT().Put(gomock.Any(), tt.returns.key, tt.in).Return(tt.returns.newKey, nil)

			cr := &ClientStorage{client: mockDSClient}
			err := cr.Put(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClientRepository_Get(t *testing.T) {
	type want struct {
		client *Client
	}

	type returns struct {
		key    datastore.Key
		client *Client
	}

	tests := []struct {
		testName string
		in       string
		want     want
		returns  returns
	}{
		{
			testName: "test1",
			in:       "sample",
			want: want{
				client: &Client{ID: "sample", Secret: "secret", RedirectUri: "redirect"},
			},
			returns: returns{
				key:    &mockKey{kind: KindClient, name: "sample"},
				client: &Client{Secret: "secret", RedirectUri: "redirect"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			mockDatastoreClient.EXPECT().NameKey(KindClient, tt.in, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Get(gomock.Any(), tt.returns.key, gomock.Any()).DoAndReturn(func(_ context.Context, _ datastore.Key, dst interface{}) error {
				val := reflect.ValueOf(dst)
				wVal := reflect.ValueOf(tt.returns.client)
				val.Elem().Set(wVal.Elem())
				return nil
			})

			cr := &ClientStorage{client: mockDatastoreClient}
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
	type returns struct {
		key datastore.Key
	}

	tests := []struct {
		testName string
		in       string
		returns  returns
	}{
		{
			testName: "test1",
			in:       "sample",
			returns: returns{
				key: &mockKey{kind: "kind", name: "sample"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			mockDatastoreClient.EXPECT().NameKey(KindClient, tt.in, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Delete(gomock.Any(), tt.returns.key).Return(nil)

			cr := &ClientStorage{client: mockDatastoreClient}
			if err := cr.Delete(context.Background(), tt.in); err != nil {
				t.Error(err)
			}
		})
	}
}
