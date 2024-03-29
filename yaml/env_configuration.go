package yaml

import (
	"gopkg.in/yaml.v3"
)

type EnvConfiguration struct {
	Version           string                                   `yaml:"version" json:"version"`
	Applications      map[string]AppConfiguration              `yaml:"apps,omitempty" json:"apps"`
	Subdomains        map[string]SubdomainConfiguration        `yaml:"subdomains,omitempty" json:"subdomains"`
	Datastores        map[string]DatastoreConfiguration        `yaml:"datastores,omitempty" json:"datastores"`
	Domains           map[string]DomainConfiguration           `yaml:"domains,omitempty" json:"domains"`
	Ingresses         map[string]IngressConfiguration          `yaml:"ingresses,omitempty" json:"ingresses"`
	Networks          map[string]NetworkConfiguration          `yaml:"networks,omitempty" json:"networks"`
	Clusters          map[string]ClusterConfiguration          `yaml:"clusters,omitempty" json:"clusters"`
	ClusterNamespaces map[string]ClusterNamespaceConfiguration `yaml:"cluster_namespaces,omitempty" json:"clusterNamespaces"`
	Blocks            map[string]BlockConfiguration            `yaml:"blocks,omitempty" json:"blocks"`
}

func ParseEnvConfiguration(data []byte) (*EnvConfiguration, error) {
	var r *EnvConfiguration
	err := yaml.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
