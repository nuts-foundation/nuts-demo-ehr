package keyring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sshKeyTest struct {
	privateKey        string
	publicKey         string
	comment           string
	sshFingerprintSHA256 string
	jwkThumbprintSHA256          string
}

func (s *sshKeyTest) checkKey(t *testing.T, data string) {
	var key Key

	t.Run("Parse", func(t *testing.T) {
		var err error
		key, err = ParseString(data)
		require.NoError(t, err)
		require.NotNil(t, key)
	})
	require.NotNil(t, key)

	t.Run("JWKThumbprintSHA256", func(t *testing.T) {
		keyID, err := key.JWKThumbprintSHA256()
		assert.NoError(t, err)
		assert.Equal(t, s.jwkThumbprintSHA256, keyID)
	})

	t.Run("SSHFingerprintSHA256", func(t *testing.T) {
		fingerprint, err := key.SSHFingerprintSHA256()
		assert.NoError(t, err)
		assert.Equal(t, s.sshFingerprintSHA256, fingerprint)
	})
}

func (s *sshKeyTest) run(t *testing.T) {
	t.Run("PrivateKey", func(t *testing.T) {
		s.checkKey(t, s.privateKey)
	})
}

