package rsa_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/IceFoxs/open-gateway/util/aes"
	rsaUtil "github.com/IceFoxs/open-gateway/util/rsa"
	"testing"
)

func TestRsa(*testing.T) {
	// 使用示例
	sign := &rsaUtil.SignAlgo{}

	param := map[string]interface{}{
		"a":    1,
		"b":    "b",
		"c":    2.0,
		"d":    "d",
		"sign": "xxx",
	}

	secret := []byte("1234567890123456")

	sortedParams := sign.SortParam(param)
	println("Sorted Params:", sortedParams)

	signature := sign.Sign(param, secret)
	println("Signature:", signature)

	param["sign"] = signature
	for key, value := range param {
		fmt.Println("Key:", key, "Value:", value)
	}
	verified := sign.Verify(param, secret)
	println("Verified:", verified)
	bytes, err := json.Marshal(param)
	if err != nil {
		return
	}
	encrypted := aes.AesEncryptECB(bytes, secret)
	println("Encrypted:", encrypted)

	decrypted := aes.AesDecryptECB(encrypted, secret)

	println("Decrypted:", string(decrypted))

	// threading.GoSafe()

	// 生成 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Failed to generate RSA private key:", err)
		return
	}
	publicKey := &privateKey.PublicKey

	// 使用公钥进行加密
	RsaEncrypted, _ := rsaUtil.RSAEncrypt(param, publicKey)
	println("RSAencrypted:", RsaEncrypted)

	// 使用私钥进行解密
	RsaDecrypted, _ := rsaUtil.RSADecrypt(RsaEncrypted, privateKey)
	println("RsaDecrypted:", string(RsaDecrypted))

	// 使用私钥进行签名
	RsaSignature, _ := rsaUtil.RSASign(param, privateKey)
	println("RsaSignature:", RsaSignature)
	param["sign"] = RsaSignature
	for key, value := range param {
		fmt.Println("Key:", key, "Value:", value)
	}
	// 验证签名的有效性
	err = rsaUtil.RSAVerify(param, publicKey)

	fmt.Println("Signature verificatio:", err)

}
