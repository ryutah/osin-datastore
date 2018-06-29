package datastore

import (
	"reflect"
	"testing"
	"time"

	"github.com/RangelReale/osin"
	"github.com/golang/mock/gomock"
)

func TestStorage_GetClient(t *testing.T) {
	type (
		in struct {
			id string
		}

		out struct {
			client *Client
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
				id: "client",
			},
			out: out{
				client: &Client{
					ID:          "client",
					Secret:      "secret",
					RedirectUri: "redirect",
					UserData:    "sample",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mch := NewMockclientGetter(ctrl)
			mch.EXPECT().Get(gomock.Any(), tt.in.id).Return(tt.out.client, nil)

			storage := &Storage{clientGetter: mch}

			got, err := storage.GetClient(tt.in.id)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.client, got) {
				t.Errorf("\nwant: %#v\n got: %#v", tt.out.client, got)
			}
		})
	}
}

func TestStorage_SaveAuthorize(t *testing.T) {
	type (
		in struct {
			authorize *osin.AuthorizeData
		}

		args struct {
			authorize *authorizeData
		}
	)
	createdAt := time.Now()

	tests := []struct {
		testName string
		in       in
		args     args
	}{
		{
			testName: "test1",
			in: in{
				authorize: &osin.AuthorizeData{
					Code:                "code",
					Client:              &Client{ID: "client"},
					ExpiresIn:           1,
					Scope:               "scope1 scope2",
					RedirectUri:         "redirect",
					State:               "state",
					CreatedAt:           createdAt,
					CodeChallenge:       "code_challenge",
					CodeChallengeMethod: "code_challenge_method",
					UserData:            "user_data",
				},
			},
			args: args{
				authorize: &authorizeData{
					Code:                "code",
					ClientKey:           "client",
					ExpiresIn:           1,
					Scope:               []string{"scope1", "scope2"},
					RedirectURI:         "redirect",
					State:               "state",
					CreatedAt:           createdAt,
					CodeChallenge:       "code_challenge",
					CodeChallengeMethod: "code_challenge_method",
					UserData:            "user_data",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mauh := NewMockauthDataHandler(ctrl)
			mauh.EXPECT().put(gomock.Any(), tt.args.authorize).Return(nil)

			storage := &Storage{authDataHandler: mauh}

			if err := storage.SaveAuthorize(tt.in.authorize); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_LoadAuthorize(t *testing.T) {
	type (
		in struct {
			code string
		}

		returns struct {
			authorize *authorizeData
			client    *Client
		}

		out struct {
			authorize *osin.AuthorizeData
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
				code: "code",
			},
			out: out{
				authorize: &osin.AuthorizeData{
					Code:     "auth",
					Client:   &Client{ID: "client"},
					UserData: "",
				},
			},
			returns: returns{
				authorize: &authorizeData{Code: "auth", ClientKey: "client"},
				client:    &Client{ID: "client"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mauh = NewMockauthDataHandler(ctrl)
				mch  = NewMockclientGetter(ctrl)
			)
			mauh.EXPECT().get(gomock.Any(), tt.in.code).Return(tt.returns.authorize, nil)
			mch.EXPECT().Get(gomock.Any(), tt.returns.authorize.ClientKey).Return(tt.returns.client, nil)

			storage := &Storage{
				authDataHandler: mauh,
				clientGetter:    mch,
			}

			got, err := storage.LoadAuthorize(tt.in.code)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.authorize, got) {
				t.Errorf("\nwant: %#v\n got: %#v", tt.out.authorize, got)
			}
		})
	}
}

func TestStorage_RemoveAuthorize(t *testing.T) {
	type (
		in struct {
			code string
		}
	)
	tests := []struct {
		testName string
		in       in
	}{
		{
			testName: "test1",
			in: in{
				code: "code",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mauh := NewMockauthDataHandler(ctrl)
			mauh.EXPECT().delete(gomock.Any(), tt.in.code).Return(nil)

			storage := &Storage{authDataHandler: mauh}

			if err := storage.RemoveAuthorize(tt.in.code); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_SaveAccess(t *testing.T) {
	type (
		in struct {
			access *osin.AccessData
		}

		args struct {
			access *accessData
		}
	)
	createdAt := time.Now()
	tests := []struct {
		testName string
		in       in
		args     args
	}{
		{
			testName: "test1",
			in: in{
				access: &osin.AccessData{
					AccessToken:   "token",
					AccessData:    &osin.AccessData{AccessToken: "token2"},
					AuthorizeData: &osin.AuthorizeData{Code: "code"},
					Client:        &Client{ID: "client"},
					RefreshToken:  "",
					ExpiresIn:     1,
					Scope:         "scope1 scope2",
					RedirectUri:   "redirect",
					CreatedAt:     createdAt,
					UserData:      "user_data",
				},
			},
			args: args{
				access: &accessData{
					AccessToken:       "token",
					ParentAccessToken: "token2",
					ClientKey:         "client",
					AuthorizeCode:     "code",
					RefreshToken:      "",
					ExpiresIn:         1,
					Scope:             []string{"scope1", "scope2"},
					RedirectURI:       "redirect",
					CreatedAt:         createdAt,
					UserData:          "user_data",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mach := NewMockaccessDataHandler(ctrl)
			mach.EXPECT().put(gomock.Any(), tt.args.access).Return(nil)

			storage := &Storage{accessDataHandler: mach}

			if err := storage.SaveAccess(tt.in.access); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_SaveAccess_WithRefreshToken(t *testing.T) {
	type (
		in struct {
			access *osin.AccessData
		}

		args struct {
			refresh *refresh
		}
	)

	tests := []struct {
		testName string
		in       in
		args     args
	}{
		{
			testName: "test1",
			in: in{
				access: &osin.AccessData{
					AccessToken:   "token",
					RefreshToken:  "refresh",
					Client:        new(Client),
					AuthorizeData: new(osin.AuthorizeData),
				},
			},
			args: args{
				refresh: &refresh{
					RefreshToken: "refresh",
					AccessToken:  "token",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mach = NewMockaccessDataHandler(ctrl)
				mrh  = NewMockrefreshHandler(ctrl)
			)
			mach.EXPECT().put(gomock.Any(), gomock.Any()).Return(nil)
			mrh.EXPECT().put(gomock.Any(), tt.args.refresh).Return(nil)

			storage := &Storage{
				accessDataHandler: mach,
				refreshHandler:    mrh,
			}

			if err := storage.SaveAccess(tt.in.access); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_LoadAccess(t *testing.T) {
	type (
		in struct {
			token string
		}

		returns struct {
			access *accessData
			auth   *authorizeData
			client *Client
		}

		out struct {
			access *osin.AccessData
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
				token: "token",
			},
			out: out{
				access: &osin.AccessData{
					Client: &Client{ID: "client"},
					AuthorizeData: &osin.AuthorizeData{
						Code:     "auth",
						Client:   &Client{ID: "client"},
						UserData: "",
					},
					AccessToken: "token",
					UserData:    "",
				},
			},
			returns: returns{
				access: &accessData{
					AccessToken:   "token",
					ClientKey:     "client",
					AuthorizeCode: "auth",
				},
				auth:   &authorizeData{Code: "auth", ClientKey: "a_client"},
				client: &Client{ID: "client"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			var (
				mach = NewMockaccessDataHandler(ctrl)
				mch  = NewMockclientGetter(ctrl)
				mauh = NewMockauthDataHandler(ctrl)
			)
			mach.EXPECT().get(gomock.Any(), tt.in.token).Return(tt.returns.access, nil)
			mch.EXPECT().Get(gomock.Any(), tt.returns.access.ClientKey).Return(tt.returns.client, nil)
			mauh.EXPECT().get(gomock.Any(), tt.returns.access.AuthorizeCode).Return(tt.returns.auth, nil)
			mch.EXPECT().Get(gomock.Any(), tt.returns.auth.ClientKey).Return(tt.returns.client, nil)

			storage := &Storage{
				accessDataHandler: mach,
				clientGetter:      mch,
				authDataHandler:   mauh,
			}

			got, err := storage.LoadAccess(tt.in.token)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.access, got) {
				t.Errorf("\nwant: %#v\n got: %#v", tt.out.access, got)
				if !reflect.DeepEqual(tt.out.access.Client, got.Client) {
					t.Errorf("client\nwant: %#v\n got: %#v", tt.out.access.Client, got.Client)
				}
				if !reflect.DeepEqual(tt.out.access.AuthorizeData, got.AuthorizeData) {
					t.Errorf("auth\nwant: %#v\n got: %#v", tt.out.access.AuthorizeData, got.AuthorizeData)
				}
				if !reflect.DeepEqual(tt.out.access.AuthorizeData.Client, got.AuthorizeData.Client) {
					t.Errorf("auth client\nwant: %#v\n got: %#v", tt.out.access.AuthorizeData.Client, got.AuthorizeData.Client)
				}
			}
		})
	}
}

func TestStorage_RemoveAccess(t *testing.T) {
	type (
		in struct {
			token string
		}
	)

	tests := []struct {
		testName string
		in       in
	}{
		{
			testName: "test1",
			in: in{
				token: "token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mach := NewMockaccessDataHandler(ctrl)
			mach.EXPECT().delete(gomock.Any(), tt.in.token).Return(nil)

			storage := &Storage{accessDataHandler: mach}

			if err := storage.RemoveAccess(tt.in.token); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_LoadRefresh(t *testing.T) {
	type (
		in struct {
			refreshToken string
		}

		out struct {
			access *osin.AccessData
		}

		returns struct {
			refresh      *refresh
			access       *accessData
			authorize    *authorizeData
			accessClient *Client
			authClient   *Client
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
				refreshToken: "refresh_token",
			},
			out: out{
				access: &osin.AccessData{
					AccessToken: "token",
					AuthorizeData: &osin.AuthorizeData{
						Code:     "auth",
						Client:   &Client{ID: "client"},
						UserData: "",
					},
					Client:       &Client{ID: "client"},
					RefreshToken: "refresh_token",
					ExpiresIn:    1,
					Scope:        "scope1 scope2",
					RedirectUri:  "redirect",
					CreatedAt:    createdAt,
					UserData:     "user_data",
				},
			},
			returns: returns{
				access: &accessData{
					AccessToken:   "token",
					AuthorizeCode: "auth",
					ClientKey:     "client",
					RefreshToken:  "refresh_token",
					ExpiresIn:     1,
					Scope:         []string{"scope1", "scope2"},
					RedirectURI:   "redirect",
					CreatedAt:     createdAt,
					UserData:      "user_data",
				},
				refresh: &refresh{
					RefreshToken: "refresh_token",
					AccessToken:  "token",
				},
				authorize: &authorizeData{
					Code:      "auth",
					ClientKey: "a_client",
				},
				accessClient: &Client{ID: "client"},
				authClient:   &Client{ID: "client"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mrh  = NewMockrefreshHandler(ctrl)
				mach = NewMockaccessDataHandler(ctrl)
				mch  = NewMockclientGetter(ctrl)
				mauh = NewMockauthDataHandler(ctrl)
			)
			mrh.EXPECT().get(gomock.Any(), tt.in.refreshToken).Return(tt.returns.refresh, nil)
			mach.EXPECT().get(gomock.Any(), tt.returns.refresh.AccessToken).Return(tt.returns.access, nil)
			mch.EXPECT().Get(gomock.Any(), tt.returns.access.ClientKey).Return(tt.returns.accessClient, nil)
			mauh.EXPECT().get(gomock.Any(), tt.returns.access.AuthorizeCode).Return(tt.returns.authorize, nil)
			mch.EXPECT().Get(gomock.Any(), tt.returns.authorize.ClientKey).Return(tt.returns.authClient, nil)

			storage := &Storage{
				refreshHandler:    mrh,
				accessDataHandler: mach,
				clientGetter:      mch,
				authDataHandler:   mauh,
			}

			got, err := storage.LoadRefresh(tt.in.refreshToken)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.out.access, got) {
				t.Errorf("\nwant: %#v\n got: %#v", tt.out.access, got)
				if !reflect.DeepEqual(tt.out.access.Client, got.Client) {
					t.Errorf("client\nwant: %#v\n got: %#v", tt.out.access.Client, got.Client)
				}
				if !reflect.DeepEqual(tt.out.access.AuthorizeData, got.AuthorizeData) {
					t.Errorf("auth\nwant: %#v\n got: %#v", tt.out.access.AuthorizeData, got.AuthorizeData)
				}
				if !reflect.DeepEqual(tt.out.access.AuthorizeData.Client, got.AuthorizeData.Client) {
					t.Errorf("auth client\nwant: %#v\n got: %#v", tt.out.access.AuthorizeData.Client, got.AuthorizeData.Client)
				}
			}
		})
	}
}

func TestStorage_RemoveRefresh(t *testing.T) {
	type (
		in struct {
			refreshToken string
		}
	)

	tests := []struct {
		testName string
		in       in
	}{
		{
			testName: "test1",
			in: in{
				refreshToken: "refresth_token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mrh := NewMockrefreshHandler(ctrl)
			mrh.EXPECT().delete(gomock.Any(), tt.in.refreshToken).Return(nil)

			storage := &Storage{refreshHandler: mrh}

			if err := storage.RemoveRefresh(tt.in.refreshToken); err != nil {
				t.Fatal(err)
			}
		})
	}
}
