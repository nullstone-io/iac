package config

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type SubdomainConfiguration struct {
	Name string `yaml:"-" json:"name"`
	// TODO: Implement DnsName in DesiredConfig
	DnsName             string                 `yaml:"dns_name" json:"dnsName"`
	ModuleSource        string                 `yaml:"module" json:"module"`
	ModuleSourceVersion *string                `yaml:"module_version" json:"moduleVersion"`
	Variables           map[string]any         `yaml:"vars" json:"vars"`
	Connections         core.ConnectionTargets `yaml:"connections" json:"connections"`
}

func (s SubdomainConfiguration) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	yamlPath := fmt.Sprintf("subdomains.%s", s.Name)
	return ValidateBlock(resolver, yamlPath, "subdomain/*/*", s.ModuleSource, *s.ModuleSourceVersion, s.Variables, s.Connections, nil)
}

func (s *SubdomainConfiguration) Normalize(resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(s.Connections, resolver)
}
