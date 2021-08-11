package ldaputil

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server string
	Bind   struct {
		DN     string
		Secret string
	}
	BaseDN string `yaml:"base_dn"`
	Listen string
}

func ParseConfig(file string) (cfg Config, err error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(b, &cfg)
	return
}
