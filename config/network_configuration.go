package config

import (
	"fmt"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type NetworkConfiguration struct {
	BlockConfiguration
}

func convertNetworkConfigurations(parsed map[string]yaml.NetworkConfiguration) map[string]NetworkConfiguration {
	result := make(map[string]NetworkConfiguration)
	for networkName, networkValue := range parsed {
		// set a default module version if not provided
		moduleVersion := "latest"
		if networkValue.ModuleSourceVersion != nil {
			moduleVersion = *networkValue.ModuleSourceVersion
		}
		network := NetworkConfiguration{
			BlockConfiguration: BlockConfiguration{
				Type:                BlockTypeNetwork,
				Name:                networkName,
				ModuleSource:        networkValue.ModuleSource,
				ModuleSourceVersion: moduleVersion,
				Variables:           networkValue.Variables,
				Connections:         convertConnections(networkValue.Connections),
			},
		}
		result[networkName] = network
	}
	return result
}

func (n NetworkConfiguration) Validate(resolver *find.ResourceResolver, configBlocks []BlockConfiguration) error {
	yamlPath := fmt.Sprintf("networks.%s", n.Name)
	contract := fmt.Sprintf("network/*/*")
	return ValidateBlock(resolver, configBlocks, yamlPath, contract, n.ModuleSource, n.ModuleSourceVersion, n.Variables, n.Connections, nil)
}
