package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"utils/pkg/stacks"
)

// TokenCarrier is a Scenario that carries the jwt token
type TokenCarrier interface {
	TokenStr() string
	SetToken(token *jwt.Token)

	stacks.Scenario
}

// Verify is a stackable that verifies the jwt token and calls the next handler if the token is valid
// This stackable uses the jwt.Keyfunc to resolve the key for the token
//
// Please note that this stackable does not validate the claims of the token
// To do that, you can add another stackable behind this stackable
func Verify[Scenario TokenCarrier](keyFunc jwt.Keyfunc) stacks.Stackable[Scenario] {
	return func(next stacks.Handler[Scenario]) stacks.Handler[Scenario] {
		return func(tc Scenario) error {
			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tc.TokenStr(), claims, keyFunc)

			if err != nil {
				return err
			}

			tc.SetToken(token)

			return next(tc)
		}
	}
}
