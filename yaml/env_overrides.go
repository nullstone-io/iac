package yaml

import (
	"gopkg.in/yaml.v3"
)

type EnvOverrides struct {
	Version           string                               `yaml:"version"`
	Applications      map[string]AppOverrides              `yaml:"apps,omitempty"`
	Subdomains        map[string]SubdomainOverrides        `yaml:"subdomains,omitempty"`
	Datastores        map[string]DatastoreOverrides        `yaml:"datastores,omitempty"`
	Domains           map[string]DomainOverrides           `yaml:"domains,omitempty"`
	Ingresses         map[string]IngressOverrides          `yaml:"ingresses,omitempty"`
	ClusterNamespaces map[string]ClusterNamespaceOverrides `yaml:"cluster_namespaces,omitempty"`
	Clusters          map[string]ClusterOverrides          `yaml:"clusters,omitempty"`
	Networks          map[string]NetworkOverrides          `yaml:"networks,omitempty"`
	Blocks            map[string]BlockOverrides            `yaml:"blocks,omitempty"`
}

func ParseEnvOverrides(data []byte) (*EnvOverrides, error) {
	var r *EnvOverrides
	err := yaml.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
