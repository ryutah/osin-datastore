package datastore

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"go.mercari.io/datastore"
)

func TestClientStorage_Put(t *testing.T) {
	type returns struct {
		key    datastore.Key
		newKey datastore.Key
	}

	tests := []struct {
		testName string
		in       *Client
		returns  returns
	}{
		{
			testName: "test1",
			in:       &Client{ID: "sample", Secret: "secret", RedirectUri: "redirect"},
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
			mockDSClient.EXPECT().NameKey(KindClient, tt.in.ID, gomock.Nil()).Return(tt.returns.key)
			mockDSClient.EXPECT().Put(gomock.Any(), tt.returns.key, tt.in).Return(tt.returns.newKey, nil)

			cr := &ClientStorage{client: mockDSClient}
			err := cr.Put(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClientStorage_Put_InValidClient(t *testing.T) {
	tests := []struct {
		testName string
		in       *Client
		want     error
	}{
		{
			testName: "test1",
			in:       &Client{Secret: "secret", RedirectUri: "redirect"},
			want:     ErrEmptyClientID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cr := &ClientStorage{}
			err := cr.Put(context.Background(), tt.in)
			if err != tt.want {
				t.Errorf("return error\nwant: %#v\n got: %#v", tt.want, err)
			}
		})
	}
}

func TestClientStorage_PutMulti(t *testing.T) {
	type returns struct {
		keys []datastore.Key
	}

	tests := []struct {
		testName string
		in       []*Client
		returns  returns
	}{
		{
			testName: "test1",
			in: []*Client{
				&Client{
					ID:          "client1",
					RedirectUri: "redirect1",
					Secret:      "secret1",
					UserData:    "user_data1",
				},
				&Client{
					ID:          "client2",
					RedirectUri: "redirect2",
					Secret:      "secret2",
					UserData:    "user_data2",
				},
				&Client{
					ID:          "client3",
					RedirectUri: "redirect3",
					Secret:      "secret3",
					UserData:    "user_data3",
				},
			},
			returns: returns{
				keys: []datastore.Key{
					&mockKey{kind: KindClient, name: "key1"},
					&mockKey{kind: KindClient, name: "key2"},
					&mockKey{kind: KindClient, name: "key3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDSClient := NewMockClient(ctrl)
			for i, client := range tt.in {
				mockDSClient.EXPECT().NameKey(KindClient, client.ID, gomock.Nil()).Return(tt.returns.keys[i])
			}
			mockDSClient.EXPECT().PutMulti(gomock.Any(), tt.returns.keys, tt.in).Return(nil, nil)

			cr := &ClientStorage{client: mockDSClient}
			err := cr.PutMulti(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClientStorage_Get(t *testing.T) {
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

func TestClientStorage_GetMulti(t *testing.T) {
	type returns struct {
		keys    []datastore.Key
		clients []*Client
	}

	tests := []struct {
		testName string
		in       []string
		wants    []*Client
		returns  returns
	}{
		{
			testName: "test1",
			in:       []string{"id1", "id2", "id3"},
			wants: []*Client{
				&Client{
					ID:          "id1",
					Secret:      "secret1",
					RedirectUri: "redirect1",
					UserData:    "user_data1",
				},
				&Client{
					ID:          "id2",
					Secret:      "secret2",
					RedirectUri: "redirect2",
					UserData:    "user_data2",
				},
				&Client{
					ID:          "id3",
					Secret:      "secret3",
					RedirectUri: "redirect3",
					UserData:    "user_data3",
				},
			},
			returns: returns{
				keys: []datastore.Key{
					&mockKey{kind: KindClient, name: "id1"},
					&mockKey{kind: KindClient, name: "id2"},
					&mockKey{kind: KindClient, name: "id3"},
				},
				clients: []*Client{
					&Client{
						Secret:      "secret1",
						RedirectUri: "redirect1",
						UserData:    "user_data1",
					},
					&Client{
						Secret:      "secret2",
						RedirectUri: "redirect2",
						UserData:    "user_data2",
					},
					&Client{
						Secret:      "secret3",
						RedirectUri: "redirect3",
						UserData:    "user_data3",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			for i, id := range tt.in {
				mockDatastoreClient.EXPECT().NameKey(KindClient, id, gomock.Nil()).Return(tt.returns.keys[i])
			}
			mockDatastoreClient.EXPECT().GetMulti(gomock.Any(), tt.returns.keys, gomock.Any()).DoAndReturn(func(_ context.Context, _ []datastore.Key, dst interface{}) error {
				dval := reflect.ValueOf(dst)
				for i := 0; i < dval.Len(); i++ {
					dval.Index(i).Set(reflect.ValueOf(tt.returns.clients[i]))
				}
				return nil
			})

			cr := &ClientStorage{client: mockDatastoreClient}
			gots, err := cr.GetMulti(context.Background(), tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.wants, gots) {
				t.Errorf("client\nwant %+v\n got %+v", tt.wants, gots)
			}
		})
	}
}

func TestClientStorage_Delete(t *testing.T) {
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

func TestClientStorage_DeleteMulti(t *testing.T) {
	tests := []struct {
		testName string
		in       []string
		returns  []datastore.Key
	}{
		{
			testName: "test1",
			in:       []string{"id1", "id2", "id3"},
			returns: []datastore.Key{
				&mockKey{kind: KindClient, name: "id1"},
				&mockKey{kind: KindClient, name: "id2"},
				&mockKey{kind: KindClient, name: "id3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			for i, id := range tt.in {
				mockDatastoreClient.EXPECT().NameKey(KindClient, id, gomock.Nil()).Return(tt.returns[i])
			}
			mockDatastoreClient.EXPECT().DeleteMulti(gomock.Any(), tt.returns).Return(nil)

			cr := &ClientStorage{client: mockDatastoreClient}
			if err := cr.DeleteMulti(context.Background(), tt.in); err != nil {
				t.Error(err)
			}
		})
	}
}
