package keyring

import (
	"testing"
)

func TestPEMECDSAprime256v1(t *testing.T) {
	testCase := &pemKeyTest{
		privateKey: `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIFkXiaRwDdX4lrzlM1UMufT57sKcDnuk7vCofSLDjBruoAoGCCqGSM49
AwEHoUQDQgAE6ctLnnYP/UXvAAm+U4nxxxMPG8FTmcndjdbWnC4bgkmjIc0J9bj7
+hAT0jFkfHI1N28aPPPSse7oDfvh8+xRHw==
-----END EC PRIVATE KEY-----`,

		publicKey: `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE6ctLnnYP/UXvAAm+U4nxxxMPG8FT
mcndjdbWnC4bgkmjIc0J9bj7+hAT0jFkfHI1N28aPPPSse7oDfvh8+xRHw==
-----END PUBLIC KEY-----`,

		keyTest: keyTest{
			jwaSignatureAlgorithm: "ES256",

			jwkThumbprintSHA256:  "Dk_8kgb9quICny5jMmM8tOSUyDzh3uLQsgS8Rtxyudc",
			sshFingerprintSHA256: "SHA256:+vuwzxZWWkUaCuYt6CsulcOLJVVBdcbkyhN75Yv+cEA",
		},
	}

	testCase.run(t)
}
