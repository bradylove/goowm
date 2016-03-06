package config

import (
	"github.com/spf13/viper"
	"gopkg.in/fsnotify.v1"
)

type Config struct {
	Display string

	Workspaces []*WorkspaceConfig

	KeyBindingNextWorkspace     string
	KeyBindingPreviousWorkspace string
	KeyBindingExecLauncher      string

	ThemeBorderColor int
	ThemeBorderWidth int
}

type WorkspaceConfig struct {
	Name string
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

	// KeyBindings
	c.KeyBindingExecLauncher = viper.GetString("key-bindings.exec_launcher")
	c.KeyBindingNextWorkspace = viper.GetString("key-bindings.next_workspace")
	c.KeyBindingPreviousWorkspace = viper.GetString("key-bindings.previous_workspace")

	// Workspaces
	for _, w := range viper.Get("workspace").([]map[string]interface{}) {
		c.Workspaces = append(c.Workspaces, &WorkspaceConfig{Name: w["name"].(string)})
	}

	// Theme related settings
	c.ThemeBorderColor = ParseHexColor(viper.GetString("theme.border_color"))
	c.ThemeBorderWidth = viper.GetInt("theme.border_width")
}

func (c *Config) onConfigChange(e fsnotify.Event) {
	c.reload()
}
