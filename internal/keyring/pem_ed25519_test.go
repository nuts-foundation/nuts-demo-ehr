package keyring

import (
	"testing"
)

func TestPEMEd25519(t *testing.T) {
	testCase := &pemKeyTest{
		privateKey: `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIIUqgZmCL47PAKNC8goFPCcesC0YrZ3q7tkYNECSMAKA
-----END PRIVATE KEY-----`,

		publicKey: `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAK1PMFx3WlYZ0NH+Yu3mELuqVTr5/3ZJmzsL3JgLyr6Q=
-----END PUBLIC KEY-----`,

		keyTest: keyTest{
			jwaSignatureAlgorithm: "EdDSA",

			jwkThumbprintSHA256:  "PInK3J6s_eDE44b3M1Hr2CS-rCI7y0D4cD8z5gBnkyI",
			sshFingerprintSHA256: "SHA256:gT/HP4shumEFWcVKQ1eyu/moz9h9WPDzAHRDUkfWz7Y",
		},
	}

	testCase.run(t)
}
