package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config interface {
	GetString(string) string
	GetInt(string) int
	GetBool(string) bool
	GetDuration(string) time.Duration
	GetFloat64(string) float64
}

type config struct {
	Env *viper.Viper
}

func NewConfig() Config {
	c := &config{Env: viper.New()}
	c.Env.AddConfigPath(".")
	c.Env.AddConfigPath("./param")
	c.Env.SetConfigName("tempokerja")
	c.Env.SetConfigType("toml")

	// Check read process
	if err := c.Env.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Config error: %s", err))
	}

	fmt.Printf("=> config file: %s\n", c.Env.ConfigFileUsed())

	return c
}

func (c *config) GetString(k string) string {
	return c.Env.GetString(k)
}

func (c *config) GetInt(k string) int {
	return c.Env.GetInt(k)
}

func (c *config) GetBool(k string) bool {
	return c.Env.GetBool(k)
}

func (c *config) GetDuration(k string) time.Duration {
	return c.Env.GetDuration(k)
}

func (c *config) GetFloat64(k string) float64 {
	return c.Env.GetFloat64(k)
}
