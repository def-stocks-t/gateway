package config

type Config struct {
	Port string `env:"GATEWAY_PORT" json:"gateway_port" default:"8080"`

	Services struct {
		Storage struct {
			Host string `env:"STORAGE_HOST" json:"storage_host" default:"127.0.0.1"`
			Port string `env:"STORAGE_PORT" json:"storage_port" default:"8081"`
		}
	}
}
