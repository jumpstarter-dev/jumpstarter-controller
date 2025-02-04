package oidc

import (
	"context"
	"crypto/ecdsa"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Signer struct {
	privatekey *ecdsa.PrivateKey
}

func NewSigner(privateKey *ecdsa.PrivateKey) Signer {
	return Signer{
		privatekey: privateKey,
	}
}

func (k *Signer) ID() string {
	return "default"
}

func (k *Signer) Algorithm() jose.SignatureAlgorithm {
	return SignatureAlgorithm
}

func (k *Signer) Use() string {
	return "sig"
}

func (k *Signer) Key() any {
	return k.privatekey.Public()
}

func (k *Signer) KeySet(context.Context) ([]op.Key, error) {
	return []op.Key{k}, nil
}

func (k *Signer) Register(group gin.IRoutes) {
	group.GET("/.well-known/openid-configuration", func(c *gin.Context) {
		op.Discover(c.Writer, &oidc.DiscoveryConfiguration{
			Issuer:  Issuer,
			JwksURI: Issuer + "/jwks",
		})
	})

	group.GET("/jwks", func(c *gin.Context) {
		op.Keys(c.Writer, c.Request, k)
	})
}

func (k *Signer) Token(
	subject string,
) (string, error) {
	return jwt.NewWithClaims(SigningMethod, jwt.RegisteredClaims{
		Issuer:    Issuer,
		Subject:   strings.TrimPrefix(subject, Prefix),
		Audience:  []string{Audience},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(365 * 24 * time.Hour)), // FIXME: rotate keys on expiration
	}).SignedString(k.privatekey)
}
