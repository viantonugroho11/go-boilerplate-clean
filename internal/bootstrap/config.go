package bootstrap

import (
	"os"

	"go-boilerplate-clean/internal/config"

	confLoader "github.com/viantonugroho11/go-config-library"
)

// LoadConfig loads and returns the application configuration.
func LoadConfig() (*config.Configuration, error) {
	c := &config.Configuration{}
	loader := confLoader.New("", "go-boilerplate-clean", os.Getenv("CONSUL_URL"),
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	if err := loader.Load(c); err != nil {
		return nil, err
	}
	return c, nil
}
