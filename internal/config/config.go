package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Secret   SecretConfig
	Nats     NatsConfig
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
}

type ServerConfig struct {
	PortServer string
}

type SecretConfig struct {
	JwtSecretKey string
}

type NatsConfig struct {
	Dns string
}

type RabbitMQConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type DatabaseConfig struct {
	Master DatabaseMaster
	Slave  DatabaseSlave
}

type DatabaseMaster struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type DatabaseSlave struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type Consul struct {
	HostConsul    string
	KeyConsul     string
	ServiceConsul string
}

func ReadFileEnv(env string) (*Consul, error) {
	err := godotenv.Load(env)
	if err != nil {
		return nil, err
	}

	data := &Consul{
		HostConsul:    os.Getenv("HostConsul"),
		KeyConsul:     os.Getenv("KeyConsul"),
		ServiceConsul: os.Getenv("ServiceConsul"),
	}
	return data, nil
}
