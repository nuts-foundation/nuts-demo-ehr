package keyring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"

	"github.com/google/uuid"
)

type keyTest struct {
	comment               string
	sshFingerprintSHA256  string
	jwaSignatureAlgorithm string
	jwkThumbprintSHA256   string
}

func (s *keyTest) run(t *testing.T, key Key) {
	t.Run("JWASignatureAlgorithm", func(t *testing.T) {
		alg, err := key.JWASignatureAlgorithm()
		assert.NoError(t, err)
		assert.Equal(t, s.jwaSignatureAlgorithm, string(alg))
	})

	t.Run("JWKThumbprintSHA256", func(t *testing.T) {
		keyID, err := key.JWKThumbprintSHA256()
		assert.NoError(t, err)
		assert.Equal(t, s.jwkThumbprintSHA256, keyID)
	})

	t.Run("JWKThumbprintURI", func(t *testing.T) {
		keyID, err := key.JWKThumbprintURI()
		assert.NoError(t, err)
		assert.Equal(t, "urn:ietf:params:oauth:jwk-thumbprint:sha-256:"+s.jwkThumbprintSHA256, keyID)
	})

	t.Run("SSHFingerprintSHA256", func(t *testing.T) {
		fingerprint, err := key.SSHFingerprintSHA256()
		assert.NoError(t, err)
		assert.Equal(t, s.sshFingerprintSHA256, fingerprint)
	})

	t.Run("JWK", func(t *testing.T) {
		jwk := key.JWK()

		t.Run("KeyID", func(t *testing.T) {
			assert.Equal(t, s.jwkThumbprintSHA256, jwk.KeyID())
		})
	})

	t.Run("Comment", func(t *testing.T) {
		comment := key.Comment()
		assert.Equal(t, s.comment, comment)
	})

	// For private keys test JWT signing
	if key.IsPrivate() {
		t.Run("PublicKey", func(t *testing.T) {
			publicKey, err := key.PublicKey()
			require.NoError(t, err)
			require.NotNil(t, publicKey)

			t.Run("KeyID", func(t *testing.T) {
				assert.Equal(t, s.jwkThumbprintSHA256, publicKey.JWK().KeyID())
			})
		})

		t.Run("SignJWT", func(t *testing.T) {
			// Generate a new JWT
			token, err := jwt.NewBuilder().
				JwtID(uuid.NewString()).
				Issuer("test@example.com").
				Subject("test@example.com").
				Build()
			require.NoError(t, err)

			// Save the payload state before signing to compare later
			payload, err := jwt.NewSerializer().Serialize(token)
			require.NoError(t, err)

			// Sign the JWT
			signed, err := key.SignJWT(token)
			require.NoError(t, err)

			// Get the public key which will be used for verification
			publicKey, err := key.PublicKey()
			require.NoError(t, err)

			// Determine the algorithm to verify with
			signatureAlgorithm, err := publicKey.JWASignatureAlgorithm()
			require.NoError(t, err)

			// Verify the signature on the resulting JWT
			verified, err := jws.Verify(signed, jws.WithKey(signatureAlgorithm, publicKey.JWK()))
			require.NoError(t, err)

			// Ensure the verified bytes match the original payload
			require.Equal(t, payload, verified)
		})
	}
}
