package model

type GatewayMethodMetadata struct {
	GatewayMethodName string `json:"gatewayMethodName"`
	InterfaceName     string `json:"interfaceName"`
	MethodName        string `json:"methodName"`
	ParameterTypeName string `json:"parameterTypeName"`
	RpcType           string `json:"rpcType"`
	Path              string `json:"path"`
	ContextPath       string `json:"contextPath"`
	AppName           string `json:"appName"`
}
