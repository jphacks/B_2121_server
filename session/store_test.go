package session_test

import (
	"testing"

	"github.com/jphacks/B_2121_server/session"
)

func TestStoreGet(t *testing.T) {
	t.Parallel()

	store := session.NewStore("store_key")
	token, err := store.New(123)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name    string
		token   string
		wantErr bool
		want    session.AuthInfo
	}{
		{
			name:  "Invalid token",
			token: "token",
			want: session.AuthInfo{
				Authenticated: false,
			},
		},
		{
			name:  "Valid token",
			token: token,
			want: session.AuthInfo{
				Authenticated: true,
				UserId:        123,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			got, err := store.Get(testCase.token)
			if (err != nil) != testCase.wantErr {
				t.Errorf("wantErr = %v, wantErr = %v", testCase.wantErr, err != nil)
			}
			if !testCase.wantErr && *got != testCase.want {
				t.Errorf("got = %v, want = %v", *got, testCase.want)
			}
		})
	}
}
