package encmgr

import (
	"common/encrypt"
	"crypto/rsa"
	"gamesrv/conf"
	"sync"

	"github.com/astaxie/beego"
)

//注意api服务器的公钥私钥不是一对,在整个服务器框架系统公钥统一使用一样的用作加密,私钥也同样使用一样的用作解密(注意定期更换RSA密钥)
type EncMgr struct {
	rsaPuk *rsa.PublicKey  //RSA 公钥加密
	rsaPrk *rsa.PrivateKey //RSA 私钥解密

	aesPuk string
	aesPrk string

	rsaCipher []byte
	aesCipher []byte
}

var sInstance *EncMgr
var once sync.Once

func Instance() *EncMgr {
	once.Do(func() {
		sInstance = &EncMgr{}
	})

	return sInstance
}

func (o *EncMgr) Init() error {
	puk, err := encrypt.DecodePuk(conf.RsaPuk)
	if err != nil {
		return err
	}
	o.rsaPuk = puk

	prk, err := encrypt.DecodePrk(conf.RsaPrk)
	if err != nil {
		return err
	}
	o.rsaPrk = prk

	o.rsaCipher = conf.RsaCipher

	beego.Info("--- Init Encrypt Mgr Done !")
	return nil
}

//封装加密(逻辑编写时不用关心key的问题)
func (o EncMgr) RsaEnc(plaintext []byte) ([]byte, error) {
	return encrypt.RsaEncrypt(o.rsaPuk, plaintext)
}

func (o EncMgr) RsaDec(cipher []byte) ([]byte, error) {
	return encrypt.RsaDecrypt(o.rsaPrk, cipher)
}

//Aes puk 加密用于服务器和客户端的通信
func (o EncMgr) AesPukEnc(plaintext []byte) ([]byte, error) {
	return encrypt.AesEnc(o.aesPuk, plaintext)
}

//Aes prk 加密用户服务器之间的通讯
func (o EncMgr) AesPrkEnc(plaintext []byte) ([]byte, error) {
	return encrypt.AesEnc(o.aesPrk, plaintext)
}

func (o EncMgr) AesPukDec(cipher []byte) ([]byte, error) {
	return encrypt.AesDec(o.aesPuk, cipher)
}

func (o EncMgr) AesPrkDec(cipher []byte) ([]byte, error) {
	return encrypt.AesDec(o.aesPrk, cipher)
}

func (o *EncMgr) SetAesPuk(puk string) {
	o.aesPuk = puk
}

func (o *EncMgr) SetAesPrk(prk string) {
	o.aesPrk = prk
}
