package config

import (
	"github.com/spf13/viper"
	"gopkg.in/fsnotify.v1"
)

type Config struct {
	Display          string
	ThemeBorderWidth int
}

func Load(name string, paths ...string) (*Config, error) {
	c := new(Config)

	viper.SetConfigName("config")
	viper.OnConfigChange(c.onConfigChange)

	for _, p := range paths {
		viper.AddConfigPath(p)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	c.reload()
	return c, err
}

func (c *Config) reload() {
	// Base Settings
	c.Display = viper.GetString("display")

	// Theme related settings
	c.ThemeBorderWidth = viper.GetInt("theme.border_width")
}

func (c *Config) onConfigChange(e fsnotify.Event) {
	c.reload()
}
