package keyring

import (
	"testing"
)

func TestSSHECDSA521(t *testing.T) {
	testCase := &sshKeyTest{
		privateKey: `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAArAAAABNlY2RzYS
1zaGEyLW5pc3RwNTIxAAAACG5pc3RwNTIxAAAAhQQAzXQ9liCgmTFCs73KFlxNQKi4yoyD
ykghy0pkAP8O1gwO08v2G01fdXLZfLHqIymNRnC7jkNtJpcRthty5Z/gxgEA3C9F5iAXnY
RNq0FqgaaW5ejGCOOyMxjSVMHvshwGA1sJFAFxpvlbkiZBqf+yJ9nZwmhHz9R7e1dDnKgr
xuc9gyMAAAEYxrfaWsa32loAAAATZWNkc2Etc2hhMi1uaXN0cDUyMQAAAAhuaXN0cDUyMQ
AAAIUEAM10PZYgoJkxQrO9yhZcTUCouMqMg8pIIctKZAD/DtYMDtPL9htNX3Vy2Xyx6iMp
jUZwu45DbSaXEbYbcuWf4MYBANwvReYgF52ETatBaoGmluXoxgjjsjMY0lTB77IcBgNbCR
QBcab5W5ImQan/sifZ2cJoR8/Ue3tXQ5yoK8bnPYMjAAAAQgHyTOzp/uBUSScbur5XavOz
c96uqDOlqCerzMlfBWQK7Q7lFeDPxHvPf+1GWAGuSpW5r1i4z1Ik/BO4TRcfPlJvuwAAAB
ZlY2RzYS1wNTIxQGV4YW1wbGUuY29tAQIDBA==
-----END OPENSSH PRIVATE KEY-----`,

		publicKey: `ecdsa-sha2-nistp521 AAAAE2VjZHNhLXNoYTItbmlzdHA1MjEAAAAIbmlzdHA1MjEAAACFBADNdD2WIKCZMUKzvcoWXE1AqLjKjIPKSCHLSmQA/w7WDA7Ty/YbTV91ctl8seojKY1GcLuOQ20mlxG2G3Lln+DGAQDcL0XmIBedhE2rQWqBppbl6MYI47IzGNJUwe+yHAYDWwkUAXGm+VuSJkGp/7In2dnCaEfP1Ht7V0OcqCvG5z2DIw==`,

		keyTest: keyTest{
			jwaSignatureAlgorithm: "ES512",

			jwkThumbprintSHA256:  "MCgcNSB-O6RMER5DRy8igppoXmNlQlgL4trVMyTSu-o",
			sshFingerprintSHA256: "SHA256:MnCgkjy8AeBFVPKriTSRPTjxtLjIt6fkefjB6KFo0e4",
			comment:              "ecdsa-p521@example.com",
		},
	}

	testCase.run(t)
}
