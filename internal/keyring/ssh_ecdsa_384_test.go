package keyring

import (
	"testing"
)

func TestSSHECDSA384(t *testing.T) {
	testCase := &sshKeyTest{
		privateKey: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAiAAAABNlY2RzYS
1zaGEyLW5pc3RwMzg0AAAACG5pc3RwMzg0AAAAYQTkKUipicAcSQvJQbiL/wiUF7OH9Jlv
ySgJUW8N1igRksNTtjFna9P+/aoly1c3fBwjbcCCRuFiX/e+pnWEgTQStvVn9VKdLWTKm+
iHHFyPasww/VRbE6GqZaUr90f9EzkAAADgRRR6h0UUeocAAAATZWNkc2Etc2hhMi1uaXN0
cDM4NAAAAAhuaXN0cDM4NAAAAGEE5ClIqYnAHEkLyUG4i/8IlBezh/SZb8koCVFvDdYoEZ
LDU7YxZ2vT/v2qJctXN3wcI23AgkbhYl/3vqZ1hIE0Erb1Z/VSnS1kypvohxxcj2rMMP1U
WxOhqmWlK/dH/RM5AAAAMQCWrhd79pw6KUmgsyAsh4XkByF3ITUrqI54g2xdV+vTw5idkF
q8ROUYZ/V/7za+DkcAAAAWZWNkc2EtcDM4NEBleGFtcGxlLmNvbQE=
-----END OPENSSH PRIVATE KEY-----`,

		publicKey: `ecdsa-sha2-nistp384 AAAAE2VjZHNhLXNoYTItbmlzdHAzODQAAAAIbmlzdHAzODQAAABhBOQpSKmJwBxJC8lBuIv/CJQXs4f0mW/JKAlRbw3WKBGSw1O2MWdr0/79qiXLVzd8HCNtwIJG4WJf976mdYSBNBK29Wf1Up0tZMqb6IccXI9qzDD9VFsToaplpSv3R/0TOQ==`,

		keyTest: keyTest{
			jwaSignatureAlgorithm: "ES384",

			jwkThumbprintSHA256:  "8Dzyd-wqHMkaPoPHL_M5jRTrbVz_30UY1ZkqQpDX1i0",
			sshFingerprintSHA256: "SHA256:YA8n5UphkaaCdqG9xaUuUjn8Niqkt631sNwH1LB7uQ8",
			comment:              "ecdsa-p384@example.com",
		},
	}

	testCase.run(t)
}
