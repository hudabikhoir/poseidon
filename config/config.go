package config

import (
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

//AppConfig Application configuration
type AppConfig struct {
	Port       int    `yaml:"port"`
	SelectedDB string `yaml:"selecteddb"`
	Database   struct {
		SQL struct {
			Driver   string `yaml:"driver"`
			Name     string `yaml:"name"`
			Address  string `yaml:"address"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"sql"`
		NOSQL struct {
			Driver   string `yaml:"driver"`
			Name     string `yaml:"name"`
			Address  string `yaml:"address"`
			Port     int    `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"nosql"`
	} `yaml:"database"`
	Cache struct {
		Driver   string `yaml:"driver"`
		Address  string `yaml:"address"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DBNumber int    `yaml:"dbnumber"`
	}
	Endpoint struct {
		Auth            string `yaml:"auth"`
		Commodities     string `yaml:"commodities"`
		ConvertCurrency string `yaml:"convertcurrency"`
		FFMSDB          string `yaml:"ffmsdb"`
	}
	Token struct {
		FFMSDB string `yaml:"ffmsdb"`
	}
}

var lock = &sync.Mutex{}
var appConfig *AppConfig

//GetConfig Initiatilize config in singleton way
func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func initConfig() *AppConfig {
	var defaultConfig AppConfig
	defaultConfig.Port = 5001
	defaultConfig.SelectedDB = "postgress"
	defaultConfig.Database.SQL.Driver = "sqlite"
	defaultConfig.Database.SQL.Name = "insanulab-test"
	defaultConfig.Database.SQL.Address = ""
	defaultConfig.Database.SQL.Port = 3306
	defaultConfig.Database.SQL.Username = ""
	defaultConfig.Database.SQL.Password = ""

	defaultConfig.Database.NOSQL.Driver = "sqlite"
	defaultConfig.Database.NOSQL.Name = "insanulab-test"
	defaultConfig.Database.NOSQL.Address = ""
	defaultConfig.Database.NOSQL.Port = 3306
	defaultConfig.Database.NOSQL.Username = ""
	defaultConfig.Database.NOSQL.Password = ""

	defaultConfig.Cache.Driver = "redis"
	defaultConfig.Cache.Address = "localhost"
	defaultConfig.Cache.Port = 6379
	defaultConfig.Cache.DBNumber = 1
	defaultConfig.Cache.Password = ""

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		log.Info("error to load config file, will use default value ", err)
		return &defaultConfig
	}

	var finalConfig AppConfig
	err := viper.Unmarshal(&finalConfig)
	if err != nil {
		log.Info("failed to extract config, will use default value ")
		return &defaultConfig
	}

	return &finalConfig
}
