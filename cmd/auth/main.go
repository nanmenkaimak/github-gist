// Auth-service API.
//
//		Schemes: https, http
//	    Host: localhost:8080
//		BasePath: /api/auth
//		Version: 0.0.1
//
//		Consumes:
//		- application/json
//
//		Produces:
//		- application/json
//
//		SecurityDefinitions:
//		  Bearer:
//		    type: apiKey
//		    name: Authorization
//		    in: header
//
// swagger:meta
package main

import (
	"fmt"
	"github.com/nanmenkaimak/github-gist/internal/auth/applicator"
	"github.com/nanmenkaimak/github-gist/internal/auth/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	//nolint:all
	defer logger.Sync()

	l := logger.Sugar()
	l = l.With(zap.String("applicator", "auth-service"))

	cfg, err := loadConfig("config/auth")
	if err != nil {
		l.Fatalf("failed to load config err: %v", err)
	}

	app := applicator.NewApp(l, &cfg)
	app.Run()
}

func loadConfig(path string) (config config.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to ReadInConfig err: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to Unmarshal config err: %w", err)
	}
	return config, nil
}
