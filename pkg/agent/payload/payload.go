// Package payload signs payloads with a given ECC private key.
package payload

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

// Signer represents a payload signer.
type Signer interface {
	Identifier() string
	PublicKey() (*ecdsa.PublicKey, error)
	Sign([]byte) ([]byte, error)
}

type payloadSigner struct {
	privateKey *ecdsa.PrivateKey
	identifier string
}

// NewSigner returns a new Signer.
func NewSigner(pubKeyStr, privKeyStr string) (Signer, error) {
	privateKey, err := parsePrivateKey(privKeyStr)
	if err != nil {
		return nil, err
	}

	id, err := generateIdentifier(pubKeyStr)
	if err != nil {
		return nil, err
	}

	return &payloadSigner{
		privateKey: privateKey,
		identifier: id,
	}, nil
}

// Identifier returns a sha256 encoded version of the public key.
func (s *payloadSigner) Identifier() string {
	return s.identifier
}

// PublicKey returns the public key that can be used to verify the signature.
func (s *payloadSigner) PublicKey() (*ecdsa.PublicKey, error) {
	return &s.privateKey.PublicKey, nil
}

// Sign returns the payload signature.
func (s *payloadSigner) Sign(body []byte) ([]byte, error) {
	hash := sha256.Sum256(body)
	return ecdsa.SignASN1(rand.Reader, s.privateKey, hash[:])
}

func generateIdentifier(key string) (string, error) {
	// Replacing {92, 110} with {10} is to replace the string "\n" with the rune '\n'.
	// This is necessary because of how go config parses env vars into the config struct.
	k := bytes.ReplaceAll([]byte(key), []byte{92, 110}, []byte{10})
	h := sha256.New()
	_, err := h.Write(k)
	if err != nil {
		return "", err
	}
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs), nil
}

func parsePrivateKey(privKeyStr string) (*ecdsa.PrivateKey, error) {
	// This ReplaceAll assumes that the key is a single line:
	// -----BEGIN EC PRIVATE KEY-----\nkey\n-----END
	privPemStr := strings.ReplaceAll(privKeyStr, "\\n", "\n")
	block, _ := pem.Decode([]byte(privPemStr))
	if block == nil {
		return nil, errors.New("failed to decode public key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// NoopSigner is a no-op signer for tests.
type NoopSigner struct{}

// Identifier returns the base64 encoded identifier of the signer.
func (ns *NoopSigner) Identifier() string { return "" }

// PublicKey returns the public key that can be used to verify the signature.
func (ns *NoopSigner) PublicKey() (*ecdsa.PublicKey, error) { return nil, nil }

// Sign returns the payload signature.
func (ns *NoopSigner) Sign([]byte) ([]byte, error) { return nil, nil }
