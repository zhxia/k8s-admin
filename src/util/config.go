package util

import (
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strings"
)

type (
	Config struct {
		ServerHost   string                  `json:"server_host"`
		ServerPort   int                     `json:"server_port"`
		K8sNamespace string                  `json:"k8s_namespace"`
		K8sConfig    string                  `json:"k8s_config"`
		LogLevel     log.Level               `json:"log_level"`
		LogDir       string                  `json:"log_dir"`
		RedisHost    string                  `json:"redis_host"`
		RedisPort    int                     `json:"redis_port"`
		RedisPass    string                  `json:"redis_pass"`
		AuthUsers    map[string]gjson.Result `json:"auth_users"`
	}
)

func (conf *Config) getLogLevel(strLevel string) log.Level {
	strLevel = strings.ToUpper(strLevel)
	level := log.InfoLevel
	switch strLevel {
	case "TRACE":
		level = log.TraceLevel
	case "DEBUG":
		level = log.DebugLevel
	case "INFO":
		level = log.InfoLevel
	case "WARN":
		level = log.WarnLevel
	case "ERROR":
		level = log.ErrorLevel
	case "FATAL":
		level = log.FatalLevel
	case "PANIC":
		level = log.PanicLevel
	}
	return level
}
func NewConfig(configFile string) *Config {
	config := &Config{}
	jsonConf := YamlFileToJson(configFile)
	if jsonConf.Exists() {
		logLevel := jsonConf.Get("log.level").String()
		logDir := jsonConf.Get("log.dir").String()
		config.ServerHost = jsonConf.Get("server.host").String()
		config.ServerPort = int(jsonConf.Get("server.port").Uint())
		config.K8sNamespace = jsonConf.Get("k8s.namespace").String()
		config.K8sConfig = jsonConf.Get("k8s.config").String()
		config.LogDir = logDir
		config.LogLevel = config.getLogLevel(logLevel)
		config.RedisHost = jsonConf.Get("redis.host").String()
		config.RedisPort = int(jsonConf.Get("redis.port").Uint())
		config.RedisPass = jsonConf.Get("redis.pass").String()
		config.AuthUsers = jsonConf.Get("auth.users").Map()
	}
	return config
}
