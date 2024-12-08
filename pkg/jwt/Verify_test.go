package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"testing"
)

type tc struct {
	token    string
	jwtToken *jwt.Token
}

func (t *tc) TokenStr() string {
	return t.token
}

func (t *tc) SetToken(token *jwt.Token) {
	t.jwtToken = token
}

func TestVerify_ValidToken(t *testing.T) {
	secret := []byte("secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	tokenString, _ := token.SignedString(secret)

	carrier := &tc{
		token: tokenString,
	}

	handler := Verify[*tc](func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	err := handler(func(scenario *tc) error {
		return nil
	})(carrier)

	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

func TestVerify_InvalidToken(t *testing.T) {
	tokenString := "invalid token"

	carrier := &tc{
		token: tokenString,
	}

	handler := Verify[*tc](func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	err := handler(func(scenario *tc) error {
		return nil
	})(carrier)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
