package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

//解析公钥
func DecodePuk(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("Public Key is not pem formate")
	}

	puk, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return puk.(*rsa.PublicKey), nil
}

//解析私钥
func DecodePrk(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("Private Key is not pem formate")
	}

	prk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return prk, nil
}

//Rsa公钥加密(错误处理交由使用部分来判断,这里只负责封装基础函数)
func RsaEncrypt(puk *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, puk, plaintext)
}

//Rsa私钥解密(错误处理交由使用部分来判断,这里只负责封装基础函数)
func RsaDecrypt(prk *rsa.PrivateKey, cipher []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, prk, cipher)
}
