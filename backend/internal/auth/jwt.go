package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Sub   string
	Email string
}

// =========================
// Validate JWT
// =========================

func Validate(
	tokenString string,
) (*Token, error) {

	if cognitoRegion == "" ||
		cognitoUserPoolID == "" ||
		cognitoClientID == "" {

		return nil, errors.New(
			"missing cognito config",
		)
	}

	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (any, error) {

			if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {

				return nil, errors.New(
					"invalid alg",
				)
			}

			kid, ok := t.Header["kid"].(string)

			if !ok {

				return nil, errors.New(
					"missing kid",
				)
			}

			return getPublicKey(kid)
		},
		jwt.WithIssuer(cognitoIssuer),
		jwt.WithAudience(cognitoClientID),
		jwt.WithValidMethods([]string{
			jwt.SigningMethodRS256.Alg(),
		}),
	)

	if err != nil || !token.Valid {

		return nil, errors.New(
			"invalid token",
		)
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {

		return nil, errors.New(
			"invalid claims",
		)
	}

	// =========================
	// token_use validation
	// =========================

	tokenUse, ok := claims["token_use"].(string)

	if !ok || tokenUse != "id" {

		return nil, errors.New(
			"invalid token_use",
		)
	}

	sub, ok := claims["sub"].(string)

	if !ok || sub == "" {

		return nil, errors.New(
			"invalid sub",
		)
	}

	email, _ := claims["email"].(string)

	return &Token{
		Sub:   sub,
		Email: email,
	}, nil
}
