package core

import (
	"fmt"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type CapabilityConfiguration struct {
	ModuleSource        string            `yaml:"module" json:"module"`
	ModuleSourceVersion *string           `yaml:"module_version" json:"moduleVersion"`
	Variables           map[string]any    `yaml:"vars" json:"vars"`
	Connections         ConnectionTargets `yaml:"connections" json:"connections"`
	Namespace           *string           `yaml:"namespace" json:"namespace"`
}

type InvalidConfigurationError struct {
	Err error
}

func (e InvalidConfigurationError) Error() string {
	return fmt.Sprintf("invalid app configuration: %s", e.Err.Error())
}

func (c CapabilityConfiguration) Normalize(resolver *find.ResourceResolver) (CapabilityConfiguration, error) {
	if err := NormalizeConnectionTargets(c.Connections, resolver); err != nil {
		return c, err
	}
	return c, nil
}
