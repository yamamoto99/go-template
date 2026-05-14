package config

import (
	"fmt"
	"os"
	"strings"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

type Config struct {
	DB     DBConfig
	TestDB DBConfig
}

func Load() (*Config, error) {
	var missing []string
	get := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	cfg := &Config{
		DB: DBConfig{
			Host:     get("DB_HOST"),
			Port:     get("DB_PORT"),
			User:     get("DB_USER"),
			Password: get("DB_PASSWORD"),
			Name:     get("DB_NAME"),
		},
	}

	if err := validate(missing); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadTest() (*Config, error) {
	var missing []string
	get := func(key string) string {
		v := os.Getenv(key)
		if v == "" {
			missing = append(missing, key)
		}
		return v
	}

	cfg := &Config{
		TestDB: DBConfig{
			Host:     get("TEST_DB_HOST"),
			Port:     get("TEST_DB_PORT"),
			User:     get("TEST_DB_USER"),
			Password: get("TEST_DB_PASSWORD"),
			Name:     get("TEST_DB_NAME"),
		},
	}

	if err := validate(missing); err != nil {
		return nil, err
	}
	return cfg, nil
}

func validate(missing []string) error {
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("missing required env: %s", strings.Join(missing, ", "))
}
