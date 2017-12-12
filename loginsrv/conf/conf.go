package conf

//注意 这个文件中得东西以后都要写到加密配置文件中,这样才方便修改

//公钥用于加密消息(所有Ctrl之外的服务器都使用这个加密发送消息到Ctrl)
var RsaPuk = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo+Low3DkBUwY+2Lr/FjT
xB+b2cXbTWH+KQ0aAzwChVVE/FESIfkGJ1QjBq8krzwgpSS/x0LRYF9Bx9u2Q88t
PSy6KSPcbf+SyKtw6eLt250nOSIDUZFkC7MZldMEWKSagIhOyUaBZnNNZ/WABPwy
BGgsPljMMf+GJACHRyvnE0X2iBm5y86RuLHhCBg8hXPYoj9kpUU1hBrwhEvcGw16
qgpSU1mZoUseLZY2uRKDV93URQoZ5tQ11ffYIz2jicGq/y/T63uqJR5OfVaAjqCR
3SC/5TAFCwE8nrNNV5pPdwvLBwquxWI21u4XiP3QIBf4Tdb3mEPvcU/QZmzkTDNR
oQIDAQAB
-----END PUBLIC KEY-----
`)

//RSA 私钥用于解密消息(所有 Ctrl 发送其他服务器的消息都使用这个解密)
var RsaPrk = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA5r1fp4nXNrfp0WRetwjnaFQ/B88Vws6WPo0GgGhZL0SpbYDG
RWSfP/eIlE2N1MvOb6Fe74N3lhCkHlgzQuS+zNLq72AbR1DRXRxYD59jdN2fX/ME
KsHFvvRPSTl201TtjJSor7A4w3YgJn1sSzbbddma1t8BVtqFDbKzKpqLQQscsvbK
lNi2y/BHbpNFQtpSFYctGsYwTrB/mVB/wMwKnoZkPAIiljeOrodT9Rce7Dgsc5WY
JsNXRxDM8jO+g554c8wy3euq3CUFKOSaxA8n6PpGDnIWsYt7O7ch6BhCiLoLQpsw
+GPdLbVxryp4ZOtk+rFr7/5iihNUDR+PnlZwTQIDAQABAoIBABo/SoVFYczgqOFf
2IJbqe8tPipGPUR2uZKN+kJbHGILHnbNYkB8jLz7DHdKRysAsA/0vFbkjpLse67T
+6jNWoL0LzNhrFi1cct0sPO9/tNJbpq8iynN9w+cvPQobELn80m9k17A3gQmCMw4
EjyQFfDW+w0cOwvFOcNwI39iKxsP0DQYHaqeJomkwO37x0M/SgreT1SahR8vo5Hx
no7fSMhZWGxPmOg2zOfH4NX7BAzoS5K/TyWC3PQlg5TVIAZJ/AXAsJ6HR/rdFR+I
9RFmNzk+SstkPWNQVALLjr4HU/YWKtywd1sCeEPCzdYy40zDssn15Ekw6i9dd3MS
fbqLgUUCgYEA/CJYR+y/novdPPFWEvY/IwL+3T8VZPKnchhar7o45VHNSR7GhsXg
OJQ5wfni2DflFKIIhFHdsTs/D9noBo1ZZ1dXYhoNVxMzOOMvDUZ+vDVsUhbeig+r
C+Df1WyIf52zwtsEvCD0LxXVJlnaRwN/Hvw/jd/5Cg7VqMbMWQJvNosCgYEA6kcN
pAssY1uJLt0XGwz01YKJD4JYPhjxx25x3Q7JjRX7Ue5/Cy7RnV8Oz6ASlpSVSGhF
RwmtJPbSZoNx1xf01uFv0tFmzw7tDz9I826MUBxxjZP06LYljIoA9+J5aCUWhvEJ
dRmUcrZrLS7lnuaXUOQm3J4TcjnwewJnYJrUp4cCgYB2JsYF1mypKFOhdlpmglxt
1L6IAULOTmOnNMBybqqw05eGd1SC3YFIIjW6r5Xcyryf4ZpqH07q+Z+AlxWC1IQb
yDMbtqefsVCkjNmEuA81tTcwdKUOP29hHpzlj3mbi9QsMKRUYIDs+6cp1JtUAdB1
PbGuk1FJpE/9SkOCRoDFJwKBgEPDQYLNaI1kkI1pjaFZYgfPte0yemubs3NH0s9p
04pnbUYJgd2uMRsfv5z2Y+oCGBvVbRRbDCXA7qKQKoFYgI0Wr81+nAoP+ymJ2IPw
2cziLUYSIaid5sZ7tEP+0bb540YsuduRBosXkHCFPA12DRZsp4DwiBdmAtTRoS0k
G5ZJAoGBAK1MqVNTqKAmziacvibC6eyf3SUBANE2d9DKAGmK+umerXR9KN1jNc4s
dkoMmo8dGa+rIT6RFeNuCEP99/R4AmmdQfn7Wl2NDXoPoZGGWg3pYpsl3oH3xXNv
nVmMHBlxlD7m0cq4X4PkDBHLS3uUvv0fdvETspRPRF/EMShez0kx
-----END RSA PRIVATE KEY-----
`)

//RSA 暗号
var RsaCipher = []byte("天王盖地虎")

//AES 暗号
var AesCipher = []byte("宝塔镇河妖")

//ctrlsrv IP(由于启动的时候必须从控制服务获取信息,所以这个只有再在每个服务器代码中写死,等项目空闲时改为加密配置文件的形式)
//var CtrlsrvIP = "http://192.168.1.182:8877"

//线上
var CtrlsrvIP = "http://192.168.1.151:8877"

//超时时间
var Timeout = "10s"

//本服务器功能
var SrvFunc = "login"
