package keyring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type sshKeyTest struct {
	keyTest
	privateKey string
	publicKey  string
}

func (s *sshKeyTest) run(t *testing.T) {
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
		// Expect an empty comment for public key checks
		oldComment := s.comment
		defer func() {
			s.comment = oldComment
		}()
		s.comment = ""

		key, err = ParseString(s.publicKey)
		require.NoError(t, err)
		require.NotNil(t, key)
		require.False(t, key.IsPrivate())

		s.keyTest.run(t, key)
	})

	t.Run("AuthorizedKey", func(t *testing.T) {
		t.Run("Simple", func(t *testing.T) {
			authorizedKey := fmt.Sprintf("foo,bar,baz %s %s", s.publicKey, s.comment)

			key, err = ParseString(authorizedKey)
			require.NoError(t, err)
			require.NotNil(t, key)
			require.False(t, key.IsPrivate())

			s.keyTest.run(t, key)
		})
	})
}
