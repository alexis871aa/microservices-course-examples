package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/ufo/internal/config/env"
)

var appConfig *config

type config struct {
	Logger  LoggerConfig
	UFOGRPC UFOGRPCConfig
	Mongo   MongoConfig
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

	appConfig = &config{
		Logger:  loggerCfg,
		UFOGRPC: ufoGRPCCfg,
		Mongo:   mongoCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
