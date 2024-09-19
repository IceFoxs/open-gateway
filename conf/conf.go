package conf

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/kr/pretty"
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

type Config struct {
	Env      string
	App      App      `yaml:"app"`
	MySQL    MySQL    `yaml:"mysql"`
	Redis    Redis    `yaml:"redis"`
	Registry Registry `yaml:"registry"`
}
type App struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Register string `yaml:"register"`
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

type Registry struct {
	RegistryAddress []string `yaml:"registry_address"`
	Username        string   `yaml:"username"`
	Password        string   `yaml:"password"`
	RegisterType    string   `yaml:"register_type"`
}

// GetConf gets configuration instance
func GetConf() *Config {
	once.Do(initConf)
	return conf
}

func initConf() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dir)
	prefix := "config"
	confFileRelPath := dir + "/" + filepath.Join(prefix, filepath.Join(GetEnv(), "conf.yaml"))
	hlog.SystemLogger().Infof("confFileRelPath - %v", confFileRelPath)
	content, err := ioutil.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		hlog.SystemLogger().Errorf("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		hlog.SystemLogger().Errorf("validate config error - %v", err)
		panic(err)
	}
	conf.Env = GetEnv()
	pretty.Printf("%+v\n", conf)
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}
