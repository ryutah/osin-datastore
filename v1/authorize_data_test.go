package datastore

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"go.mercari.io/datastore"
)

func TestAuthorizeDataStorage_Put(t *testing.T) {
	type (
		in struct {
			authorize *authorizeData
		}

		returns struct {
			key datastore.Key
		}
	)

	createAt := time.Now()
	tests := []struct {
		testName string
		in       in
		returns  returns
	}{
		{
			testName: "test1",
			in: in{
				authorize: &authorizeData{
					Code:                "code",
					ClientKey:           "clientKey",
					ExpiresIn:           123,
					Scope:               []string{"scope", "scope2"},
					RedirectURI:         "redirect",
					State:               "state",
					CreatedAt:           createAt,
					UserData:            "userdata",
					CodeChallenge:       "code_challenge",
					CodeChallengeMethod: "code_challenge_method",
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
			mockDSClient.EXPECT().NameKey(KindAuthorizeData, tt.in.authorize.Code, gomock.Nil()).Return(tt.returns.key)
			mockDSClient.EXPECT().Put(gomock.Any(), tt.returns.key, tt.in.authorize).Return(tt.returns.key, nil)

			storage := &authorizeDataStorage{client: mockDSClient}
			err := storage.put(context.Background(), tt.in.authorize)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAuthorizeDataStorage_Get(t *testing.T) {
	type (
		in struct {
			code string
		}

		out struct {
			authorize *authorizeData
		}

		returns struct {
			key       datastore.Key
			authorize *authorizeData
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
				code: "code",
			},
			out: out{
				authorize: &authorizeData{
					Code:                "code",
					ClientKey:           "clientKey",
					ExpiresIn:           123,
					Scope:               []string{"scope", "scope2"},
					RedirectURI:         "redirect",
					State:               "state",
					CreatedAt:           createdAt,
					UserData:            "userdata",
					CodeChallenge:       "code_challenge",
					CodeChallengeMethod: "code_challenge_method",
				},
			},
			returns: returns{
				key: &mockKey{kind: "kind", name: "code"},
				authorize: &authorizeData{
					ClientKey:           "clientKey",
					ExpiresIn:           123,
					Scope:               []string{"scope", "scope2"},
					RedirectURI:         "redirect",
					State:               "state",
					CreatedAt:           createdAt,
					UserData:            "userdata",
					CodeChallenge:       "code_challenge",
					CodeChallengeMethod: "code_challenge_method",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDatastoreClient := NewMockClient(ctrl)
			mockDatastoreClient.EXPECT().NameKey(KindAuthorizeData, tt.in.code, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Get(gomock.Any(), tt.returns.key, gomock.Any()).DoAndReturn(func(_ context.Context, _ datastore.Key, dst interface{}) error {
				val := reflect.ValueOf(dst)
				wVal := reflect.ValueOf(tt.returns.authorize)
				val.Elem().Set(wVal.Elem())
				return nil
			})

			storage := &authorizeDataStorage{client: mockDatastoreClient}
			got, err := storage.get(context.Background(), tt.in.code)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.authorize, got) {
				t.Errorf("authorize data\nwant %#v\n got %#v", tt.out.authorize, got)
			}
		})
	}
}

func TestAuthorizeDataStorage_Delete(t *testing.T) {
	type (
		in struct {
			code string
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
				code: "sample",
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
			mockDatastoreClient.EXPECT().NameKey(KindAuthorizeData, tt.in.code, gomock.Nil()).Return(tt.returns.key)
			mockDatastoreClient.EXPECT().Delete(gomock.Any(), tt.returns.key).Return(nil)

			storage := &authorizeDataStorage{client: mockDatastoreClient}
			if err := storage.delete(context.Background(), tt.in.code); err != nil {
				t.Error(err)
			}
		})
	}
}
