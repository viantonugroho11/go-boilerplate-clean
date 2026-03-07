package bootstrap

import (
	"os"

	"go-boilerplate-clean/internal/config"

	confLoader "github.com/viantonugroho11/go-config-library"
)

var cfg *config.Configuration

// LoadConfig memuat konfigurasi sekali dan simpan global (per proses). Wajib dipanggil sebelum Init* / Run*.
func LoadConfig() error {
	c, err := loadConfig()
	if err != nil {
		return err
	}
	cfg = c
	return nil
}

// Config mengembalikan config yang sudah di-load. Panic bila LoadConfig belum dipanggil.
func Config() *config.Configuration {
	if cfg == nil {
		panic("bootstrap: LoadConfig() harus dipanggil dulu")
	}
	return cfg
}

func loadConfig() (*config.Configuration, error) {
	c := &config.Configuration{}
	loader := confLoader.New("", "go-boilerplate-clean", os.Getenv("CONSUL_URL"),
		confLoader.WithConfigFileSearchPaths("./config"),
	)
	if err := loader.Load(c); err != nil {
		return nil, err
	}
	return c, nil
}
