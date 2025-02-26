package config

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type EnvVariables struct {
	SbiUserName         string `envconfig:"SBI_USERNAME"`
	SbiPassword         string `envconfig:"SBI_PASSWORD"`
	SbiTorihikiPassword string `envconfig:"SBI_TORIHIKI_PASSWORD"`
	LogLevel            string `envconfig:"LOG_LEVEL" default:"INFO"` // INFO or DEBUG
	Headless            bool   `envconfig:"HEADLESS" default:"true"`
	GCPProjectID        string `envconfig:"GCP_PROJECT_ID"`
	Env                 string `envconfig:"ENV" default:"local"` // local or gcp
}

func (e *EnvVariables) isLocal() bool {
	return e.Env == "local"
}

func (e *EnvVariables) isGCP() bool {
	return e.Env == "gcp"
}

func (e *EnvVariables) loadSecrets(ctx context.Context, client *secretmanager.Client) error {
	if e.Env == "gcp" {
		// Fetch secrets from GCP Secret Manager
		secretName := fmt.Sprintf("projects/%s/secrets/%s-password/versions/latest", e.GCPProjectID, e.SbiUserName)
		var err error
		e.SbiPassword, err = getSecret(ctx, client, secretName)
		if err != nil {
			return errors.Wrapf(err, "failed to get password secret: %s", secretName)
		}

		secretName = fmt.Sprintf("projects/%s/secrets/%s-torihiki-password/versions/latest", e.GCPProjectID, e.SbiUserName)
		e.SbiTorihikiPassword, err = getSecret(ctx, client, secretName)
		if err != nil {
			return errors.Wrapf(err, "failed to get torihiki password secret: %s", secretName)
		}
	}
	return nil
}

func getSecret(ctx context.Context, client *secretmanager.Client, name string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return "", err
	}
	return string(result.Payload.Data), nil
}

func LoadEnvVariables() (*EnvVariables, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create secret manager client")
	}
	defer client.Close()

	var ev EnvVariables
	if err := envconfig.Process("", &ev); err != nil {
		return nil, errors.Wrap(err, "failed to process envconfig")
	}

	if ev.isGCP() {
		if err := ev.loadSecrets(ctx, client); err != nil {
			return nil, err
		}
	}

	if err := ev.validate(); err != nil {
		return nil, err
	}

	return &ev, nil
}

func (ev *EnvVariables) validate() error {
	checks := []struct {
		bad    bool
		errMsg string
	}{
		{
			ev.SbiUserName == "",
			"SBI_USERNAME is required",
		},
		{
			ev.isLocal() && ev.SbiPassword == "",
			"SBI_PASSWORD is required",
		},
		{
			ev.isLocal() && ev.SbiTorihikiPassword == "",
			"SBI_TORIHIKI_PASSWORD is required",
		},
		{
			!(ev.LogLevel == "INFO" || ev.LogLevel == "DEBUG" || ev.LogLevel == ""),
			fmt.Sprintf("invalid LOG_LEVEL is specified: %s", ev.LogLevel),
		},
	}

	for _, check := range checks {
		if check.bad {
			return errors.Errorf(check.errMsg)
		}
	}

	return nil
}
