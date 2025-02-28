package g

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

// IsExist checks whether a file or directory exists
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// ReadFileAndTrim reads the file content and returns it as a trimmed string
func readFileAndTrim(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

type HttpConfig struct {
	Enable bool   `json:"enable"`
	Listen string `json:"listen"`
	Token  string `json:"token"`
}

type MailConfig struct {
	Enable            bool   `json:"enable"`
	SendConcurrent    int    `json:"sendConcurrent"`
	MaxQueueSize      int    `json:"maxQueueSize"`
	FromUser          string `json:"fromUser"`
	MailServerHost    string `json:"mailServerHost"`
	MailServerPort    int    `json:"mailServerPort"`
	MailServerAccount string `json:"mailServerAccount"`
	MailServerPasswd  string `json:"mailServerPasswd"`
}

type GlobalConfig struct {
	Debug bool        `json:"debug"`
	Http  *HttpConfig `json:"http"`
	Mail  *MailConfig `json:"mail"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func GetConfig() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func LoadConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !isExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := readFileAndTrim(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	log.Println("g.ParseConfig ok, file ", cfg)
}
