package oidc

import (
	"context"
	"crypto/ecdsa"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Signer struct {
	pvivateKey *ecdsa.PrivateKey
}

func NewSigner(privateKey *ecdsa.PrivateKey) Signer {
	return Signer{
		pvivateKey: privateKey,
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
	return k.pvivateKey.Public()
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
