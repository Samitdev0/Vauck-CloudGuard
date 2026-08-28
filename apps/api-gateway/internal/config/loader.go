package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	env := strings.TrimSpace(os.Getenv("APP_ENV"))

	if env == "" {
		env = "development"
	}

	v := viper.New()

	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// ==========================
	// Environment Overrides
	// ==========================

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// JWT
	v.BindEnv(
		"jwt.secret",
		"CLOUDGUARD_JWT_SECRET",
	)

	v.BindEnv(
		"jwt.issuer",
		"CLOUDGUARD_JWT_ISSUER",
	)

	// ==========================
	// Unmarshal
	// ==========================

	var cfg Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// ==========================
	// Validation
	// ==========================

	if strings.TrimSpace(cfg.JWT.Secret) == "" {
		return nil, fmt.Errorf("jwt secret is not configured")
	}

	if len(cfg.JWT.Secret) < 32 {
		return nil, fmt.Errorf(
			"jwt secret must contain at least 32 characters",
		)
	}

	if strings.TrimSpace(cfg.JWT.Issuer) == "" {
		cfg.JWT.Issuer = "CloudGuard"
	}

	return &cfg, nil
}
