package overrides

import (
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/yaml.v3"
)

type ConfigurationOverrides struct {
	Version      string                          `yaml:"version"`
	Applications map[string]ApplicationOverrides `yaml:"apps,omitempty"`
	Subdomains   map[string]SubdomainOverrides   `yaml:"subdomains,omitempty"`
	Datastores   map[string]DatastoreOverrides   `yaml:"datastores,omitempty"`
}

func (c *ConfigurationOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	ve := errors.ValidationErrors{}
	for _, subdomain := range c.Subdomains {
		verrs, err := subdomain.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, datastore := range c.Datastores {
		verrs, err := datastore.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, app := range c.Applications {
		verrs, err := app.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	return ve, nil
}

func (c *ConfigurationOverrides) Normalize(resolver *find.ResourceResolver) error {
	for key, subdomain := range c.Subdomains {
		if err := subdomain.Normalize(resolver); err != nil {
			return err
		}
		c.Subdomains[key] = subdomain
	}
	for key, datastore := range c.Datastores {
		if err := datastore.Normalize(resolver); err != nil {
			return err
		}
		c.Datastores[key] = datastore
	}
	for key, app := range c.Applications {
		if err := app.Normalize(resolver); err != nil {
			return err
		}
		c.Applications[key] = app
	}
	return nil
}

func ParseConfigurationOverrides(data []byte) (*ConfigurationOverrides, error) {
	var r *ConfigurationOverrides
	err := yaml.Unmarshal(data, &r)
	for k, ao := range r.Applications {
		ao.Name = k
		r.Applications[k] = ao
	}
	for k, so := range r.Subdomains {
		so.Name = k
		r.Subdomains[k] = so
	}
	return r, err
}

type ApplicationOverrides struct {
	Name         string              `yaml:"-"`
	Variables    map[string]any      `yaml:"vars"`
	EnvVariables map[string]string   `yaml:"environment"`
	Capabilities CapabilityOverrides `yaml:"capabilities"`
}

func (a *ApplicationOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	// TODO: Implement: How do we validate if we don't have a module to resolve
	return errors.ValidationErrors{}, nil
}

func (a *ApplicationOverrides) Normalize(resolver *find.ResourceResolver) error {
	return a.Capabilities.Normalize(resolver)
}

type SubdomainOverrides struct {
	Name        string                 `yaml:"-"`
	Variables   map[string]any         `yaml:"vars"`
	Connections core.ConnectionTargets `yaml:"connections"`
}

func (s *SubdomainOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	// TODO: Implement: How do we validate if we don't have a module to resolve
	return errors.ValidationErrors{}, nil
}

func (s *SubdomainOverrides) Normalize(resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(s.Connections, resolver)
}

type DatastoreOverrides struct {
	Name        string                 `yaml:"-"`
	Variables   map[string]any         `yaml:"vars"`
	Connections core.ConnectionTargets `yaml:"connections"`
}

func (d *DatastoreOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	// TODO: Implement: How do we validate if we don't have a module to resolve
	return errors.ValidationErrors{}, nil
}

func (d *DatastoreOverrides) Normalize(resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(d.Connections, resolver)
}
