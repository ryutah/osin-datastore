package datastore

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"go.mercari.io/datastore"
)

func TestAccessDataStorage_Put(t *testing.T) {
	createAt := time.Now()
	type (
		in struct {
			accessData *accessData
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
				&accessData{
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
			mockDSClient.EXPECT().NameKey(KindAccessData, tt.in.accessData.AccessToken, gomock.Nil()).Return(tt.returns.key)
			mockDSClient.EXPECT().Put(gomock.Any(), tt.returns.key, tt.in.accessData).Return(tt.returns.key, nil)

			storage := &accessDataStorage{client: mockDSClient}
			err := storage.put(context.Background(), tt.in.accessData)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAccessDataStorage_Get(t *testing.T) {
	type (
		in struct {
			token string
		}

		out struct {
			accessData *accessData
		}

		returns struct {
			key        datastore.Key
			accessData *accessData
		}
	)

	createdAt := time.Now()
	tests := []struct {
		testName string
		in       in
		out      out
		returns  returns
	}{
		{
			testName: "test1",
			in: in{
				token: "token",
			},
			out: out{
				accessData: &accessData{
					AccessToken:       "token",
					ParentAccessToken: "parentToken",
					ClientKey:         "client",
					AuthorizeCode:     "authCode",
					RefreshToken:      "refreshToken",
					ExpiresIn:         123,
					Scope:             []string{"scope1", "scope2"},
					RedirectURI:       "redirect",
					CreatedAt:         createdAt,
					UserData:          "userData",
				},
			},
			returns: returns{
				key: &mockKey{kind: "kind", name: "token"},
				accessData: &accessData{
					ParentAccessToken: "parentToken",
					ClientKey:         "client",
					AuthorizeCode:     "authCode",
					RefreshToken:      "refreshToken",
					ExpiresIn:         123,
					Scope:             []string{"scope1", "scope2"},
					RedirectURI:       "redirect",
					CreatedAt:         createdAt,
					UserData:          "userData",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			mockDatastoreClient.EXPECT().NameKey(KindAccessData, tt.in.token, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Get(gomock.Any(), tt.returns.key, gomock.Any()).DoAndReturn(func(_ context.Context, _ datastore.Key, dst interface{}) error {
				val := reflect.ValueOf(dst)
				wVal := reflect.ValueOf(tt.returns.accessData)
				val.Elem().Set(wVal.Elem())
				return nil
			})

			storage := &accessDataStorage{client: mockDatastoreClient}
			got, err := storage.get(context.Background(), tt.in.token)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.accessData, got) {
				t.Errorf("access data\nwant %#v\n got %#v", tt.out.accessData, got)
			}
		})
	}
}

func TestAccessDataStorage_Get_Error(t *testing.T) {
	type (
		in struct {
			token string
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
				token: "token",
			},
			out: out{
				err: datastore.ErrNoSuchEntity,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			mockDatastoreClient.EXPECT().NameKey(gomock.Any(), gomock.Any(), gomock.Any()).Return(new(mockKey))
			mockDatastoreClient.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(datastore.ErrNoSuchEntity)

			storage := &accessDataStorage{client: mockDatastoreClient}
			_, err := storage.get(context.Background(), tt.in.token)
			if err != tt.out.err {
				t.Errorf("\nwant %#v\n got %#v", tt.out.err, err)
			}
		})
	}
}

func TestAccessDataStorage_Delete(t *testing.T) {
	type (
		in struct {
			token string
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
				token: "sample",
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
			mockDatastoreClient.EXPECT().NameKey(KindAccessData, tt.in.token, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Delete(gomock.Any(), tt.returns.key).Return(nil)

			storage := &accessDataStorage{client: mockDatastoreClient}
			if err := storage.delete(context.Background(), tt.in.token); err != nil {
				t.Error(err)
			}
		})
	}
}
