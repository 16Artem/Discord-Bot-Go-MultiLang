package config

import (
    "io/ioutil"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Token       string `yaml:"token"`
    DefaultLang string `yaml:"default_lang"`
    OwnerID     string `yaml:"owner_id"`
}

func LoadConfig(path string) (*Config, error) {
    b, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var c Config
    if err := yaml.Unmarshal(b, &c); err != nil {
        return nil, err
    }
    if c.DefaultLang == "" {
        c.DefaultLang = "en"
    }
    return &c, nil
}
