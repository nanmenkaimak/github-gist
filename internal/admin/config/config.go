package config

import "time"

type Config struct {
	HttpServer HttpServer `yaml:"HttpServer"`
	Databases  Databases  `yaml:"Databases"`
	Auth       Auth       `yaml:"JwtSecretKey"`
	Storage    Storage    `yaml:"Storage"`
}

type HttpServer struct {
	Port            int           `yaml:"Port"`
	ShutdownTimeout time.Duration `yaml:"ShutdownTimeout"`
}

type Databases struct {
	Gist Database `yaml:"Gist"`
	User Database `yaml:"User"`
}

type Database struct {
	Main    DbNode `yaml:"Main"`
	Replica DbNode `yaml:"Replica"`
}

type DbNode struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Name     string `yaml:"Name"`
	SslMode  string `yaml:"SslMode"`
}

type Auth struct {
	JwtSecretKey string `yaml:"JwtSecretKey"`
}

type Storage struct {
	Interval time.Duration `yaml:"Interval"`
}
