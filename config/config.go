package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/go-chujang/demo-aa/common/logx"
	"github.com/go-chujang/demo-aa/common/utils/ternary"
	"gopkg.in/yaml.v3"
)

const (
	commonTag = "common"
	local     = "local"
	demo      = "demo"
)

var (
	apptag, envtag string
	confHolder     = map[key]string{}
	ErrEmptyString = errors.New("empty string")
)

func IsRelease() bool { return envtag == demo }
func AppTag() string  { return apptag }
func EnvTag() string  { return envtag }

type conf struct {
	Default map[string]string `yaml:"default"`
	Local   map[string]string `yaml:"local"`
	Demo    map[string]string `yaml:"demo"`
}

func init() {
	if apptag = os.Getenv(APP_TAG.String()); apptag == "" {
		dir, _ := os.Getwd()
		apptag = filepath.Base(dir)
		// panic("APP_TAG must not be empty")
	}
	if envtag = os.Getenv(ENV_TAG.String()); envtag == "" {
		envtag = local
	}

	path := ternary.Cond(envtag == local, "../../.config/config.yaml", "config.yaml")
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	parsed := make(map[string]conf)
	if err = yaml.Unmarshal(bytes, parsed); err != nil {
		panic(err)
	}

	var (
		common map[string]string
		appenv map[string]string
	)
	switch envtag {
	case local:
		common = parsed[commonTag].Local
		appenv = parsed[apptag].Local
	case demo:
		common = parsed[commonTag].Demo
		appenv = parsed[apptag].Demo
	default:
		panic("unsupported environment tag")
	}

	// set common
	for k, v := range parsed[commonTag].Default {
		confHolder[key(k)] = v
	}
	for k, v := range common {
		confHolder[key(k)] = v
	}
	// set appenv
	for k, v := range parsed[apptag].Default {
		confHolder[key(k)] = v
	}
	for k, v := range appenv {
		confHolder[key(k)] = v
	}
	// print env
	logx.Write(apptag, "APP_TAG\t%s", apptag)
	logx.Write(apptag, "ENV_TAG\t%s", envtag)
	for k, v := range confHolder {
		if envtag == demo {
			logx.Write(apptag, "%s\t%s", k, v)
		}
		os.Setenv(k.String(), v)
	}
}
