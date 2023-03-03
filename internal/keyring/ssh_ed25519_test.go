package keyring

import (
	"testing"
)

func TestSSHEd25519(t *testing.T) {
	testCase := &sshKeyTest{
		privateKey: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACAAd91dPfRMknEScEgFkjM9QoKhl7YzUH43w4KZLR0XNQAAAJgAsjLQALIy
0AAAAAtzc2gtZWQyNTUxOQAAACAAd91dPfRMknEScEgFkjM9QoKhl7YzUH43w4KZLR0XNQ
AAAECfuVISGTugvoO5sW/xhDTDZ9sluBymljJcib5zDRvaUAB33V099EyScRJwSAWSMz1C
gqGXtjNQfjfDgpktHRc1AAAAE2VkMjU1MTlAZXhhbXBsZS5jb20BAg==
-----END OPENSSH PRIVATE KEY-----`,

		publicKey: `ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAB33V099EyScRJwSAWSMz1CgqGXtjNQfjfDgpktHRc1`,

		keyTest: keyTest{
			jwaSignatureAlgorithm: "EdDSA",

			jwkThumbprintSHA256:  "z_dnv54m05AgzoTOAljQG2Pfiu_4uZ9Zf2MUxNYjPUA",
			sshFingerprintSHA256: "SHA256:HbghS6X2Xe4GgtBIV+aGl0RZmwvxnTO/CcvjLFmvMcU",

			comment: "ed25519@example.com",
		},
	}

	testCase.run(t)
}
