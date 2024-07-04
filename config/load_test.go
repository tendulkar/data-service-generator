package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"stellarsky.ai/platform/codegen/data-service-generator/base"
)

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig()
	assert.NotNil(t, cfg)
	base.LOG.Info("Config loaded", "config", cfg)
}
