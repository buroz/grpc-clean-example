package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type ArangoConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	Name     string `json:"name"`
}

type AmqpConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int64  `json:"port"`
}

type SmtpConfig struct {
	SenderName  string `json:"sender_name"`
	SenderEmail string `json:"sender_email"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        int64  `json:"port"`
	Secure      bool   `json:"secure"`
}

type Config struct {
	Db   ArangoConfig `json:"arango"`
	Amqp AmqpConfig   `json:"amqp"`
	Smtp SmtpConfig   `json:"smtp"`
}

func NewConfig() (*Config, error) {
	conf := &Config{}

	var configPath = ""

	flag.StringVar(&configPath, "config", "", "Config file path")

	flag.Parse()

	if configPath == "" {
		return nil, fmt.Errorf("Config file path is empty")
	}

	conf = &Config{}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(byteValue, &conf)

	return conf, nil
}
