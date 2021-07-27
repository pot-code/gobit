package util

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// LoadConfig init app config using viper
func LoadConfig(envPrefix string, def interface{}) error {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// env
	viper.AutomaticEnv()
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// config files
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/go/") // path to look for the config file in
	viper.AddConfigPath("$HOME/")   // call multiple times to add many search paths
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("[skipped]config file dose not present.")
		} else {
			return err
		}
	}

	if err := viper.Unmarshal(def); err != nil {
		return err
	}
	if err := validateConfig(def); err != nil {
		return err
	}

	return nil
}

func validateConfig(config interface{}) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" || name == "" {
			name = fld.Tag.Get("yaml")
			if name == "-" || name == "" {
				return ""
			}
		}
		return name
	})

	err := validate.Struct(config)
	if err == nil {
		return nil
	} else if errors.Is(err, &validator.InvalidValidationError{}) {
		return fmt.Errorf("failed to validate config: %s", err)
	}

	var msg []string
	for _, field := range err.(validator.ValidationErrors) {
		fieldName := field.Namespace()
		switch field.Tag() {
		case "required":
			msg = append(msg, fmt.Sprintf("%s is required", fieldName))
		case "oneof":
			msg = append(msg, fmt.Sprintf("%s must be one of (%s)", fieldName, field.Param()))
		default:
			msg = append(msg, field.Error())
		}
	}
	if len(msg) > 0 {
		return fmt.Errorf("failed to validate config: \n%s", strings.Join(msg, "\n"))
	}
	return err
}
