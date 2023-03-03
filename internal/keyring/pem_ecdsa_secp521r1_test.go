package keyring

import (
	"testing"
)

func TestPEMECDSAsecp521r1(t *testing.T) {
	testCase := &pemKeyTest{
		privateKey: `-----BEGIN EC PRIVATE KEY-----
MIHcAgEBBEIBUCI/WGmRoL8mtBexBKAUrFi9s4mZcfIeNtp0ILiMBBq/yypK3FFv
8ezCNBF+owuAKb0yM68ENuJ7TfbgLBpaq6+gBwYFK4EEACOhgYkDgYYABAF4zKyc
MP2K4HeDbBGpiGlsSdo2cBZ5wPH/33PRIQHYj8RL6aSLcJ6+EIuKuhD22NZaKZ1E
dyO4t5b1MusJvkzm8QGxrgaoLITad0SiKrnaOQtHBuFgsLHeb3Av3R1wr5MvOjZl
wrXsDlmJpPgVy9YeFdJp3W7QntzHBvQRIMTz4Icc7A==
-----END EC PRIVATE KEY-----`,

		publicKey: `-----BEGIN PUBLIC KEY-----
MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBeMysnDD9iuB3g2wRqYhpbEnaNnAW
ecDx/99z0SEB2I/ES+mki3CevhCLiroQ9tjWWimdRHcjuLeW9TLrCb5M5vEBsa4G
qCyE2ndEoiq52jkLRwbhYLCx3m9wL90dcK+TLzo2ZcK17A5ZiaT4FcvWHhXSad1u
0J7cxwb0ESDE8+CHHOw=
-----END PUBLIC KEY-----`,

		keyTest: keyTest{
			jwaSignatureAlgorithm: "ES512",

			jwkThumbprintSHA256:  "7fYmHFyUxdB8qmfi6VZdqW94qVe8fxmK6XjyBRWUXt4",
			sshFingerprintSHA256: "SHA256:UZe3gXjTwrnGrLdwyqxC61p/w0kCtFTnBqcaTYss5Ug",
		},
	}

	testCase.run(t)
}
