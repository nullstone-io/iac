package yaml

import (
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"gopkg.in/yaml.v3"
)

type EnvConfiguration struct {
	Version           string                                   `yaml:"version" json:"version"`
	Applications      map[string]AppConfiguration              `yaml:"apps,omitempty" json:"apps"`
	Blocks            map[string]BlockConfiguration            `yaml:"blocks,omitempty" json:"blocks"`
	Clusters          map[string]ClusterConfiguration          `yaml:"clusters,omitempty" json:"clusters"`
	ClusterNamespaces map[string]ClusterNamespaceConfiguration `yaml:"cluster_namespaces,omitempty" json:"clusterNamespaces"`
	Datastores        map[string]DatastoreConfiguration        `yaml:"datastores,omitempty" json:"datastores"`
	Domains           map[string]DomainConfiguration           `yaml:"domains,omitempty" json:"domains"`
	Ingresses         map[string]IngressConfiguration          `yaml:"ingresses,omitempty" json:"ingresses"`
	Networks          map[string]NetworkConfiguration          `yaml:"networks,omitempty" json:"networks"`
	Subdomains        map[string]SubdomainConfiguration        `yaml:"subdomains,omitempty" json:"subdomains"`
}

func ParseEnvConfiguration(data []byte) (*EnvConfiguration, error) {
	var r *EnvConfiguration
	err := yaml.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func EnvConfigurationFromWorkspaceConfig(block types.Block, config types.WorkspaceConfig) EnvConfiguration {
	result := EnvConfiguration{}
	result.Version = "0.1"

	switch block.Type {
	case string(types.BlockTypeApplication):
		result.Applications[block.Name] = AppConfigurationFromWorkspaceConfig(config)
	case string(types.BlockTypeDatastore):
		result.Datastores[block.Name] = DatastoreConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(config)}
	case string(types.BlockTypeSubdomain):
		result.Subdomains[block.Name] = SubdomainConfigurationFromWorkspaceConfig(block, config)
	case string(types.BlockTypeDomain):
		result.Domains[block.Name] = DomainConfigurationFromWorkspaceConfig(block, config)
	case string(types.BlockTypeIngress):
		result.Ingresses[block.Name] = IngressConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(config)}
	case string(types.BlockTypeClusterNamespace):
		result.ClusterNamespaces[block.Name] = ClusterNamespaceConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(config)}
	case string(types.BlockTypeCluster):
		result.Clusters[block.Name] = ClusterConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(config)}
	case string(types.BlockTypeNetwork):
		result.Networks[block.Name] = NetworkConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(config)}
	default:
		result.Blocks[block.Name] = BlockConfigurationFromWorkspaceConfig(config)
	}

	return result
}
