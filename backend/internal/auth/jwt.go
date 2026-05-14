package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Sub   string
	Email string
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

var (
	cacheMu    sync.RWMutex
	cachedKeys map[string]*rsa.PublicKey
	cachedAt   time.Time
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

// =========================
// Validate JWT
// =========================

func Validate(
	tokenString string,
) (*Token, error) {

	region := os.Getenv("AWS_REGION")
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	clientID := os.Getenv("COGNITO_APP_CLIENT_ID")

	if region == "" ||
		userPoolID == "" ||
		clientID == "" {

		return nil, errors.New(
			"missing cognito config",
		)
	}

	issuer := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s",
		region,
		userPoolID,
	)

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
		jwt.WithIssuer(issuer),
		jwt.WithAudience(clientID),
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

// =========================
// JWKS fetch + cache
// =========================

func getPublicKey(
	kid string,
) (*rsa.PublicKey, error) {

	cacheMu.RLock()

	if key, ok := cachedKeys[kid]; ok &&
		time.Since(cachedAt) < 30*time.Minute {

		cacheMu.RUnlock()

		return key, nil
	}

	cacheMu.RUnlock()

	region := os.Getenv("AWS_REGION")
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")

	if region == "" || userPoolID == "" {

		return nil, errors.New(
			"missing cognito config",
		)
	}

	url := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		region,
		userPoolID,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return nil, fmt.Errorf(
			"jwks fetch failed: status=%d",
			resp.StatusCode,
		)
	}

	var set jwks

	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return nil, err
	}

	newKeys := make(
		map[string]*rsa.PublicKey,
		len(set.Keys),
	)

	for _, k := range set.Keys {

		nBytes, err := base64.RawURLEncoding.DecodeString(k.N)

		if err != nil {
			continue
		}

		eBytes, err := base64.RawURLEncoding.DecodeString(k.E)

		if err != nil {
			continue
		}

		e := 0

		for _, b := range eBytes {
			e = e<<8 + int(b)
		}

		if e == 0 {
			e = 65537
		}

		newKeys[k.Kid] = &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: e,
		}
	}

	cacheMu.Lock()

	cachedKeys = newKeys
	cachedAt = time.Now()

	cacheMu.Unlock()

	key, ok := newKeys[kid]

	if !ok {

		return nil, errors.New(
			"key not found",
		)
	}

	return key, nil
}
