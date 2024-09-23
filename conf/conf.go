package conf

import (
	"fmt"
	"github.com/dubbogo/gost/log/logger"
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
	Host       string `yaml:"host"`
	Name       string `yaml:"name"`
	StaticPath string `yaml:"static_path"`
	Register   string `yaml:"register"`
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
	logger.Infof("confFileRelPath - %v", confFileRelPath)
	content, err := ioutil.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}
	conf = new(Config)
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		logger.Errorf("parse yaml error - %v", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		logger.Errorf("validate config error - %v", err)
		panic(err)
	}
	conf.Env = GetEnv()
}

func GetEnv() string {
	e := os.Getenv("GO_ENV")
	if len(e) == 0 {
		return "test"
	}
	return e
}
