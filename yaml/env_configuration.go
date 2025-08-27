package yaml

import (
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"gopkg.in/yaml.v3"
)

type EnvConfiguration struct {
	Version string `yaml:"version" json:"version"`

	Events EventConfigurations `yaml:"events,omitempty" json:"events"`

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

// EnvConfigurationFromWorkspaceConfig is used to generate an IaC configuration from types.WorkspaceConfig
// Deprecated - This needs reworked
func EnvConfigurationFromWorkspaceConfig(stackId, envId int64, block types.Block, config types.WorkspaceConfig) EnvConfiguration {
	result := EnvConfiguration{Version: "0.1"}

	switch block.Type {
	case string(types.BlockTypeApplication):
		if result.Applications == nil {
			result.Applications = map[string]AppConfiguration{}
		}
		result.Applications[block.Name] = AppConfigurationFromWorkspaceConfig(stackId, envId, config)
	case string(types.BlockTypeDatastore):
		if result.Datastores == nil {
			result.Datastores = map[string]DatastoreConfiguration{}
		}
		result.Datastores[block.Name] = DatastoreConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config)}
	case string(types.BlockTypeSubdomain):
		if result.Subdomains == nil {
			result.Subdomains = map[string]SubdomainConfiguration{}
		}
		result.Subdomains[block.Name] = SubdomainConfigurationFromWorkspaceConfig(stackId, envId, config)
	case string(types.BlockTypeDomain):
		if result.Domains == nil {
			result.Domains = map[string]DomainConfiguration{}
		}
		result.Domains[block.Name] = DomainConfigurationFromWorkspaceConfig(stackId, envId, config)
	case string(types.BlockTypeIngress):
		if result.Ingresses == nil {
			result.Ingresses = map[string]IngressConfiguration{}
		}
		result.Ingresses[block.Name] = IngressConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config)}
	case string(types.BlockTypeClusterNamespace):
		if result.ClusterNamespaces == nil {
			result.ClusterNamespaces = map[string]ClusterNamespaceConfiguration{}
		}
		result.ClusterNamespaces[block.Name] = ClusterNamespaceConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config)}
	case string(types.BlockTypeCluster):
		if result.Clusters == nil {
			result.Clusters = map[string]ClusterConfiguration{}
		}
		result.Clusters[block.Name] = ClusterConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config)}
	case string(types.BlockTypeNetwork):
		if result.Networks == nil {
			result.Networks = map[string]NetworkConfiguration{}
		}
		result.Networks[block.Name] = NetworkConfiguration{BlockConfiguration: BlockConfigurationFromWorkspaceConfig(stackId, envId, config)}
	default:
		if result.Blocks == nil {
			result.Blocks = map[string]BlockConfiguration{}
		}
		result.Blocks[block.Name] = BlockConfigurationFromWorkspaceConfig(stackId, envId, config)
	}

	return result
}
