package overrides

import (
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/yaml.v3"
)

type ConfigurationOverrides struct {
	Version           string                               `yaml:"version"`
	Applications      map[string]ApplicationOverrides      `yaml:"apps,omitempty"`
	Subdomains        map[string]SubdomainOverrides        `yaml:"subdomains,omitempty"`
	Datastores        map[string]DatastoreOverrides        `yaml:"datastores,omitempty"`
	Domains           map[string]DomainOverrides           `yaml:"domains,omitempty"`
	Ingresses         map[string]IngressOverrides          `yaml:"ingresses,omitempty"`
	ClusterNamespaces map[string]ClusterNamespaceOverrides `yaml:"cluster_namespaces,omitempty"`
	Clusters          map[string]ClusterOverrides          `yaml:"clusters,omitempty"`
	Networks          map[string]NetworkOverrides          `yaml:"networks,omitempty"`
	Blocks            map[string]BlockOverrides            `yaml:"blocks,omitempty"`
}

func (c *ConfigurationOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	ve := errors.ValidationErrors{}
	for _, block := range c.Blocks {
		verrs, err := block.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, network := range c.Networks {
		verrs, err := network.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, cluster := range c.Clusters {
		verrs, err := cluster.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, clusterNamespace := range c.ClusterNamespaces {
		verrs, err := clusterNamespace.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, ingress := range c.Ingresses {
		verrs, err := ingress.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, domain := range c.Domains {
		verrs, err := domain.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
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
	for key, block := range c.Blocks {
		if err := block.Normalize(resolver); err != nil {
			return err
		}
		c.Blocks[key] = block
	}
	for key, network := range c.Networks {
		if err := network.Normalize(resolver); err != nil {
			return err
		}
		c.Networks[key] = network
	}
	for key, cluster := range c.Clusters {
		if err := cluster.Normalize(resolver); err != nil {
			return err
		}
		c.Clusters[key] = cluster
	}
	for key, clusterNamespace := range c.ClusterNamespaces {
		if err := clusterNamespace.Normalize(resolver); err != nil {
			return err
		}
		c.ClusterNamespaces[key] = clusterNamespace
	}
	for key, ingress := range c.Ingresses {
		if err := ingress.Normalize(resolver); err != nil {
			return err
		}
		c.Ingresses[key] = ingress
	}
	for key, domain := range c.Domains {
		if err := domain.Normalize(resolver); err != nil {
			return err
		}
		c.Domains[key] = domain
	}
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
	if err != nil {
		return nil, core.InvalidYamlError("previews.yml", err)
	}

	for k, ao := range r.Applications {
		ao.Name = k
		r.Applications[k] = ao
	}
	for k, so := range r.Subdomains {
		so.Name = k
		r.Subdomains[k] = so
	}
	for k, do := range r.Datastores {
		do.Name = k
		r.Datastores[k] = do
	}
	for k, bo := range r.Domains {
		bo.Name = k
		r.Domains[k] = bo
	}
	for k, bo := range r.Ingresses {
		bo.Name = k
		r.Ingresses[k] = bo
	}
	for k, bo := range r.ClusterNamespaces {
		bo.Name = k
		r.ClusterNamespaces[k] = bo
	}
	for k, bo := range r.Clusters {
		bo.Name = k
		r.Clusters[k] = bo
	}
	for k, bo := range r.Networks {
		bo.Name = k
		r.Networks[k] = bo
	}
	for k, bo := range r.Blocks {
		bo.Name = k
		r.Blocks[k] = bo
	}
	return r, nil
}
