package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type ufoGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type ufoGRPCConfig struct {
	raw ufoGRPCEnvConfig
}

func NewUFOGRPCConfig() (*ufoGRPCConfig, error) {
	var raw ufoGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &ufoGRPCConfig{raw: raw}, nil
}

func (cfg *ufoGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
