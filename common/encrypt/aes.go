package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

//这个测试完成后 试试更简单的方法有一篇收藏在浏览器里面的帖子
//空了来研究加密

//由于目前项目使用的是json传输数据,而json.Marshal之后是byte所以这里可以直接写出byte[]的接口
//秘钥长度需要时AES-128(16bytes)或者AES-256(32bytes) 目前采用16bytes
//原文必须填充至blocksize的整数倍，填充方法可以参见https://tools.ietf.org/html/rfc5246#section-6.2.3.2
func AesEnc(key string, plaintext []byte) ([]byte, error) {
	//beego.Debug("原文原本长度 : ", len(plaintext))
	//补齐原文
	plaintext = PKCS5Padding(plaintext, aes.BlockSize)
	//beego.Debug("补齐后原文长度 : ", len(plaintext))
	//beego.Debug("补齐后原文内容 : ")
	//fmt.Printf("%s\n", plaintext)
	//beego.Debug("-------")

	//检查plaintext长度是否合法
	if len(plaintext)%aes.BlockSize != 0 { //块大小在aes.BlockSize中定义
		return nil, errors.New("plaintext is not a multiple of the block size")
	}

	//将明文key 生成加密用的16位block
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}

	// 对IV有随机性要求，但没有保密性要求，所以常见的做法是将IV包含在加密文本当中
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	//随机一个block大小作为IV
	//采用不同的IV时相同的秘钥将会产生不同的密文，可以理解为一次加密的session
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// 谨记密文需要认证(i.e. by using crypto/hmac)
	//fmt.Printf("原文加密后%x\n", ciphertext)
	return ciphertext, nil
}

func AesDec(key string, ciphertext []byte) ([]byte, error) {
	keyByte := []byte(key)
	//ciphertext, _ := hex.DecodeString("59eaa875ea5d05bf01258a0e8bd963a3ac27d9819931190e64723aac51cd051f")

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks可以原地更新
	mode.CryptBlocks(ciphertext, ciphertext)

	//fmt.Printf("%s\n", ciphertext)

	//beego.Debug("未去补原文长度 : ", len(ciphertext))

	plaintext := PKCS5UnPadding(ciphertext)

	//beego.Debug("完成解密后原文长度 : ", len(plaintext))
	//beego.Debug("完成解密后原文 : ")
	//fmt.Printf("%s\n", plaintext)
	return plaintext, nil
}

//将原文补齐到 block size 的长度
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//原文去补
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
