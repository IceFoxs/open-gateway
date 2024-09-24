package dubbo

import (
	"context"
	dcn "dubbo.apache.org/dubbo-go/v3/common/constant"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/IceFoxs/open-gateway/util"
	"github.com/dubbogo/gost/log/logger"
	"sync"
)
import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance/ringhash"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)
import (
	dg "dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	hessian "github.com/apache/dubbo-go-hessian2"
)

var (
	dubboClient        *Client
	onceClient         = sync.Once{}
	defaultApplication = &dg.ApplicationConfig{
		Organization: constant.SERVEICE_NAME,
		Name:         constant.SERVEICE_NAME,
		Module:       constant.SERVEICE_NAME,
		Version:      "",
		Owner:        constant.SERVEICE_NAME,
		Environment:  "pro",
	}
)

// Client client to generic invoke dubbo
type Client struct {
	lock               sync.RWMutex
	GenericServicePool map[string]*generic.GenericService
	rootConfig         *dg.RootConfig
}

// SingletonDubboClient singleton dubbo client
func SingletonDubboClient() *Client {
	if dubboClient == nil {
		onceClient.Do(func() {
			InitDefaultDubboClient()
		})
	}
	return dubboClient
}

// InitDefaultDubboClient init default dubbo client
func InitDefaultDubboClient() {
	dubboClient = NewDubboClient()
	if err := dubboClient.Apply(); err != nil {
		logger.Warnf("dubbo client apply error %s", err)
	}
}

func NewDubboClient() *Client {
	return &Client{
		lock:               sync.RWMutex{},
		GenericServicePool: make(map[string]*generic.GenericService, 4),
	}
}

// Close clear GenericServicePool.
func (dc *Client) Close() error {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	for k := range dc.GenericServicePool {
		delete(dc.GenericServicePool, k)
	}
	return nil
}

// Apply init dubbo, config mapping can do here
func (dc *Client) Apply() error {
	registry := conf.GetConf().Registry.Register
	registryType := conf.GetConf().Registry.RegisterType
	address := conf.GetConf().Registry.RegistryAddress[0]
	username := conf.GetConf().Registry.Username
	password := conf.GetConf().Registry.Password
	registryConfig := &dg.RegistryConfig{
		Protocol:     registry,
		Address:      address,
		Username:     username,
		Password:     password,
		Simplified:   true,
		RegistryType: registryType,
	}
	metadata := &dg.MetadataReportConfig{
		Protocol: registry,
		Address:  address,
		Username: username,
		Password: password,
	}
	rootConfig := dg.NewRootConfigBuilder().
		SetApplication(defaultApplication).
		AddRegistry(registry, registryConfig).
		SetMetadataReport(metadata).
		Build()
	if err := dg.Load(dg.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	dc.rootConfig = rootConfig
	return nil
}
func (dc *Client) get(key string) *generic.GenericService {
	dc.lock.RLock()
	defer dc.lock.RUnlock()
	return dc.GenericServicePool[key]
}
func (dc *Client) check(key string) bool {
	dc.lock.RLock()
	defer dc.lock.RUnlock()
	if _, ok := dc.GenericServicePool[key]; ok {
		return true
	}
	return false
}

// Get find a dubbo GenericService
func (dc *Client) Get(key string, iface string) *generic.GenericService {
	logger.Infof("key :%s,iface: %s", key, iface)
	if dc.check(key) {
		return dc.get(key)
	}
	return dc.create(key, iface)
}
func (dc *Client) create(key string, iface string) *generic.GenericService {
	check := false
	registry := conf.GetConf().Registry.Register
	retryNum := conf.GetConf().Registry.Retries
	requestTimeout := conf.GetConf().Registry.RequestTimeout
	refConf := dg.ReferenceConfig{
		InterfaceName:  iface,
		RegistryIDs:    []string{registry},
		Protocol:       dubbo.DUBBO,
		Generic:        "true",
		Retries:        retryNum,
		RequestTimeout: requestTimeout,
		Check:          &check,
	}
	logger.Debugf("[opengateway] client dubbo timeout val %v", refConf.RequestTimeout)
	dc.lock.Lock()
	defer dc.lock.Unlock()
	if service, ok := dc.GenericServicePool[key]; ok {
		return service
	}
	if err := dg.Load(dg.WithRootConfig(dc.rootConfig)); err != nil {
		panic(err)
	}
	_ = refConf.Init(dc.rootConfig)
	refConf.GenericLoad(iface)
	clientService := refConf.GetRPCService().(*generic.GenericService)
	dc.GenericServicePool[key] = clientService
	return clientService
}
func (dc *Client) Invoke(gmm model.GatewayMethodMetadata, param interface{}, wrapResp string) (interface{}, error) {
	paramTypes := []string{}
	params := []hessian.Object{}
	if len(gmm.ParameterTypeName) != 0 {
		paramTypes = append(paramTypes, gmm.ParameterTypeName)
		params = append(params, param)
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, dcn.AttachmentKey, map[string]interface{}{
		constant.WRAP_RESP_HEADER: wrapResp,
	})
	gs := dc.Get(gmm.GetReferenceKey(), gmm.InterfaceName)
	resp, err := gs.Invoke(
		ctx,
		gmm.MethodName,
		paramTypes,
		params,
	)
	//if err != nil {
	//	return nil, err
	//}
	logger.Infof("Invoke method,%s res: %+v", gmm.GatewayMethodName, resp)
	if resp != nil {
		return util.DealResp(resp, false)
	}
	return nil, err
}
