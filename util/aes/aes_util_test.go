package aes_test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/IceFoxs/open-gateway/util/aes"
	"log"
	"testing"
)

func TestAes(t *testing.T) {
	origData := []byte("赵云涛")         // 待加密的数据
	key := []byte("1234567890123456") // 加密的密钥
	log.Println("原文：", string(origData))
	fmt.Println("------------------ CBC模式 --------------------")
	encrypted := aes.AesEncryptCBC(origData, key)
	fmt.Println("密文(hex)：", hex.EncodeToString(encrypted))
	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := aes.AesDecryptCBC(encrypted, key)
	fmt.Println("解密结果：", string(decrypted))
	fmt.Println("------------------ ECB模式 --------------------")
	encrypted = aes.AesEncryptECB(origData, key)
	fmt.Println("密文(hex)：", hex.EncodeToString(encrypted))
	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = aes.AesDecryptECB(encrypted, key)
	fmt.Println("解密结果：", string(decrypted))
	fmt.Println("------------------ CFB模式 --------------------")
	encrypted = aes.AesEncryptCFB(origData, key)
	fmt.Println("密文(hex)：", hex.EncodeToString(encrypted))
	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted = aes.AesDecryptCFB(encrypted, key)
	fmt.Println("解密结果：", string(decrypted))
}
