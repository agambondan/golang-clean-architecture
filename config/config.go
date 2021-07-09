package config

import "os"

type Configuration struct {
	DBDriver      string
	DBHost        string
	DBName        string
	DBUser        string
	DBPassword    string
	DBPort        string
	Port          string
	RedisHost     string
	RedisPort     string
	RedisPassword string
}

func (config *Configuration) Init() *Configuration {
	config.DBDriver = os.Getenv("DB_DRIVER")
	config.DBHost = os.Getenv("DB_HOST")
	config.DBPassword = os.Getenv("DB_PASSWORD")
	config.DBUser = os.Getenv("DB_USER")
	config.DBName = os.Getenv("DB_NAME")
	config.DBPort = os.Getenv("DB_PORT")
	config.Port = os.Getenv("PORT")
	config.RedisHost = os.Getenv("REDIS_HOST")
	config.RedisPort = os.Getenv("REDIS_PORT")
	config.RedisPassword = os.Getenv("REDIS_PASSWORD")
	return config
}
