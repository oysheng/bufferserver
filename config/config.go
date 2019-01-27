package config

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

func NewConfig() *Config {
	if len(os.Args) <= 1 {
		log.Fatal("Please setup the config file path")
	}

	return NewConfigWithPath(os.Args[1])
}

func NewConfigWithPath(path string) *Config {
	configFile, err := os.Open(path)
	if err != nil {
		log.WithFields(log.Fields{"err": err, "file_path": os.Args[1]}).Fatal("fail to open config file")
	}
	defer configFile.Close()

	cfg := &Config{}
	if err := json.NewDecoder(configFile).Decode(cfg); err != nil {
		log.WithField("err", err).Fatal("fail to decode config file")
	}

	return cfg
}

type Config struct {
	GinGonic GinGonic `json:"gin-gonic"`
	MySQL    MySQL    `json:"mysql"`
	Redis    Redis    `json:"redis"`
	API      API      `json:"api"`
	Updater  Updater  `json:"updater"`
}

type GinGonic struct {
	ListeningPort uint64 `json:"listening_port"`
	IsReleaseMode bool   `json:"is_release_mode"`
}

type MySQL struct {
	Master  MySQLConnection `json:"master"`
	LogMode bool            `json:"log_mode"`
}

type MySQLConnection struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DbName   string `json:"database"`
}

type Redis struct {
	Endpoint       string `json:"endpoint"`
	PoolSize       int    `json:"pool_size"`
	Database       int    `json:"database"`
	Password       string `json:"password"`
	CacheSeconds   int    `json:"cache_seconds"`
	LongCacheHours int    `json:"long_cache_hours"`
}

type MySQLConnCfg struct {
	MaxOpenConns int `json:"max_open_conns"`
	MaxIdleConns int `json:"max_idle_conns"`
}

type API struct {
	MySQLConnCfg MySQLConnCfg `json:"mysql_conns"`
}

type Updater struct {
	SyncSeconds  int          `json:"sync_seconds"`
	URL          string       `json:"url"`
	MySQLConnCfg MySQLConnCfg `json:"mysql_conns"`
}
