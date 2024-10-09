package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
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

func (n NetworkConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, ic core.IacContext) errors.ValidationErrors {
	pc := core.NewYamlPathContext("networks", n.Name)
	contract := fmt.Sprintf("network/*/*")
	return ValidateBlock(ctx, resolver, ic, pc, contract, n.ModuleSource, n.ModuleSourceVersion, n.Variables, n.Connections, nil, nil)
}
