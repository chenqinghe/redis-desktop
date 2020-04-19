package config

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type Config struct {
	Lang      string    `toml:"lang"`
	LogConfig LogConfig `toml:"log"`
}

type LogConfig struct {
	Level string `toml:"level"`
}

var config = Config{
	Lang: "zh_cn",
	LogConfig: LogConfig{
		Level: "info",
	},
}

func Get() Config {
	return config
}

func Load(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return toml.Unmarshal(data, &config)
}

func Save(file string) error {
	buf := bytes.NewBuffer(nil)
	encoder := toml.NewEncoder(buf)
	if err := encoder.Encode(config); err != nil {
		return err
	}

	return ioutil.WriteFile(file, buf.Bytes(), 0644)
}
