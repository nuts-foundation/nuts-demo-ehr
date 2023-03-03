package keyring

import (
	"testing"
)

func TestPEMECDSAsecp384r1(t *testing.T) {
	testCase := &pemKeyTest{
		privateKey: `-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDBqwjZtjOqlBVcY0GinaFX3T0kdbuBGmjBr65pts0DiUXQXFL8Dc/Fl
TIOAhYExXY+gBwYFK4EEACKhZANiAASpvW5b6uaSjBqAsV+GSn1QgnQiN4cYz0E3
rUH6NSH90J0VLiIQvQAiwZ52ZKh0i/LFVEKYh7ft++Hfc6Kwbt4xKc9WL/QCpOSb
JJdUcnK832oKmmug5x+grPGElEHVEVs=
-----END EC PRIVATE KEY-----`,

		publicKey: `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEqb1uW+rmkowagLFfhkp9UIJ0IjeHGM9B
N61B+jUh/dCdFS4iEL0AIsGedmSodIvyxVRCmIe37fvh33OisG7eMSnPVi/0AqTk
mySXVHJyvN9qCpproOcfoKzxhJRB1RFb
-----END PUBLIC KEY-----`,

		keyTest: keyTest{
			jwaSignatureAlgorithm: "ES384",

			jwkThumbprintSHA256:  "daUv4ZmBNTCR9A4xENC7JrAyr8PVdtYi2_YsUbyIVFk",
			sshFingerprintSHA256: "SHA256:yMlhnbo7LWcgLupRz93nRXALNX/Mz4xh/CFglOiAboQ",
		},
	}

	testCase.run(t)
}
