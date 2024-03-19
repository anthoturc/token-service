package main

import (
	"fmt"
	"os"
	"path"

	"dario.cat/mergo"
	"gopkg.in/yaml.v3"
)

type TimeoutConf struct {
	Read  uint `yaml:"read"`
	Write uint `yaml:"write"`
}

type ServerConf struct {
	Port        uint16 `yaml:"port"`
	Address     string `yaml:"address"`
	TimeoutConf `yaml:"timeout"`
}

type DbConf struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SslMode  string `yaml:"sslMode"`
}

// AppConf will hold basic application configuration (including secrets)
type AppConf struct {
	ServerConf `yaml:"server"`
	DbConf     `yaml:"db"`
}

// NewAppConf will return an instance of AppConf that is populated by the sources
// within the pathToConfDir directory. Any errors encountered are returned.
func NewAppConf(pathToConfDir string, sources ...string) (AppConf, error) {

	var dest AppConf
	for _, source := range sources {
		conf, err := unmarshalSource(path.Join(pathToConfDir, source))
		if err != nil {
			return AppConf{}, fmt.Errorf("new app conf - failed to unmarshal source: %v", err)
		}
		mergo.Merge(&dest, conf, mergo.WithOverride)
	}

	return dest, nil
}

// unmarshalSource will translate the pathToSource yaml file and any references to environment variables into an AppConf. This function should never
// be used directly, because the resultant AppConf may not have all the fields (and nested fields) populated.
func unmarshalSource(pathToSource string) (AppConf, error) {
	data, err := os.ReadFile(pathToSource)
	if err != nil {
		return AppConf{}, fmt.Errorf("unmarshal source - failed to read %s: %v", pathToSource, err)
	}

	// Apply the environment variables to this config if it applies.
	data = []byte(os.ExpandEnv(string(data)))

	var conf AppConf
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return AppConf{}, fmt.Errorf("unmarshal source - failed to unmarshal data in %s: %v", pathToSource, err)
	}

	return conf, nil
}
