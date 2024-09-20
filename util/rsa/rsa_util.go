package rsa

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	// "github.com/zeromicro/go-zero/core/threading"
)

type SignAlgo struct{}

func (s *SignAlgo) SortParam(param map[string]interface{}) string {
	keys := make([]string, 0, len(param))
	for key := range param {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var sortedParams []string
	for _, key := range keys {

		value := param[key]

		if key == "sign" || value == nil || value == "" {
			continue
		}

		switch v := value.(type) {
		case int, uint, int16, int32, int64:
			sortedParams = append(sortedParams, fmt.Sprintf("%s=%d", key, v))
		case float64, float32:
			sortedParams = append(sortedParams, fmt.Sprintf("%s=%f", key, v))
		default:
			sortedParams = append(sortedParams, key+"="+value.(string))
		}

	}
	return strings.Join(sortedParams, "&")
}

func (s *SignAlgo) Sign(data map[string]interface{}, secret []byte) string {
	str := s.SortParam(data)
	str += "&key=" + string(secret)
	hash := md5.Sum(secret)
	return hex.EncodeToString(hash[:])
}

func (s *SignAlgo) Verify(data map[string]interface{}, secret []byte) bool {
	signature := data["sign"].(string)
	delete(data, "sign")
	return s.Sign(data, secret) == signature
}

// PKCS5Padding 对数据进行PKCS5填充
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS5Unpadding 去除PKCS5填充
func PKCS5Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("PKCS5 unpadding error: data is empty")
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return nil, errors.New("PKCS5 unpadding error: invalid padding size")
	}
	return data[:length-unpadding], nil
}

// ZeroPadding 使用ZeroPadding填充数据
func ZeroPadding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{0}, padding)
	return append(data, padText...)
}

// ZeroUnpadding 去除ZeroPadding填充数据
func ZeroUnpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("ZeroUnpadding error: data is empty")
	}
	unpadding := 0
	for i := length - 1; i >= 0; i-- {
		if data[i] == 0 {
			unpadding++
		} else {
			break
		}
	}
	if unpadding == 0 {
		return nil, errors.New("ZeroUnpadding error: no padding bytes found")
	}
	return data[:length-unpadding], nil
}

// RSAEncrypt 使用 RSA 公钥对数据进行加密
func RSAEncrypt(data map[string]interface{}, publicKey *rsa.PublicKey) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	encryptData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptData), nil
}

// RSADecrypt 使用 RSA 私钥对数据进行解密
func RSADecrypt(ciphertext string, privateKey *rsa.PrivateKey) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, decoded)
}

// RSASign 使用 RSA 私钥对数据进行签名
func RSASign(data map[string]interface{}, privateKey *rsa.PrivateKey) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	hashed := sha256.Sum256(bytes)
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

func RSASignByString(data string, privateKey *rsa.PrivateKey) (string, error) {
	hashed := sha256.Sum256([]byte(data))
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

// RSAVerify 验证 RSA 签名的有效性
func RSAVerify(data map[string]interface{}, publicKey *rsa.PublicKey) error {
	signature := data["sign"].(string)
	decoded, err := base64.StdEncoding.DecodeString(signature)
	delete(data, "sign")
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256(bytes)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], decoded)
}

// RSAVerify 验证 RSA 签名的有效性
func RSAVerifyByString(data string, signature string, publicKey *rsa.PublicKey) error {
	decoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(data))
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], decoded)
}

func Base64PublicKeyToRSA(base64PublicKey string) (*rsa.PublicKey, error) {
	// 解码Base64字符串
	publicKeyDecoded, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		return nil, err
	}
	// 解析PEM格式的公钥
	pubInterface, err := x509.ParsePKIXPublicKey(publicKeyDecoded)
	if err != nil {
		return nil, err
	}
	// 类型断言，确保得到的是*rsa.PublicKey
	publicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key type")
	}

	return publicKey, nil
}

func Base64ToPrivateKey(base64EncodedKey string) (*rsa.PrivateKey, error) {
	// 解码Base64编码的密钥
	decodedKey, err := base64.StdEncoding.DecodeString(base64EncodedKey)
	if err != nil {
		return nil, err
	}

	// 解析PEM格式的密钥为*rsa.PrivateKey
	privateKey, err := x509.ParsePKCS1PrivateKey(decodedKey)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func Base64ToPrivateKeyByPkcs8(base64EncodedKey string) (*rsa.PrivateKey, error) {
	// 解码Base64编码的密钥
	decodedKey, err := base64.StdEncoding.DecodeString(base64EncodedKey)
	if err != nil {
		return nil, err
	}
	// 解析PEM格式的密钥为*rsa.PrivateKey
	privateKey, err := x509.ParsePKCS8PrivateKey(decodedKey)
	if err != nil {
		return nil, err
	}
	//// 类型断言，确保得到的是*rsa.PublicKey
	privateKey1, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key type")
	}

	return privateKey1, nil
}

func PrivateKeyToBytes(privateKey *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(privateKey)
}

func PublicKeyToBytes(publicKey *rsa.PublicKey) []byte {
	x, _ := x509.MarshalPKIXPublicKey(publicKey)
	return x
}

func PrivateKeyPkcs8ToBytes(privateKey *rsa.PrivateKey) []byte {
	cs8, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	return cs8
}
