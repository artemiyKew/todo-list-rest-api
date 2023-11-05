package api

type Config struct {
	BindAddr    string `yaml:"bind_addr"`
	DataBaseURL string `yaml:"db_url"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
	}
}
