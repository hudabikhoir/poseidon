package config

import (
	"sync"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

//AppConfig Application configuration
type AppConfig struct {
	App struct {
		Name string `toml:"name"`
		Env  string `toml:"env"`
		Port int    `toml:"port"`
	} `toml:"app"`
	Database struct {
		Driver   string `toml:"driver"`
		Address  string `toml:"address"`
		Port     int    `toml:"port"`
		Username string `toml:"username"`
		Password string `toml:"password"`
		Name     string `toml:"name"`
	} `toml:"database"`
	Cache struct {
		Driver   string `toml:"driver"`
		Address  string `toml:"address"`
		Port     int    `toml:"port"`
		Password string `toml:"password"`
		Dbnumber int    `toml:"dbnumber"`
	} `toml:"cache"`
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
	// var defaultConfig AppConfig
	// defaultConfig.App.Port = 5001
	// defaultConfig.App.Name = "poseidon"
	// defaultConfig.App.Env = "local"

	// defaultConfig.Database.Driver = "sqlite"
	// defaultConfig.Database.Name = "insanulab-test"
	// defaultConfig.Database.Address = ""
	// defaultConfig.Database.Port = 3306
	// defaultConfig.Database.Username = ""
	// defaultConfig.Database.Password = ""

	// defaultConfig.Cache.Driver = "redis"
	// defaultConfig.Cache.Address = "localhost"
	// defaultConfig.Cache.Port = 6379
	// defaultConfig.Cache.Password = ""

	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		log.Info("error to load config file, will use default value ", err)
		// return &defaultConfig
	}

	var finalConfig AppConfig
	err := viper.Unmarshal(&finalConfig)
	if err != nil {
		log.Info("failed to extract config, will use default value ")
		// return &defaultConfig
	}

	return &finalConfig
}
