package cmd

import (
	"errors"
	"github.com/IceFoxs/open-gateway/cache/appmetadata"
	"github.com/IceFoxs/open-gateway/cache/gatewayconfig"
	"github.com/IceFoxs/open-gateway/cache/gatewaymethod"
	"github.com/IceFoxs/open-gateway/cache/gatewaysystem"
	"github.com/IceFoxs/open-gateway/conf"
	"github.com/IceFoxs/open-gateway/constant"
	"github.com/IceFoxs/open-gateway/db"
	"github.com/IceFoxs/open-gateway/registry"
	"github.com/IceFoxs/open-gateway/rpc/dubbo"
	"github.com/IceFoxs/open-gateway/rpc/http"
	"github.com/IceFoxs/open-gateway/server/consul"
	na "github.com/IceFoxs/open-gateway/server/nacos"
	"github.com/IceFoxs/open-gateway/server/router"
	zk "github.com/IceFoxs/open-gateway/server/zookeeper"
	"github.com/IceFoxs/open-gateway/sync"
	"github.com/IceFoxs/open-gateway/sync/config/nacos"
	"github.com/IceFoxs/open-gateway/sync/config/zookeeper"
	"github.com/cloudwego/hertz/pkg/app/server"
	re "github.com/cloudwego/hertz/pkg/app/server/registry"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func getEncoder(z *zap.Config) zapcore.Encoder {
	if z.Encoding == "json" {
		return zapcore.NewJSONEncoder(z.EncoderConfig)
	} else if z.Encoding == "console" {
		return zapcore.NewConsoleEncoder(z.EncoderConfig)
	}
	return nil
}

// getLogWriter get Lumberjack writer by LumberjackConfig
func getLogWriter(l *lumberjack.Logger) zapcore.WriteSyncer {
	return zapcore.AddSync(l)
}
func Start() {
	log := conf.GetConf().Logger
	// 提供压缩和删除
	zapLoggerEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	config := &zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Encoding:         log.Encoding,
		EncoderConfig:    zapLoggerEncoderConfig,
		OutputPaths:      []string{"stdout", log.FileName},
		ErrorOutputPaths: []string{"stderr", log.FileName},
	}
	lumberjackLogger := &lumberjack.Logger{
		Filename:   log.FileName,
		MaxSize:    log.MaxSize,
		MaxBackups: log.MaxBackups,
		MaxAge:     log.MaxAge,
		Compress:   true,
	}
	logger := hertzzap.NewLogger(hertzzap.WithZapOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			getEncoder(config),
			getLogWriter(lumberjackLogger),
			config.Level)
	})))
	hlog.SetLogger(logger)
	nacos.GetConfChangeClient()
	zookeeper.GetConfChangeClient()
	dubbo.InitDefaultDubboClient()
	dsn := conf.GetConf().MySQL.DSN
	db.Init(dsn)
	registry.GetRegisterClient()
	c := gatewayconfig.GetGatewayConfigCache()
	c.RefreshCache()
	gsc := gatewaysystem.GetGatewaySystemCache()
	gsc.RefreshCache()
	appNames := gsc.GetAllAppName()
	amc := appmetadata.GetAppMetadataCache()
	amc.RefreshCacheByAppName(appNames)
	//amc.AddListen()
	gmc := gatewaymethod.GetGatewayMethodCache()
	methods := amc.GetAllMethods()
	gmc.RefreshAllCache(methods)
	http.GetHttpClient()
	sync.GetConfChangeClientHelper()
	h, err := CreateServer()
	if err != nil {
		hlog.SystemLogger().Errorf("create server failed: %s", err.Error())
	}
	//pprof.Register(h)
	staticPath := conf.GetConf().BaseDir
	if len(staticPath) == 0 {
		staticPath, _ = os.Getwd()
	}
	hlog.Infof("static path is %s", staticPath)
	router.AddRouter(h, staticPath)
	h.Spin()

}

func CreateServer() (*server.Hertz, error) {
	register := conf.GetConf().App.Register
	if register == "" {
		panic("app register can not empty, please check your config")
	}
	appName := conf.GetConf().App.Name
	host := conf.GetConf().App.Host
	var r re.Registry
	var err error
	if register == constant.REGISTRY_NACOS {
		r, err = na.CreateRegistry()
	} else if register == constant.REGISTRY_CONSUL {
		r, err = consul.CreateRegistry()
	} else if register == constant.REGISTRY_ZOOKEEPER {
		r, err = zk.CreateRegistry()
	} else {
		return nil, errors.New("register[" + register + "]is not supported")
	}
	if err != nil {
		return nil, err
	}
	weight := conf.GetConf().App.Weight
	if weight <= 0 {
		weight = 1
	}
	h := server.Default(server.WithHostPorts(host), server.WithRegistry(r, &re.Info{
		ServiceName: appName,
		Addr:        utils.NewNetAddr("tcp", host),
		Weight:      weight,
		Tags:        nil,
	}))
	return h, nil
}
