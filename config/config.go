package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type EnvVariables struct {
	SbiUserName         string `envconfig:"SBI_USERNAME"`
	SbiPassword         string `envconfig:"SBI_PASSWORD"`
	SbiTorihikiPassword string `envconfig:"SBI_TORIHIKI_PASSWORD"`
	LogLevel            string `envconfig:"LOG_LEVEL" default:"INFO"` // INFO or DEBUG
	Headless            bool   `envconfig:"HEADLESS" default:"true"`
}

func LoadEnvVariables() (*EnvVariables, error) {
	var c EnvVariables
	if err := envconfig.Process("", &c); err != nil {
		return nil, errors.Wrap(err, "failed to process envconfig")
	}

	if err := c.validate(); err != nil {
		return nil, err
	}

	return &c, nil
}

func (e *EnvVariables) validate() error {
	checks := []struct {
		bad    bool
		errMsg string
	}{
		{
			e.SbiUserName == "",
			"SBI_USERNAME is required",
		},
		{
			e.SbiPassword == "",
			"SBI_PASSWORD is required",
		},
		{
			e.SbiTorihikiPassword == "",
			"SBI_TORIHIKI_PASSWORD is required",
		},
		{
			!(e.LogLevel == "INFO" || e.LogLevel == "DEBUG" || e.LogLevel == ""),
			fmt.Sprintf("invalid LOG_LEVEL is specified: %s", e.LogLevel),
		},
	}

	for _, check := range checks {
		if check.bad {
			return errors.Errorf(check.errMsg)
		}
	}

	return nil
}
