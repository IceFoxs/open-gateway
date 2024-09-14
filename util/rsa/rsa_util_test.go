package rsa_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
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
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
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

	fmt.Println("Signature verification:", err)

	// 编码私钥为Base64
	// 将私钥转换为PKCS8编码
	encodedPublicKey := rsaUtil.PublicKeyToBytes(publicKey)
	// 将私钥编码为Base64字符串
	publicBase64 := base64.StdEncoding.EncodeToString(encodedPublicKey)
	fmt.Println("Public Key (Base64):", publicBase64)
	pk, _ := rsaUtil.Base64PublicKeyToRSA(publicBase64)
	// 使用公钥进行加密
	RsaEncrypted1, _ := rsaUtil.RSAEncrypt(param, pk)
	println("64-RSAencrypted:", RsaEncrypted1)

	// 编码私钥为Base64
	// 将私钥转换为PKCS8编码
	encodedPrivateKey := rsaUtil.PrivateKeyToBytes(privateKey)
	// 将私钥编码为Base64字符串
	privateKeyBase64 := base64.StdEncoding.EncodeToString(encodedPrivateKey)
	fmt.Println("Private Key (Base64):", privateKeyBase64)
	// 使用私钥进行解密
	pk2, _ := rsaUtil.Base64ToPrivateKey(privateKeyBase64)
	RsaDecrypted1, _ := rsaUtil.RSADecrypt(RsaEncrypted1, pk2)
	println("64RsaDecrypted:", string(RsaDecrypted1))

	// 编码私钥为Base64
	// 将私钥转换为PKCS8编码
	encodedPrivateKey1 := rsaUtil.PrivateKeyPkcs8ToBytes(privateKey)
	// 将私钥编码为Base64字符串
	privateKeyBase641 := base64.StdEncoding.EncodeToString(encodedPrivateKey1)
	fmt.Println("Pkcs8 Private Key (Base64):", privateKeyBase641)
	// 使用私钥进行解密
	pk3, _ := rsaUtil.Base64ToPrivateKeyByPkcs8(privateKeyBase641)
	RsaDecrypted2, _ := rsaUtil.RSADecrypt(RsaEncrypted1, pk3)
	println("Pkcs8 64RsaDecrypted:", string(RsaDecrypted2))

}
