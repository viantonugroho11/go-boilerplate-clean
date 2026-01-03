package config

type App struct {
	Port string `json:"port"`
	Environment string `json:"environment"`
	LogLevel string `json:"log_level"`
}

func (a App) IsProduction() bool {
	return a.Environment == "production"
}

func (a App) IsDevelopment() bool {
	return a.Environment == "development"
}

func (a App) IsStaging() bool {
	return a.Environment == "staging"
}