package datastore

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"go.mercari.io/datastore"
)

func TestRefreshRepository_Put(t *testing.T) {
	type (
		in struct {
			refresh *refresh
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
				refresh: &refresh{
					RefreshToken: "refresh",
					AccessToken:  "access",
				},
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
			mockDSClient.EXPECT().NameKey(KindRefresh, tt.in.refresh.RefreshToken, gomock.Nil()).Return(tt.returns.key)
			mockDSClient.EXPECT().Put(gomock.Any(), tt.returns.key, tt.in.refresh).Return(tt.returns.key, nil)

			storage := &refreshStorage{client: mockDSClient}
			err := storage.put(context.Background(), tt.in.refresh)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestRefreshRepository_Get(t *testing.T) {
	type (
		in struct {
			refreshToken string
		}

		out struct {
			refresh *refresh
		}

		returns struct {
			key     datastore.Key
			refresh *refresh
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
				refreshToken: "refresh",
			},
			out: out{
				refresh: &refresh{
					RefreshToken: "refresh",
					AccessToken:  "access",
				},
			},
			returns: returns{
				key: &mockKey{kind: "kind", name: "token"},
				refresh: &refresh{
					AccessToken: "access",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			mockDatastoreClient.EXPECT().NameKey(KindRefresh, tt.in.refreshToken, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Get(gomock.Any(), tt.returns.key, gomock.Any()).DoAndReturn(func(_ context.Context, _ datastore.Key, dst interface{}) error {
				val := reflect.ValueOf(dst)
				wVal := reflect.ValueOf(tt.returns.refresh)
				val.Elem().Set(wVal.Elem())
				return nil
			})

			storage := &refreshStorage{client: mockDatastoreClient}
			got, err := storage.get(context.Background(), tt.in.refreshToken)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.refresh, got) {
				t.Errorf("refresh data\nwant %#v\n got %#v", tt.out.refresh, got)
			}
		})
	}
}

func TestRefreshRepository_Delete(t *testing.T) {
	type (
		in struct {
			refreshToken string
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
				refreshToken: "sample",
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
			mockDatastoreClient.EXPECT().NameKey(KindRefresh, tt.in.refreshToken, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Delete(gomock.Any(), tt.returns.key).Return(nil)

			storage := &refreshStorage{client: mockDatastoreClient}
			if err := storage.delete(context.Background(), tt.in.refreshToken); err != nil {
				t.Error(err)
			}
		})
	}
}
