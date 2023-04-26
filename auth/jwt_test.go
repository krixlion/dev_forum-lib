package auth

import (
	"testing"
	"time"

	"github.com/krixlion/dev_forum-lib/internal/gentest"
	"github.com/lestrrat-go/jwx/jwt"
)

const (
	signedAccessToken  = "eyJhbGciOiJIUzI1NiIsImtpZCI6InRlc3QiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE2ODI1MTc1NzIsImlhdCI6MTY4MjUxNzI4NiwiaXNzIjoidGVzdCIsImp0aSI6InRlc3QiLCJzdWIiOiJ0ZXN0LWlkIiwidHlwZSI6ImFjY2Vzcy10b2tlbiJ9.wxoMBhYMLxZo_0il-EeQOnfcYUXfyuGWI--3IiYupbY"
	signedRefreshToken = "eyJhbGciOiJIUzI1NiIsImtpZCI6InRlc3QiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE2ODI1MTc1NzIsImlhdCI6MTY4MjUxNzI4NiwiaXNzIjoidGVzdCIsImp0aSI6InRlc3QiLCJzdWIiOiJ0ZXN0LWlkIiwidHlwZSI6InJlZnJlc2gtdG9rZW4ifQ.uiDFSRVO5urzRb5u4aXD4fn15hmNZN9w8ArDDdbLC5Q"
	testIssuer         = "test"
)

var (
	testHMACKey = []byte("key")

	testClockFunc = jwt.ClockFunc(func() time.Time {
		return time.Unix(1682517486, 0)
	})

	testKey = Key{
		Id:        "test",
		Algorithm: "HS256",
		Raw:       testHMACKey,
	}
)

func setUpTokenValidator(keys []Key, clockFunc jwt.Clock) TokenValidator {
	v, err := MakeTokenValidator(keys, Config{
		Issuer: testIssuer,
		Clock:  clockFunc,
	})
	if err != nil {
		panic(err)
	}
	return v
}

func TestTokenValidator_VerifyJWT(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name      string
		args      args
		keySet    []Key
		clockFunc jwt.Clock
		wantErr   bool
	}{
		{
			name: "Test if correctly parses a valid token",
			args: args{
				token: signedAccessToken,
			},
			keySet:    []Key{testKey},
			clockFunc: testClockFunc,
			wantErr:   false,
		},
		{
			name: "Test if fails on invalid token type",
			args: args{
				token: signedRefreshToken,
			},
			keySet:    []Key{testKey},
			clockFunc: testClockFunc,
			wantErr:   true,
		},
		{
			name: "Test if fails on invalid algorithm",
			args: args{
				token: signedRefreshToken,
			},
			keySet: []Key{{
				Id:        "test",
				Type:      "HMAC",
				Algorithm: "HS256",
				Raw:       testHMACKey,
			}},
			clockFunc: testClockFunc,
			wantErr:   true,
		},
		{
			name: "Test if fails on expired token",
			args: args{
				token: signedRefreshToken,
			},
			keySet: []Key{testKey},
			clockFunc: jwt.ClockFunc(func() time.Time {
				return time.Now().Add(time.Hour * 24)
			}),
			wantErr: true,
		},
		{
			name: "Test if fails on malformed token",
			args: args{
				token: gentest.RandomString(50),
			},
			keySet: []Key{testKey},
			clockFunc: jwt.ClockFunc(func() time.Time {
				return time.Now().Add(time.Hour * 24)
			}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := setUpTokenValidator(tt.keySet, tt.clockFunc)

			if err := v.VerifyJWT(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("TokenManager.VerifyJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
