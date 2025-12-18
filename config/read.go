package config

import (
	"fmt"
	"reflect"

	"github.com/Hossara/quera_bootcamp_chatapp_backend/pkg/constants"
	"github.com/spf13/viper"
)

func ReadConfig(configPath string) (*Config, error) {
	viper.SetConfigName(constants.ConfigName)
	viper.SetConfigType(constants.ConfigFormat)
	viper.AddConfigPath(configPath)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal into the Config struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %v", err)
	}

	// Merge with defaults
	mergeWithDefaults(&config)

	return &config, nil
}

// mergeWithDefaults automatically applies default values from DefaultConfig to the given config.
func mergeWithDefaults(config *Config) {
	applyDefaults(reflect.ValueOf(config), reflect.ValueOf(&DefaultConfig).Elem())
}

func applyDefaults(target, defaults reflect.Value) {
	for i := 0; i < target.Elem().NumField(); i++ {
		targetField := target.Elem().Field(i)
		defaultField := defaults.Field(i)

		if targetField.Kind() == reflect.Struct {
			applyDefaults(targetField.Addr(), defaultField)
		} else if isZeroValue(targetField) {
			targetField.Set(defaultField)
		}
	}
}

func isZeroValue(value reflect.Value) bool {
	return value.IsZero()
}

func MustReadConfig(path string) *Config {
	config, err := ReadConfig(path)
	if err != nil {
		panic(err)
	}
	return config
}
