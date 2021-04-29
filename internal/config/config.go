package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"os"
)

const filePath = "internal/config/config.yaml"

type Configuration struct {
	Host     string                `yaml:"host"`
	Service  Service               `yaml:"service"`
	Database DatabaseConfiguration `yaml:"database"`
}
type Service struct {
	Host string `yaml:"host"`
}
type DatabaseConfiguration struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"dbname"`
}

func NewConfig() (*Configuration, error) {
	conf := &Configuration{}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(wd + "/" + filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}

func (c *Configuration) GetdbUrl() *url.URL {
	url := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.Database.User, c.Database.Password),
		Host:   c.Database.Host,
		Path:   c.Database.Database,
	}
	q := url.Query()
	q.Add("sslmode", "disable")
	url.RawQuery = q.Encode()
	return url
}
