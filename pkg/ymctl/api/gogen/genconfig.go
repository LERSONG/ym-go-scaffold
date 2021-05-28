package gogen

import (
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
)

const configTemplate = `package config

type Config struct {
	Addr string
	RegistryAddr string
}

func NewConfig(opts ...Option) *Config {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}
	return config
}


type Option = func(c *Config)

func Addr(addr string) Option {
	return func(c *Config) {
		c.Addr = addr
	}
}

func RegistryAddr(registryAddr string) Option {
	return func(c *Config) {
		c.RegistryAddr = registryAddr
	}
}
`

func genConfig(dir string, api *spec.ApiSpec) error {
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          configDir,
		filename:        "config.go",
		templateName:    "configTemplate",
		category:        category,
		templateFile:    configTemplateFile,
		builtinTemplate: configTemplate,
		data:            map[string]string{},
	})
}
