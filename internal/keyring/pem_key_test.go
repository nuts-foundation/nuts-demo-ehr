package keyring

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type pemKeyTest struct {
	keyTest
	privateKey string
	publicKey  string
}

func (s *pemKeyTest) run(t *testing.T) {
	var key Key
	var err error

	t.Run("PrivateKey", func(t *testing.T) {
		key, err = ParseString(s.privateKey)
		require.NoError(t, err)
		require.NotNil(t, key)
		require.True(t, key.IsPrivate())

		s.keyTest.run(t, key)
	})

	t.Run("PublicKey", func(t *testing.T) {
		key, err = ParseString(s.publicKey)
		require.NoError(t, err)
		require.NotNil(t, key)
		require.False(t, key.IsPrivate())

		s.keyTest.run(t, key)
	})
}
