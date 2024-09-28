package conf

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	conf *Config
	once sync.Once
)

type Logger struct {
	FileName   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	//json console
	Encoding string `yaml:"encoding"`
}
type Config struct {
	Env        string
	App        App        `yaml:"app"`
	SyncConfig SyncConfig `yaml:"sync_config"`
	MySQL      MySQL      `yaml:"mysql"`
	Redis      Redis      `yaml:"redis"`
	Dubbo      Dubbo      `yaml:"dubbo"`
	Discovery  Discovery  `yaml:"discovery"`
	Registry   Registry   `yaml:"registry"`
	Logger     Logger     `yaml:"logger"`
	BaseDir    string     `yaml:"base_dir"`
	Zookeeper  Zookeeper  `yaml:"zookeeper"`
	Nacos      Nacos      `yaml:"nacos"`
	Consul     Consul     `yaml:"consul"`
}
type SyncConfig struct {
	ConfigType string `yaml:"config_type"`
}

type Nacos struct {
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

type Consul struct {
	Address []string `yaml:"address"`
	Token   string   `yaml:"token"`
	Schema  string   `yaml:"schema"`
}
type Zookeeper struct {
	Address        []string `yaml:"address"`
	SessionTimeout int64    `yaml:"session_timeout"`
	Username       string   `yaml:"username"`
	Password       string   `yaml:"password"`
}

type App struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Register string `yaml:"register"`
	Weight   int    `yaml:"weight"`
	WebPath  string `yaml:"web_path"`
}
type MySQL struct {
	DSN string `yaml:"dsn"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}
type Discovery struct {
	HttpRegister string `yaml:"http_register"`
}
type Registry struct {
	Register string `yaml:"register"`
}

type Dubbo struct {
	RegisterType   string `yaml:"register_type"`
	Register       string `yaml:"register"`
	WrapResp       string `yaml:"wrap_resp"`
	Retries        string `yaml:"retries"`
	RequestTimeout string `yaml:"request_timeout"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	confPath := os.Getenv("CONF_PATH")
	var confFileRelPath string
	var dir string
	if len(confPath) == 0 {
		dir = os.Getenv("BASE_DIR")
		if len(dir) == 0 {
			var err error
			dir, err = os.Getwd()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		hlog.Infof("BASE_DIR:%s", dir)
		prefix := "config"
		confFileRelPath = dir + "/" + filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	} else {
		hlog.Infof("CONF_PATH:%s", confPath)
		confFileRelPath = confPath
	}
	hlog.Infof("confFileRelPath - %v", confFileRelPath)
	content, err := ioutil.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		hlog.Errorf("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		hlog.Errorf("validate config error - %v", err)
		panic(err)
	}
	conf.Env = GetEnv()
	if len(conf.BaseDir) == 0 {
		conf.BaseDir = dir
	}
}
func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}
