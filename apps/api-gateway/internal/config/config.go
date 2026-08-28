package config

import "fmt"

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	HTTP     HTTPConfig     `mapstructure:"http"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Service     string `mapstructure:"service"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

type HTTPConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"readTimeout"`
	WriteTimeout int    `mapstructure:"writeTimeout"`
	IdleTimeout  int    `mapstructure:"idleTimeout"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxOpen  int32  `mapstructure:"maxOpen"`
	MinOpen  int32  `mapstructure:"minOpen"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Issuer string `mapstructure:"issuer"`
}

func (d DatabaseConfig) DSN() string {
	return "postgres://" +
		d.User + ":" +
		d.Password + "@" +
		d.Host + ":" +
		fmt.Sprint(d.Port) + "/" +
		d.Name +
		"?sslmode=" +
		d.SSLMode
}
