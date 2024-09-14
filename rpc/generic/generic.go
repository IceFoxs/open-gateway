package generic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IceFoxs/open-gateway/util"
	"log"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/dubbogo/gost/log/logger"

	tpconst "github.com/dubbogo/triple/pkg/common/constant"
)

const appName = "dubbo.io"

// export DUBBO_GO_CONFIG_PATH= PATH_TO_SAMPLES/generic/default/go-client/conf/dubbogo.yml
func main() {
	// register POJOs
	// generic invocation samples using hessian serialization on Dubbo protocol
	dubboRefConf := NewRefConf("org.apache.dubbo.samples.UserProvider", dubbo.DUBBO)
	CallGetUser(dubboRefConf)
	//callGetOneUser(dubboRefConf)
	CallGetUsers(dubboRefConf)
	CallGetUsersMap(dubboRefConf)
	CallQueryUser(dubboRefConf)
	CallQueryUsers(dubboRefConf)
	//callQueryAll(dubboRefConf)

	// generic invocation samples using hessian serialization on Triple protocol
	tripleRefConf := NewRefConf("org.apache.dubbo.samples.UserProviderTriple", tpconst.TRIPLE)
	CallGetUser(tripleRefConf)
	//callGetOneUser(tripleRefConf)
	CallGetUsers(tripleRefConf)
	CallGetUsersMap(tripleRefConf)
	CallQueryUser(tripleRefConf)
	CallQueryUsers(tripleRefConf)
	//callQueryAll(tripleRefConf)

}

func CallGetUser(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser1",
		[]string{"java.lang.String"},
		[]hessian.Object{"A003"},
	)

	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser1(userId string) res: %+v", resp)

	resp, err = refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser2",
		[]string{"java.lang.String", "java.lang.String"},
		[]hessian.Object{"A003", "lily"},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser2(userId string, name string) res: %+v", resp)

	resp, err = refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser3",
		[]string{"int"},
		[]hessian.Object{1},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser3(userCode int) res: %+v", resp)

	resp, err = refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUser4",
		[]string{"int", "java.lang.String"},
		[]hessian.Object{1, "zhangsan"},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUser4(userCode int, name string) res: %+v", resp)
}

// nolint
func CallGetOneUser(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetOneUser",
		[]string{},
		// TODO go-go []hessian.Object{}, go-java []string{}
		[]hessian.Object{},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetOneUser() res: %+v", resp)
}

func CallGetUsers(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUsers",
		[]string{"java.util.List"},
		[]hessian.Object{
			[]hessian.Object{
				"001", "002", "003", "004",
			},
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUsers1(userIdList []*string) res: %+v", resp)
}

func CallGetUsersMap(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"GetUsersMap",
		[]string{"java.util.List"},
		[]hessian.Object{
			[]hessian.Object{
				"001", "002", "003", "004",
			},
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("GetUserMap(userIdList []*string) res: %+v", resp)
}

func CallQueryUser(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryUser",
		[]string{"org.apache.dubbo.samples.User"},
		// the map represents a User object:
		// &User {
		// 		ID: "3213",
		// 		Name: "panty",
		// 		Age: 25,
		// 		Time: time.Now(),
		// }
		[]hessian.Object{
			map[string]hessian.Object{
				"iD":   "3213",
				"name": "panty",
				"age":  25,
				"time": time.Now(),
			}},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("queryUser(user *User) res: %+v", resp)
}

func CallQueryUsers(refConf config.ReferenceConfig) {
	var resp, err = refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryUsers",
		[]string{"java.util.ArrayList"},
		[]hessian.Object{
			[]hessian.Object{
				map[string]hessian.Object{
					"id":    "3212",
					"name":  "XavierNiu",
					"age":   24,
					"time":  time.Now().Add(4),
					"class": "org.apache.dubbo.samples.User",
				},
				map[string]hessian.Object{
					"iD":    "3213",
					"name":  "zhangsan",
					"age":   21,
					"time":  time.Now().Add(4),
					"class": "org.apache.dubbo.samples.User",
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("queryUsers(users []*User) res: %+v", resp)
}

// nolint
func CallQueryAll(refConf config.ReferenceConfig) {
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"queryAll",
		[]string{},
		// TODO go-go []hessian.Object{}, go-java []string{}
		//[]hessian.Object{},
		[]hessian.Object{},
	)
	if err != nil {
		panic(err)
	}
	logger.Infof("queryAll() res: %+v", resp)
}

func NewRefConf(iface, protocol string) config.ReferenceConfig {
	registryConfig := &config.RegistryConfig{
		Protocol: "zookeeper",
		Address:  "127.0.0.1:2181",
	}

	refConf := config.ReferenceConfig{
		InterfaceName: iface,
		Cluster:       "failover",
		RegistryIDs:   []string{"zk"},
		Protocol:      protocol,
		Generic:       "true",
	}

	rootConfig := config.NewRootConfigBuilder().
		AddRegistry("zk", registryConfig).
		Build()
	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	_ = refConf.Init(rootConfig)
	refConf.GenericLoad(appName)
	return refConf
}

func NewRefConf1(iface, registry string, registryType string, protocol string, address string, username string, password string) config.ReferenceConfig {
	registryConfig := &config.RegistryConfig{
		Protocol:     registry,
		Address:      address,
		Username:     username,
		Password:     password,
		Simplified:   true,
		RegistryType: registryType,
	}
	refConf := config.ReferenceConfig{
		InterfaceName: iface,
		RegistryIDs:   []string{registry},
		Protocol:      protocol,
		Generic:       "true",
	}
	metadata := &config.MetadataReportConfig{
		Protocol: registry,
		Address:  address,
		Username: username,
		Password: password,
	}
	app := &config.ApplicationConfig{}
	rootConfig := config.NewRootConfigBuilder().
		SetApplication(app).
		AddRegistry(registry, registryConfig).
		SetMetadataReport(metadata).
		Build()
	if err := config.Load(config.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	_ = refConf.Init(rootConfig)
	refConf.GenericLoad(iface)
	return refConf
}

func ConfRefresh(refConf config.ReferenceConfig) (interface{}, error) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recover the panic:", e)
		}
	}()
	resp, err := refConf.GetRPCService().(*generic.GenericService).Invoke(
		context.TODO(),
		"confRefresh",
		[]string{"com.hundsun.manager.model.req.ConfRefreshRequest"},
		[]hessian.Object{
			map[string]hessian.Object{
				"confType":    "BANK_TEST",
				"confContent": "TEST|20240930",
			}},
	)
	if err != nil {

		return nil, err
	}
	logger.Infof("confRefresh res: %+v", resp)
	data := util.ConvertMap(resp)
	by, err := json.Marshal(data)
	log.Println("output json:", string(by), err)
	return data, nil
}
