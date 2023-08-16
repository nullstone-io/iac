package config

import (
	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/yaml.v3"
)

type EnvConfiguration struct {
	Version      string                            `yaml:"version" json:"version"`
	Subdomains   map[string]SubdomainConfiguration `yaml:"subdomains" json:"subdomains"`
	Applications map[string]AppConfiguration       `yaml:"apps" json:"apps"`
}

func ParseEnvConfiguration(data []byte) (*EnvConfiguration, error) {
	var r *EnvConfiguration
	err := yaml.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	newApps := make(map[string]AppConfiguration)
	for appName, appValue := range r.Applications {
		appValue.Name = appName
		// set a default module version if not provided
		if appValue.ModuleSourceVersion == nil {
			latest := "latest"
			appValue.ModuleSourceVersion = &latest
		}
		newCaps := make([]CapabilityConfiguration, len(appValue.Capabilities))
		for i, capValue := range appValue.Capabilities {
			// set a default module version if not provided
			if capValue.ModuleSourceVersion == nil {
				latest := "latest"
				capValue.ModuleSourceVersion = &latest
			}
			newCaps[i] = capValue
		}
		appValue.Capabilities = newCaps
		newApps[appName] = appValue
	}
	r.Applications = newApps

	newSubdomains := make(map[string]SubdomainConfiguration)
	for subdomainName, subdomainValue := range r.Subdomains {
		subdomainValue.Name = subdomainName
		// set a default module version if not provided
		if subdomainValue.ModuleSourceVersion == nil {
			latest := "latest"
			subdomainValue.ModuleSourceVersion = &latest
		}
		newSubdomains[subdomainName] = subdomainValue
	}
	r.Subdomains = newSubdomains

	return r, err
}

func (e EnvConfiguration) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	ve := errors.ValidationErrors{}
	for _, subdomain := range e.Subdomains {
		verrs, err := subdomain.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, app := range e.Applications {
		verrs, err := app.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	return ve, nil
}

func (e *EnvConfiguration) Normalize(resolver *find.ResourceResolver) error {
	for key, subdomain := range e.Subdomains {
		if err := subdomain.Normalize(resolver); err != nil {
			return err
		}
		e.Subdomains[key] = subdomain
	}
	for key, app := range e.Applications {
		if err := app.Normalize(resolver); err != nil {
			return err
		}
		e.Applications[key] = app
	}
	return nil
}
