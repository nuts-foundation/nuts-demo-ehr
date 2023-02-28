// keyring provides utilities for loading and using crypto keys in different formats
package keyring

import (
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

type Key interface {
	JWASignatureAlgorithm() (jwa.SignatureAlgorithm, error)
	JWKThumbprintSHA256() (string, error)
	JWTSerializer() (*jwt.Serializer, error)
	SSHFingerprintSHA256() (string, error)
}

type keyImpl struct {
	comment string
	jwk     jwk.Key
}

// Open reads a key from file, returning a Key
func Open(path string) (Key, error) {
	// Read the key file
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read %v: %w", path, err)
	}

	return Parse(keyData)
}

// ParseString is a shortcut for calling Parse with a string
func ParseString(s string) (Key, error) {
	return Parse([]byte(s))
}

// Parse parsses raw key bytes
func Parse(bytes []byte) (Key, error) {
	// Try to parse the key as a JWK
	if key, err := jwk.ParseKey(bytes); err == nil {
		return &keyImpl{jwk: key}, nil
	}

	// Try to parse as PEM
	if key, err := jwk.ParseKey(bytes, jwk.WithPEM(true)); err == nil {
		return &keyImpl{jwk: key}, nil
	}

	// Try to get a raw crypto/* type from an OpenSSH private key
	if key, err := ssh.ParseRawPrivateKey(bytes); err == nil {
		// Load the raw crypto/* type
		if converted, err := from(key); err == nil {
			// Return the loaded key
			return converted, nil
		} else {
			return nil, err
		}
	}

	return nil, errors.New("failed to parse as JWK, PEM, or OpenSSH private key")
}

// From creates a Key from go standard library key types
func from(raw interface{}) (*keyImpl, error) {
	// Handle other types more directly
	switch raw := raw.(type) {
	case *ed25519.PrivateKey:
		if convertedToJWK, err := jwk.FromRaw(*raw); err == nil {
			return &keyImpl{jwk: convertedToJWK}, nil
		}
	default:
		if convertedToJWK, err := jwk.FromRaw(raw); err == nil {
			return &keyImpl{jwk: convertedToJWK}, nil
		}
	}
	return nil, fmt.Errorf("failed to convert %T", raw)
}

// Raw returns the raw crypto/* type of the key
func (k *keyImpl) Raw() (interface{}, error) {
	// Convert the key to its raw key type from crypto/*
	var rawKey interface{}
	if err := k.jwk.Raw(&rawKey); err != nil {
		return nil, fmt.Errorf("unable to convert to raw key: %w", err)
	}

	return rawKey, nil
}

// JWKThumbprintSHA256 returns the base64 encoded SHA256 JWK thumbprint of the key per RFC7638
func (k *keyImpl) JWKThumbprintSHA256() (string, error) {
	// Generate the SHA256 thumbprint
	hash, err := k.jwk.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("failed to generate thumbprint: %v", err)
	}

	// Return the base64 hash
	return base64.RawURLEncoding.EncodeToString(hash), nil
}

// JWASignatureAlgorithm() returns a suitable signing algorithm for the key
func (k *keyImpl) JWASignatureAlgorithm() (jwa.SignatureAlgorithm, error) {
	// Inspect the key type
	switch k.jwk.KeyType() {
	case jwa.Ed25519:
		return jwa.EdDSA, nil
	default:
		return "", fmt.Errorf("unexpected key type: %v", k.jwk.KeyType())
	}
}

// JWTSerializer returns a JWT serializer for the key
func (k *keyImpl) JWTSerializer() (*jwt.Serializer, error) {
	sigAlg, err := k.JWASignatureAlgorithm()
	if err != nil {
		return nil, err
	}

	signOption := jwt.WithKey(sigAlg, k.jwk)
	return jwt.NewSerializer().Sign(signOption), nil
}

// SSHPublicKey converts the public key component of the key into an ssh.PublicKey
func (k *keyImpl) SSHPublicKey() (ssh.PublicKey, error) {
	// Get the raw crypto/* key
	raw, err := k.Raw()
	if err != nil {
		return nil, err
	}

	// For the following steps we will ultimately need an ssh public key
	var publicKey ssh.PublicKey

	// Convert the rawKey into an SSH signer, which works if rawKey is a private key, but not public key
	if signer, err := ssh.NewSignerFromKey(raw); err == nil {
		publicKey = signer.PublicKey()
	} else {
		// Since creating a signer from the key failed perhaps we were given a public key, so try
		// creating an ssh.PublicKey directly from the rawKey
		if publicKey, err = ssh.NewPublicKey(raw); err != nil {
			return nil, fmt.Errorf("unable to create ssh.Signer or ssh.PublicKey from %T", raw)
		}
	}

	return publicKey, nil
}

// SSHFingerprintSHA256 returns the ssh SHA256 fingerprint of the key
func (k *keyImpl) SSHFingerprintSHA256() (string, error) {
	publicKey, err := k.SSHPublicKey()
	if err != nil {
		return "", err
	}

	return ssh.FingerprintSHA256(publicKey), nil
}

