package config

import (
	"context"
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
		network := NetworkConfiguration{
			BlockConfiguration: blockConfigFromYaml(networkName, networkValue.BlockConfiguration, BlockTypeNetwork),
		}
		result[networkName] = network
	}
	return result
}

func (n NetworkConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, repoName, filename string) error {
	yamlPath := fmt.Sprintf("networks.%s", n.Name)
	contract := fmt.Sprintf("network/*/*")
	return ValidateBlock(ctx, resolver, repoName, filename, yamlPath, contract, n.ModuleSource, n.ModuleSourceVersion, n.Variables, n.Connections, nil, nil)
}
