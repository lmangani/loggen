package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL       string            `yaml:"url"`
	APIKey    string            `yaml:"api_key"`
	APISecret string            `yaml:"api_secret"`
	Labels    map[string]string `yaml:"labels"`
	Rate      int               `yaml:"rate"`
	Timeout   time.Duration     `yaml:"timeout"`
}

var (
	basePath       = fmt.Sprintf("%s/.loggen", os.Getenv("HOME"))
	configFilename = fmt.Sprintf("%s/config.yaml", basePath)
	c              = &Config{}
)

func Load() {
	f, err := os.Open(configFilename)
	if err != nil {
		fmt.Println("Creating default config...")
		if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
			log.Printf("unable create config file directory: %v", err)
			return
		}

		c = getDefaultConfig()
		b, _ := yaml.Marshal(c)
		if err := os.WriteFile(configFilename, b, os.ModePerm); err != nil {
			log.Printf("unable to create config file: %v", err)
		}
		return
	}

	b, _ := io.ReadAll(f)
	if err := yaml.Unmarshal(b, c); err != nil {
		fmt.Println(err)
	}
	_ = f.Close()
}

func Get() *Config {
	return c
}

func getDefaultConfig() *Config {
	return &Config{
		URL:     "https://qryn.gigapipe.com",
		Rate:    100,
		Timeout: 30 * time.Second,
	}
}
