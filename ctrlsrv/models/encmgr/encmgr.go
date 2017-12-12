package encmgr

import (
	"common/encrypt"
	"crypto/rsa"
	"ctrlsrv/conf"
	"ctrlsrv/models/dbmgr"
	"sync"

	"github.com/astaxie/beego"
)

type EncMgr struct {
	//RsaPuk string 目前没有找到再mongodb中存贮 RSA KEY 的方法等功能做完来找
	//RsaPrk string
	rsaCipher string
	aesCipher string
	aesPuk    string //AsePuk 是用于客户端与服务器之间通信
	aesPrk    string //AsePrk 是用于服务器之间通信

	rsaPuk *rsa.PublicKey
	rsaPrk *rsa.PrivateKey
}

var sInstance *EncMgr
var once sync.Once

func Instance() *EncMgr {
	once.Do(func() {
		sInstance = &EncMgr{}
		sInstance.init()
	})
	//beego.Debug("*****EncMgr  instance Done")
	return sInstance
}

func (o *EncMgr) init() {
	ret := Key{}
	err := dbmgr.Instance().EncryptCollection.Find(nil).One(&ret)
	if err != nil {
		beego.Error(err)
		return
	}

	o.rsaCipher = ret.Rsa_cipher
	o.aesCipher = ret.Aes_cipher
	o.aesPuk = ret.Aes_puk
	o.aesPrk = ret.Aes_prk

	//解析公钥
	rsaPuk, err := encrypt.DecodePuk(conf.RsaPuk)
	if err != nil {
		beego.Error(err)
		return
	}
	o.rsaPuk = rsaPuk

	rsaPrk, err := encrypt.DecodePrk(conf.RsaPrk)
	if err != nil {
		beego.Error(err)
		return
	}
	o.rsaPrk = rsaPrk

	beego.Info("--- Init Encrypt Mgr  Done !")
}

func (o EncMgr) RsaDec(cipher []byte) ([]byte, error) {
	return encrypt.RsaDecrypt(o.rsaPrk, cipher)
}

func (o EncMgr) RsaEnc(plaintext []byte) ([]byte, error) {
	return encrypt.RsaEncrypt(o.rsaPuk, plaintext)
}

func (o EncMgr) GetRsaCipher() string {
	return o.rsaCipher
}

func (o EncMgr) GetAesCipher() string {
	return o.aesCipher
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

func (o EncMgr) GetAesPuk() string {
	return o.aesPuk
}

func (o EncMgr) GetAesPrk() string {
	return o.aesPrk
}

//目前没有找到再mongodb中存贮 RSA KEY 的方法等功能做完来找
type Key struct {
	Rsa_cipher string
	Aes_cipher string

	Aes_puk string
	Aes_prk string
}
