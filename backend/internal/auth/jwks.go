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
	"sync"
	"time"
)

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

	if cognitoRegion == "" ||
		cognitoUserPoolID == "" {

		return nil, errors.New(
			"missing cognito config",
		)
	}

	url := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		cognitoRegion,
		cognitoUserPoolID,
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
