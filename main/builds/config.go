package builds

import (
	e "XrayHelper/main/errors"
	"XrayHelper/main/log"
	"os"

	"gopkg.in/yaml.v3"
)

const tagConfig = "config"

var ConfigFilePath *string

// Config the program configuration, yml
var Config struct {
	XrayHelper struct {
		CoreType      string   `yaml:"coreType"`
		CorePath      string   `yaml:"corePath"`
		CoreConfig    string   `yaml:"coreConfig"`
		DataDir       string   `yaml:"dataDir"`
		RunDir        string   `yaml:"runDir"`
	} `yaml:"xrayHelper"`
	Proxy struct {
		Method          string   `yaml:"method"`
		TproxyPort      string   `yaml:"tproxyPort"`
		EnableIPv6      bool     `yaml:"enableIPv6"`
		Mode            string   `yaml:"mode"`
		PkgList         []string `yaml:"pkgList"`
		ApList          []string `yaml:"apList"`
		IgnoreList      []string `yaml:"ignoreList"`
		IntraList       []string `yaml:"intraList"`
	} `yaml:"proxy"`
}

// LoadConfig load program configuration file, should be called before any command Execute
func LoadConfig() error {
	configFile, err := os.ReadFile(*ConfigFilePath)
	if err != nil {
		return e.New("load config failed, ", err).WithPrefix(tagConfig)
	}
	if err := yaml.Unmarshal(configFile, &Config); err != nil {
		return e.New("unmarshal config failed, ", err).WithPrefix(tagConfig)
	}
	log.HandleDebug(Config.XrayHelper)
	log.HandleDebug(Config.Proxy)
	return nil
}
