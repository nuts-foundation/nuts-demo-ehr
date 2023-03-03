package keyring

import (
	"testing"
)

func TestSSHECDSA256(t *testing.T) {
	testCase := &sshKeyTest{
		privateKey: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAaAAAABNlY2RzYS
1zaGEyLW5pc3RwMjU2AAAACG5pc3RwMjU2AAAAQQRZB/pmz2AMyCEhqoiaO+i80HAFtwF5
torqQPGGM3rC1eLPYM7xnAvEpEk5o363/ILhvhaV2IBdrGqZO+qw9S9+AAAAsIFR1FeBUd
RXAAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFkH+mbPYAzIISGq
iJo76LzQcAW3AXm2iupA8YYzesLV4s9gzvGcC8SkSTmjfrf8guG+FpXYgF2sapk76rD1L3
4AAAAgDuL/Lu6MTdMqK1/ri7T0nYwpF/brdpLx78Y/K7Zdv28AAAAVZWNkc2EtMjU2QGV4
YW1wbGUuY29tAQID
-----END OPENSSH PRIVATE KEY-----`,

		publicKey: `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFkH+mbPYAzIISGqiJo76LzQcAW3AXm2iupA8YYzesLV4s9gzvGcC8SkSTmjfrf8guG+FpXYgF2sapk76rD1L34=`,

		keyTest: keyTest{
			jwaSignatureAlgorithm: "ES256",

			jwkThumbprintSHA256:  "tPWIY0wU1CCmwZBIOEYGhVoC9bOJNOoQjRlAyIHoFxo",
			sshFingerprintSHA256: "SHA256:tecEJxiXYIPisGAIHw52wnETyNamO+49OIG3hJydR5k",
			comment:              "ecdsa-256@example.com",
		},
	}

	testCase.run(t)
}
