package configs

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type ConfigFiber struct {
	Type   string `yaml:"type" env-default:porn`
	BindIp string `yaml:"bind_ip" env-default:0.0.0.0`
	Port   string `yaml:"port" env-default:8081`
}

type ConfigMongoDB struct {
	Host string `yaml:"host" env-default:localhost`
	Port string `yaml:"port" env-default:27017`
	Database string `yaml:"database"`
	Collection string `yaml:"collection"`
}

type Config struct {
	ModeLog string `yaml: modelog env-default: jsonInfo`
	Fiber  ConfigFiber `yaml: fiber`
	MongoDB ConfigMongoDB  `yaml:"mongodb"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("config_path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not open: %s", configPath)
	}
	once.Do(func () {
		cfg = &Config{}
		
		if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
			cleanenv.GetDescription(cfg, nil)
		}

		if mongoHost := os.Getenv("MONGODB_HOST"); mongoHost != "" {
			cfg.MongoDB.Host = mongoHost
		}

	})
	return cfg
}