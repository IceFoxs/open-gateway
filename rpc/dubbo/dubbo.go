package dubbo

import (
	"context"
	dcn "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	"encoding/json"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/model"
	"github.com/IceFoxs/open-gateway/util"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"path/filepath"
	"sync"
)
import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance/ringhash"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)
import (
	dg "dubbo.apache.org/dubbo-go/v3/config"
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
		//版本不能设置
		Version:     "",
		Owner:       constant.SERVEICE_NAME,
		Environment: "pro",
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
			dubboClient = NewDubboClient()
		})
	}
	return dubboClient
}

// InitDefaultDubboClient init default dubbo client
func InitDefaultDubboClient() {
	dubboClient = NewDubboClient()
	if err := dubboClient.Apply(); err != nil {
		hlog.Warnf("dubbo client apply error %s", err)
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
	var address string
	var username string
	var password string
	if registry == constant.REGISTRY_NACOS {
		address = conf.GetConf().Nacos.Address[0]
		username = conf.GetConf().Nacos.Username
		password = conf.GetConf().Nacos.Password
	}
	if registry == constant.REGISTRY_ZOOKEEPER {
		address = conf.GetConf().Zookeeper.Address[0]
		username = conf.GetConf().Zookeeper.Username
		password = conf.GetConf().Zookeeper.Password
	}
	registryConfig := &dg.RegistryConfig{
		Protocol:     registry,
		Address:      address,
		Username:     username,
		Password:     password,
		Simplified:   true,
		RegistryType: registryType,
		Params: map[string]string{
			"nacos.cacheDir": conf.GetConf().BaseDir + string(filepath.Separator) + "logs" + string(filepath.Separator) + "dubbo" + string(filepath.Separator) + "cache",
			"nacos.logDir":   conf.GetConf().BaseDir + string(filepath.Separator) + "logs" + string(filepath.Separator) + "dubbo" + string(filepath.Separator) + "log",
		},
	}
	metadata := &dg.MetadataReportConfig{
		Protocol: registry,
		Address:  address,
		Username: username,
		Password: password,
		Params: map[string]string{
			"nacos.cacheDir": conf.GetConf().BaseDir + string(filepath.Separator) + "logs" + string(filepath.Separator) + "metadata" + string(filepath.Separator) + "cache",
			"nacos.logDir":   conf.GetConf().BaseDir + string(filepath.Separator) + "logs" + string(filepath.Separator) + "metadata" + string(filepath.Separator) + "log",
		},
	}
	logger := dg.NewLoggerConfigBuilder().SetDriver("zap").
		SetLevel("info").
		SetFileName(conf.GetConf().BaseDir + string(filepath.Separator) + "logs" + string(filepath.Separator) + "dubbo.log").
		SetFileMaxAge(10).
		SetFileMaxSize(100).
		SetFileMaxBackups(5).
		SetFileCompress(true).
		SetFormat("text").
		SetAppender("file").
		Build()
	err := logger.Init()
	if err != nil {
		return err
	}
	rootConfig := dg.NewRootConfigBuilder().
		SetApplication(defaultApplication).
		AddRegistry(registry, registryConfig).
		SetMetadataReport(metadata).
		SetLogger(logger).
		SetCacheFile(conf.GetConf().BaseDir + string(filepath.Separator) + "logs" + string(filepath.Separator) + "dubbo").
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
	hlog.Infof("key :%s,iface: %s", key, iface)
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
	hlog.Debugf("[opengateway] client dubbo timeout val %v", refConf.RequestTimeout)
	dc.lock.Lock()
	defer dc.lock.Unlock()
	if service, ok := dc.GenericServicePool[key]; ok {
		return service
	}
	if err := dg.Load(dg.WithRootConfig(dc.rootConfig)); err != nil {
		panic(err)
	}
	_ = refConf.Init(dc.rootConfig)
	refConf.GenericLoad(key)
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
	var body interface{}
	if resp != nil {
		body, err = util.DealResp(resp, false)
		bodyStr, _ := json.Marshal(body)
		hlog.Infof("Invoke method,%s res: %s", gmm.GatewayMethodName, bodyStr)
		return body, err
	}
	return nil, err
}
