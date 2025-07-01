package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/config/env"
)

var appConfig *config

type config struct {
	Logger  LoggerConfig
	UFOGRPC UFOGRPCConfig
	Mongo   MongoConfig
	Redis   RedisConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	ufoGRPCCfg, err := env.NewUFOGRPCConfig()
	if err != nil {
		return err
	}

	mongoCfg, err := env.NewMongoConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:  loggerCfg,
		UFOGRPC: ufoGRPCCfg,
		Mongo:   mongoCfg,
		Redis:   redisCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
