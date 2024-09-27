package model

type GatewayConfig struct {
	AppId               string `json:"appId" column:"app_id" gorm:"primarykey"`
	AppName             string `json:"appName" column:"app_name"`
	AesKey              string `json:"aesKey" column:"aes_key"`
	AesType             string `json:"aesType" column:"aes_type"`
	RsaPublicKey        string `json:"rsaPublicKey" column:"rsa_public_key"`
	RsaPrivateKey       string `json:"rsaPrivateKey" column:"rsa_private_key"`
	SignType            string `json:"signType" column:"sign_type"`
	IsEnable            int    `json:"isEnable" column:"is_enable"`
	ClientRsaPrivateKey string `json:"clientRsaPrivateKey" column:"client_rsa_private_key"`
	ClientRsaPublicKey  string `json:"clientRsaPublicKey" column:"client_rsa_public_key"`
	CallbackUrl         string `json:"callbackUrl" column:"callback_url"`
}

func (c *GatewayConfig) TableName() string {
	return "gateway_channel_config"
}
