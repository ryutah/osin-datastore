package datastore

import (
	"reflect"
	"testing"
	"time"

	"github.com/RangelReale/osin"
	"github.com/golang/mock/gomock"
)

func TestStorage_GetClient(t *testing.T) {
	tests := []struct {
		testName string
		in       string
		want     *Client
	}{
		{
			testName: "test1",
			in:       "client",
			want: &Client{
				ID:          "client",
				Secret:      "secret",
				RedirectUri: "redirect",
				UserData:    "sample",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mch := NewMockclientHandler(ctrl)
			mch.EXPECT().get(gomock.Any(), tt.in).Times(1).Return(tt.want, nil)

			storage := &Storage{clientHandler: mch}

			got, err := storage.GetClient(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestStorage_SaveAuthorize(t *testing.T) {
	createdAt := time.Now()

	tests := []struct {
		testName string
		in       *osin.AuthorizeData
		want     *authorizeData
	}{
		{
			testName: "test1",
			in: &osin.AuthorizeData{
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
			want: &authorizeData{
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
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mauh := NewMockauthDataHandler(ctrl)
			mauh.EXPECT().put(gomock.Any(), tt.want).Return(nil)

			storage := &Storage{authDataHandler: mauh}

			if err := storage.SaveAuthorize(tt.in); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_LoadAuthorize(t *testing.T) {
	type returns struct {
		auth   *authorizeData
		client *Client
	}
	tests := []struct {
		testName string
		in       string
		want     *osin.AuthorizeData
		returns  returns
	}{
		{
			testName: "test1",
			in:       "code",
			want: &osin.AuthorizeData{
				Code:     "auth",
				Client:   &Client{ID: "client"},
				UserData: "",
			},
			returns: returns{
				auth:   &authorizeData{Code: "auth", ClientKey: "client"},
				client: &Client{ID: "client"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mauh = NewMockauthDataHandler(ctrl)
				mch  = NewMockclientHandler(ctrl)
			)
			mauh.EXPECT().get(gomock.Any(), tt.in).Return(tt.returns.auth, nil)
			mch.EXPECT().get(gomock.Any(), tt.returns.auth.ClientKey).Return(tt.returns.client, nil)

			storage := &Storage{
				authDataHandler: mauh,
				clientHandler:   mch,
			}

			got, err := storage.LoadAuthorize(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("\nwant: %#v\n got: %#v", tt.want, got)
			}
		})
	}
}

func TestStorage_RemoveAuthorize(t *testing.T) {
	tests := []struct {
		testName string
		in       string
	}{
		{
			testName: "test1",
			in:       "code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mauh := NewMockauthDataHandler(ctrl)
			mauh.EXPECT().delete(gomock.Any(), tt.in).Return(nil)

			storage := &Storage{authDataHandler: mauh}

			if err := storage.RemoveAuthorize(tt.in); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_SaveAccess(t *testing.T) {
	createdAt := time.Now()
	tests := []struct {
		testName string
		in       *osin.AccessData
		want     *accessData
	}{
		{
			testName: "test1",
			in: &osin.AccessData{
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
			want: &accessData{
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
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mach := NewMockaccessDataHandler(ctrl)
			mach.EXPECT().put(gomock.Any(), tt.want).Return(nil)

			storage := &Storage{accessDataHandler: mach}

			if err := storage.SaveAccess(tt.in); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_SaveAccess_WithRefreshToken(t *testing.T) {
	tests := []struct {
		testName string
		in       *osin.AccessData
		want     *refresh
	}{
		{
			testName: "test1",
			in: &osin.AccessData{
				AccessToken:   "token",
				RefreshToken:  "refresh",
				Client:        new(Client),
				AuthorizeData: new(osin.AuthorizeData),
			},
			want: &refresh{
				RefreshToken: "refresh",
				AccessToken:  "token",
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
			mrh.EXPECT().put(gomock.Any(), tt.want).Return(nil)

			storage := &Storage{
				accessDataHandler: mach,
				refreshHandler:    mrh,
			}

			if err := storage.SaveAccess(tt.in); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_LoadAccess(t *testing.T) {
	type returns struct {
		access *accessData
		auth   *authorizeData
		client *Client
	}
	tests := []struct {
		testName string
		in       string
		want     *osin.AccessData
		returns  returns
	}{
		{
			testName: "test1",
			in:       "token",
			want: &osin.AccessData{
				Client: &Client{ID: "client"},
				AuthorizeData: &osin.AuthorizeData{
					Code:     "auth",
					Client:   &Client{ID: "client"},
					UserData: "",
				},
				AccessToken: "token",
				UserData:    "",
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
				mch  = NewMockclientHandler(ctrl)
				mauh = NewMockauthDataHandler(ctrl)
			)
			mach.EXPECT().get(gomock.Any(), tt.in).Return(tt.returns.access, nil)
			mch.EXPECT().get(gomock.Any(), tt.returns.access.ClientKey).Return(tt.returns.client, nil)
			mauh.EXPECT().get(gomock.Any(), tt.returns.access.AuthorizeCode).Return(tt.returns.auth, nil)
			mch.EXPECT().get(gomock.Any(), tt.returns.auth.ClientKey).Return(tt.returns.client, nil)

			storage := &Storage{
				accessDataHandler: mach,
				clientHandler:     mch,
				authDataHandler:   mauh,
			}

			got, err := storage.LoadAccess(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("\nwant: %#v\n got: %#v", tt.want, got)
				if !reflect.DeepEqual(tt.want.Client, got.Client) {
					t.Errorf("client\nwant: %#v\n got: %#v", tt.want.Client, got.Client)
				}
				if !reflect.DeepEqual(tt.want.AuthorizeData, got.AuthorizeData) {
					t.Errorf("auth\nwant: %#v\n got: %#v", tt.want.AuthorizeData, got.AuthorizeData)
				}
				if !reflect.DeepEqual(tt.want.AuthorizeData.Client, got.AuthorizeData.Client) {
					t.Errorf("auth client\nwant: %#v\n got: %#v", tt.want.AuthorizeData.Client, got.AuthorizeData.Client)
				}
			}
		})
	}
}

func TestStorage_RemoveAccess(t *testing.T) {
	tests := []struct {
		testName string
		in       string
	}{
		{
			testName: "test1",
			in:       "token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mach := NewMockaccessDataHandler(ctrl)
			mach.EXPECT().delete(gomock.Any(), tt.in).Return(nil)

			storage := &Storage{accessDataHandler: mach}

			if err := storage.RemoveAccess(tt.in); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestStorage_LoadRefresh(t *testing.T) {
	type want struct {
		refresh *refresh
		access  *osin.AccessData
	}
	type returns struct {
		access *accessData
		auth   *authorizeData
		client *Client
	}

	createdAt := time.Now()
	tests := []struct {
		testName string
		in       string
		want     want
		returns  returns
	}{
		{
			testName: "test1",
			in:       "refresh_token",
			want: want{
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
				refresh: &refresh{
					RefreshToken: "refresh_token",
					AccessToken:  "token",
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
				auth: &authorizeData{
					Code:      "auth",
					ClientKey: "a_client",
				},
				client: &Client{ID: "client"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// TODO implement test codes.
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			var (
				mrh  = NewMockrefreshHandler(ctrl)
				mach = NewMockaccessDataHandler(ctrl)
				mch  = NewMockclientHandler(ctrl)
				mauh = NewMockauthDataHandler(ctrl)
			)
			mrh.EXPECT().get(gomock.Any(), tt.in).Return(tt.want.refresh, nil)
			mach.EXPECT().get(gomock.Any(), tt.want.refresh.AccessToken).Return(tt.returns.access, nil)
			mch.EXPECT().get(gomock.Any(), tt.returns.access.ClientKey).Return(tt.returns.client, nil)
			mauh.EXPECT().get(gomock.Any(), tt.returns.access.AuthorizeCode).Return(tt.returns.auth, nil)
			mch.EXPECT().get(gomock.Any(), tt.returns.auth.ClientKey).Return(tt.returns.client, nil)

			storage := &Storage{
				refreshHandler:    mrh,
				accessDataHandler: mach,
				clientHandler:     mch,
				authDataHandler:   mauh,
			}

			got, err := storage.LoadRefresh(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(tt.want.access, got) {
				t.Errorf("\nwant: %#v\n got: %#v", tt.want, got)
				if !reflect.DeepEqual(tt.want.access.Client, got.Client) {
					t.Errorf("client\nwant: %#v\n got: %#v", tt.want.access.Client, got.Client)
				}
				if !reflect.DeepEqual(tt.want.access.AuthorizeData, got.AuthorizeData) {
					t.Errorf("auth\nwant: %#v\n got: %#v", tt.want.access.AuthorizeData, got.AuthorizeData)
				}
				if !reflect.DeepEqual(tt.want.access.AuthorizeData.Client, got.AuthorizeData.Client) {
					t.Errorf("auth client\nwant: %#v\n got: %#v", tt.want.access.AuthorizeData.Client, got.AuthorizeData.Client)
				}
			}
		})
	}
}

func TestStorage_RemoveRefresh(t *testing.T) {
	tests := []struct {
		testName string
		in       string
	}{
		{
			testName: "test1",
			in:       "refresth_token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mrh := NewMockrefreshHandler(ctrl)
			mrh.EXPECT().delete(gomock.Any(), tt.in).Return(nil)

			storage := &Storage{refreshHandler: mrh}

			if err := storage.RemoveRefresh(tt.in); err != nil {
				t.Fatal(err)
			}
		})
	}
}
