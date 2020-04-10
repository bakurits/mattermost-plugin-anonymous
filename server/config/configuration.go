package config

import (
	"context"
	"reflect"
)

// Config captures the plugin's external Config as exposed in the Mattermost server
// Config, as well as values computed from the Config. Any public fields will be
// deserialized from the Mattermost server Config in OnConfigurationChange.
//
// As plugins are inherently concurrent (hooks being called asynchronously), and the plugin
// Config can change at any time, access to the Config must be synchronized. The
// strategy used in this plugin is to guard a pointer to the Config, and clone the entire
// struct whenever it changes. You may replace this with whatever strategy you choose.
//
// If you add non-reference types to your Config struct, be sure to rewrite Clone as a deep
// copy appropriate for your types.
type Config struct {
	PluginID      string
	PluginVersion string
}

// Clone shallow copies the Config. Your implementation may require a deep copy if
// your Config has reference types.
func (c *Config) Clone() *Config {
	var clone = *c
	return &clone
}

var contextKey = reflect.TypeOf(Config{})

// Context sets config object in context
func Context(ctx context.Context, conf *Config) context.Context {
	return context.WithValue(ctx, contextKey, conf)
}

// FromContext loads context object from context
func FromContext(ctx context.Context) *Config {
	return ctx.Value(contextKey).(*Config)
}
