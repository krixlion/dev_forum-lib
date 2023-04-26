package auth

import (
	"errors"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

var ErrInvalidTokenType = errors.New("invalid token type")

type Key struct {
	Id        string
	Algorithm string
	Type      string
	Raw       interface{}
}

type Config struct {
	Issuer string
	Clock  jwt.Clock
}

type TokenValidator struct {
	config Config
	keySet jwk.Set
}

// MakeTokenValidator returns a new instance or a non-nil error if provided keys
// fail to serialize into keyset.
// If no Clock is provided in the config time.Now() is used by default.
func MakeTokenValidator(keys []Key, config Config) (TokenValidator, error) {
	keySet := jwk.NewSet()

	for _, key := range keys {
		jwKey, err := jwk.New(key.Raw)
		if err != nil {
			return TokenValidator{}, err
		}

		if err := jwKey.Set(jwk.KeyIDKey, key.Id); err != nil {
			return TokenValidator{}, err
		}

		if err := jwKey.Set(jwk.KeyTypeKey, key.Type); err != nil {
			return TokenValidator{}, err
		}

		if err := jwKey.Set(jwk.AlgorithmKey, jwa.HS256); err != nil {
			return TokenValidator{}, err
		}

		keySet.Add(jwKey)
	}

	if config.Clock == nil {
		config.Clock = jwt.ClockFunc(time.Now)
	}

	return TokenValidator{
		config: config,
		keySet: keySet,
	}, nil
}

// VerifyJWT returns a non-nil error if the token is expired,
// signature is invalid or any of the token's claims are different than expected.
// Eg. token was issued in the future or specified 'kid' does not exist.
func (validator TokenValidator) VerifyJWT(token string) error {
	jwToken, err := jwt.ParseString(token, jwt.WithKeySet(validator.keySet))
	if err != nil {
		return err
	}

	validateOptions := []jwt.ValidateOption{
		jwt.WithIssuer(validator.config.Issuer),
		jwt.WithClock(validator.config.Clock),
	}

	if err := jwt.Validate(jwToken, validateOptions...); err != nil {
		return err
	}

	tokenType, ok := jwToken.Get("type")
	if !ok || tokenType != "access-token" {
		return ErrInvalidTokenType
	}

	return nil
}
