package payload

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"strings"
)

// Authenticator is for authenticating a payload using ECDSA.
type Authenticator interface {
	IsValid(ctx context.Context, body []byte, identifier, sig string) (bool, error)
}

// NewAuthenticator returns a new Authenticator.
func NewAuthenticator(pubKey string) (Authenticator, error) {
	k, err := parsePubKey(pubKey)
	if err != nil {
		return nil, err
	}

	return &authenticator{k}, nil
}

type authenticator struct {
	pubKey *ecdsa.PublicKey
}

// IsValid checks if the payload is valid.
func (a *authenticator) IsValid(ctx context.Context, data []byte, identifier, sig string) (bool, error) {
	// Parse the Webhook Signature
	parsedSig := asn1Signature{}
	asnSig, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false, err
	}
	rest, err := asn1.Unmarshal(asnSig, &parsedSig)
	if err != nil || len(rest) != 0 {
		return false, err
	}

	// Verify the SHA256 encoded payload against the signature with GitHub's Key
	digest := sha256.Sum256(data)
	keyOk := ecdsa.Verify(a.pubKey, digest[:], parsedSig.R, parsedSig.S)

	return keyOk, nil
}

func parsePubKey(pubKey string) (*ecdsa.PublicKey, error) {
	pubPemStr := strings.ReplaceAll(pubKey, "\\n", "\n")
	// Decode the Public Key
	block, _ := pem.Decode([]byte(pubPemStr))
	if block == nil {
		return nil, errors.New("error parsing PEM block with GitHub public key")
	}

	// Create our ECDSA Public Key
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	// Because of documentation, we know it's a *ecdsa.PublicKey
	ecdsaKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("GitHub key is not ECDSA")
	}

	return ecdsaKey, nil
}

// asn1Signature is a struct for ASN.1 serializing/parsing signatures.
type asn1Signature struct {
	R *big.Int
	S *big.Int
}
