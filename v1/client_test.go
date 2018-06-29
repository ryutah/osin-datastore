package datastore

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"go.mercari.io/datastore"
)

func TestClientStorage_Put(t *testing.T) {
	type (
		in struct {
			client *Client
		}

		returns struct {
			key datastore.Key
		}
	)

	tests := []struct {
		testName string
		in       in
		returns  returns
	}{
		{
			testName: "test1",
			in: in{
				client: &Client{ID: "sample", Secret: "secret", RedirectUri: "redirect"},
			},
			returns: returns{
				key: &mockKey{kind: "test", name: "sample"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDSClient := NewMockClient(ctrl)
			mockDSClient.EXPECT().NameKey(KindClient, tt.in.client.ID, gomock.Nil()).Return(tt.returns.key)
			mockDSClient.EXPECT().Put(gomock.Any(), tt.returns.key, tt.in.client).Return(tt.returns.key, nil)

			cr := &ClientStorage{client: mockDSClient}
			err := cr.Put(context.Background(), tt.in.client)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClientStorage_Put_InValidClient(t *testing.T) {
	type (
		in struct {
			client *Client
		}

		out struct {
			err error
		}
	)
	tests := []struct {
		testName string
		in       in
		out      out
	}{
		{
			testName: "test1",
			in: in{
				client: &Client{Secret: "secret", RedirectUri: "redirect"},
			},
			out: out{
				err: ErrEmptyClientID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cr := &ClientStorage{}
			err := cr.Put(context.Background(), tt.in.client)
			if err != tt.out.err {
				t.Errorf("return error\nwant: %#v\n got: %#v", tt.out.err, err)
			}
		})
	}
}

func TestClientStorage_PutMulti(t *testing.T) {
	type (
		in struct {
			clients []*Client
		}

		returns struct {
			keys []datastore.Key
		}
	)
	tests := []struct {
		testName string
		in       in
		returns  returns
	}{
		{
			testName: "test1",
			in: in{
				clients: []*Client{
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
			for i, client := range tt.in.clients {
				mockDSClient.EXPECT().NameKey(KindClient, client.ID, gomock.Nil()).Return(tt.returns.keys[i])
			}
			mockDSClient.EXPECT().PutMulti(gomock.Any(), tt.returns.keys, tt.in.clients).Return(nil, nil)

			cr := &ClientStorage{client: mockDSClient}
			err := cr.PutMulti(context.Background(), tt.in.clients)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClientStorage_Get(t *testing.T) {
	type (
		in struct {
			id string
		}

		out struct {
			client *Client
		}

		returns struct {
			key    datastore.Key
			client *Client
		}
	)

	tests := []struct {
		testName string
		in       in
		out      out
		returns  returns
	}{
		{
			testName: "test1",
			in: in{
				id: "sample",
			},
			out: out{
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
			mockDatastoreClient.EXPECT().NameKey(KindClient, tt.in.id, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Get(gomock.Any(), tt.returns.key, gomock.Any()).DoAndReturn(func(_ context.Context, _ datastore.Key, dst interface{}) error {
				val := reflect.ValueOf(dst)
				wVal := reflect.ValueOf(tt.returns.client)
				val.Elem().Set(wVal.Elem())
				return nil
			})

			cr := &ClientStorage{client: mockDatastoreClient}
			got, err := cr.Get(context.Background(), tt.in.id)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.client, got) {
				t.Errorf("client\nwant %#v\n got %#v", tt.out.client, got)
			}
		})
	}
}

func TestClientStorage_GetMulti(t *testing.T) {
	type (
		in struct {
			ids []string
		}

		out struct {
			clients []*Client
		}

		returns struct {
			keys    []datastore.Key
			clients []*Client
		}
	)

	tests := []struct {
		testName string
		in       in
		out      out
		returns  returns
	}{
		{
			testName: "test1",
			in: in{
				ids: []string{"id1", "id2", "id3"},
			},
			out: out{
				clients: []*Client{
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
			for i, id := range tt.in.ids {
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
			gots, err := cr.GetMulti(context.Background(), tt.in.ids)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.clients, gots) {
				t.Errorf("client\nwant %+v\n got %+v", tt.out.clients, gots)
			}
		})
	}
}

func TestClientStorage_Delete(t *testing.T) {
	type (
		in struct {
			id string
		}

		returns struct {
			key datastore.Key
		}
	)

	tests := []struct {
		testName string
		in       in
		returns  returns
	}{
		{
			testName: "test1",
			in: in{
				id: "sample",
			},
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
			mockDatastoreClient.EXPECT().NameKey(KindClient, tt.in.id, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Delete(gomock.Any(), tt.returns.key).Return(nil)

			cr := &ClientStorage{client: mockDatastoreClient}
			if err := cr.Delete(context.Background(), tt.in.id); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestClientStorage_DeleteMulti(t *testing.T) {
	type (
		in struct {
			ids []string
		}

		returns struct {
			keys []datastore.Key
		}
	)
	tests := []struct {
		testName string
		in       in
		returns  returns
	}{
		{
			testName: "test1",
			in: in{
				ids: []string{"id1", "id2", "id3"},
			},
			returns: returns{
				keys: []datastore.Key{
					&mockKey{kind: KindClient, name: "id1"},
					&mockKey{kind: KindClient, name: "id2"},
					&mockKey{kind: KindClient, name: "id3"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			for i, id := range tt.in.ids {
				mockDatastoreClient.EXPECT().NameKey(KindClient, id, gomock.Nil()).Return(tt.returns.keys[i])
			}
			mockDatastoreClient.EXPECT().DeleteMulti(gomock.Any(), tt.returns.keys).Return(nil)

			cr := &ClientStorage{client: mockDatastoreClient}
			if err := cr.DeleteMulti(context.Background(), tt.in.ids); err != nil {
				t.Error(err)
			}
		})
	}
}
