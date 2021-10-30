package ticket

import (
	"dxkite.cn/gateway/util"
	"fmt"
	"testing"
)

const priKeyPem = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAzXI2aH4v77zn7zpKfyz8TQEXGojj661ywmWUW5VPWV7oEYXP
wP4yPzN4UhesVWTILA3mvCOBDqYl3SSAoSnGbFw4HvWXjU/8EBJJOFg4yinD6SvA
17YQzjvDBJJzXPlgYxsb/gBkkXzxewJz0ZH3YvJqUrUk0n2PdEG+Mi9yEnZ1us9h
gpCgXScpbi1Qgk8N2IiGIDN1oQPibfDKDrSY9z/p7f+eEE24GcbGrirQiEVvRTcl
TMaOcpNkK+Qy2/oNuDinNTjNvwBJqbt74dGjbjRbKzZu+STpurHwlRSoW2bShI+X
v5ODz5iuZO2PCdMJ1/zXdHfzoI8KkZUWDjUcoQIDAQABAoIBABOqJ/0OfwYifczP
Now/wRKZ1R1MFwza+E4VLQMrIoI6bFopBb7CVgrooU6yR2ORFvvohLpjZ5oAW0V0
lf9XIVAD77c+6Le0/TqIlYuYHi9zmfX6oyXCno8u9za70eEHWauCz8YbQYUPgYm1
rWJU52oILBcokZK4Q//rvwnMss2Dbz4IAHQ4UZ2vADlJBbXwbVIrHTZmrZnZP2pT
t8RFKy4/loEEudnjzx0nVSSAbnYy2S1uRFAJb4RcJywTZS3aFNu2qjLlBO6ZeaAv
XwQ2++MC0LvNyEfriTs8F6cuwX6xqd0LUGdpx2uKak9oIPUFj6dKdkbiC+lRYn9s
mgMqnbECgYEA8sM2oJaHpzxHedtaDxwdRfW+YXiL7q4zwYMprDSb6hVGMp6l3DdU
Qugo25wQGqKhncqZYc1oved5f0q35mcpOdCciJQHZuqlbH3ZHeQt8UmQ6ykaFCTc
b7mLyn0Elf/0Y7V3bxgB6qlENOZl40/2d2RrjusCClhX1kpCF5SIiEUCgYEA2KYW
6vJzkznn5HOOCswXSNyHKKIb1cWFMOQe0CaJLWq9FTBI8uqoseucHEtL2NhnpjWM
REfqi4rUNcFb6c8TkEqNZBCtL1f9pSdF2m/TRqZ3dHjYpdmOF4os2KAwCfkMTc5F
W5iqatu8gUFhAlF9YELq4OggTPcyMS9S73d8Tq0CgYBVLKzz9xytTnb9iDq25nRW
4Xvkkvj1y2UZVj2+z86MeN5iUEt9UmRb/TyooL79uWXfCQB70igXySlVwg935WYP
hOQG/3kBYP6dbCJLXI3KBLe16nvd6Xj2MjGb3/VF88H5Yef/sHqrrKvjq5rAAIRH
K5KZWFck7g4Tf4Zk45ZryQKBgHq9NhSrgVDyqG7kDKAPWk28KpKZrN1ihv/Y7aAN
hQAHDdKYRWviB+qsyge5nOHgUHB4u9vfRoECCRHfVvxShgnkQtBjJkrBNgFAC7Ii
UncfTmPdJxhm9bpeXOPpdO3he9gEuYSYLExX6ybrbFNM6ZQEtV7wA4S3M2dsITdr
4TANAoGBAJGttoSNVwOTDFXBPEB7xp+rxlok02kH1F5B4YaWg6Cxn1/IBGDboS+O
KANKiUkSSrvEQlEQzLXpoXuDHRrI1qo5xK6dQ1FzQp9rT4CAWRiaTET0kSUu1qIJ
HaukToaj7WuWtuBtKApg7VFOjKy7QJ2ISCSBQS02rdjFT1ZepCMn
-----END RSA PRIVATE KEY-----
`

const certPem = `-----BEGIN CERTIFICATE-----
MIIDLzCCAhegAwIBAgIRANymJ0/SWMYyLc/LrHiKLF4wDQYJKoZIhvcNAQELBQAw
HzELMAkGA1UEBhMCQ04xEDAOBgNVBAMTB1JPT1QgQ0EwHhcNMjExMDIzMTQyMzEw
WhcNMzExMDIxMTQyMzEwWjAfMQswCQYDVQQGEwJDTjEQMA4GA1UEAxMHZ2F0ZXdh
eTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAM1yNmh+L++85+86Sn8s
/E0BFxqI4+utcsJllFuVT1le6BGFz8D+Mj8zeFIXrFVkyCwN5rwjgQ6mJd0kgKEp
xmxcOB71l41P/BASSThYOMopw+krwNe2EM47wwSSc1z5YGMbG/4AZJF88XsCc9GR
92LyalK1JNJ9j3RBvjIvchJ2dbrPYYKQoF0nKW4tUIJPDdiIhiAzdaED4m3wyg60
mPc/6e3/nhBNuBnGxq4q0IhFb0U3JUzGjnKTZCvkMtv6Dbg4pzU4zb8ASam7e+HR
o240Wys2bvkk6bqx8JUUqFtm0oSPl7+Tg8+YrmTtjwnTCdf813R386CPCpGVFg41
HKECAwEAAaNmMGQwDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMB
BggrBgEFBQcDAjAfBgNVHSMEGDAWgBTNjFY8dlTstlww1vbwQbAlu8d6bjASBgNV
HREECzAJggdnYXRld2F5MA0GCSqGSIb3DQEBCwUAA4IBAQBYC2XhMNJ90RbmgRu1
tB7tzO0lBT8bwUaZIg55SUx9FBQZLpVdoPAUFcHv/dxR6nc/bW3rINglDQ64B6lK
T1zNzsSJ43oUlccGnq36DO9fcAITL6yB7eo+FJVC832Bt66NF0d6JuxrnJ8jCLac
P8bAHX1QmWEcyeKTRtT9KyzxuioWamlinfUx5F+TXrruXa8rcGLD4pSNhGX8lsoM
uw82ZqF7M7s1h4R+RjB1/KmEAANNLI/a+zVQgfR5njgQ0eSFN9mHwpqrda9UwLin
SMuaywmuSsgribV/+am78E/eJF0ulA7NQnMoO1hZ81EoJ9diOR5bMT/VMXRm/jO8
puwS
-----END CERTIFICATE-----
`

func TestRsaTicket_Decode(t *testing.T) {
	priKey, err := util.ParsePrivateKey([]byte(priKeyPem))
	if err != nil {
		t.Error(err)
		return
	}

	pubKey, err := util.ParsePublicKeyFromCertificate([]byte(certPem))
	if err != nil {
		t.Error(err)
		return
	}

	rt := &RsaTicket{
		pri: priKey,
		pub: pubKey,
	}

	var uin uint64 = 1008611
	ticket, err := rt.Encode(&SessionData{
		Uin:        uin,
		CreateTime: 0,
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("ticket", ticket)
	tt, err := rt.Decode(ticket)
	if err != nil {
		t.Error(err)
		return
	}

	if tt.Uin != uin {
		t.Error("uin error", uin, tt.Uin)
	}
	fmt.Println("ticket", tt.Uin, tt.CreateTime)
}
