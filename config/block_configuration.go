package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type BlockType string

const (
	BlockTypeApplication      BlockType = "Application"
	BlockTypeDatastore        BlockType = "Datastore"
	BlockTypeSubdomain        BlockType = "Subdomain"
	BlockTypeDomain           BlockType = "Domain"
	BlockTypeIngress          BlockType = "Ingress"
	BlockTypeClusterNamespace BlockType = "ClusterNamespace"
	BlockTypeCluster          BlockType = "Cluster"
	BlockTypeNetwork          BlockType = "Network"
	BlockTypeBlock            BlockType = "Block"
)

type BlockConfiguration struct {
	Type                BlockType               `json:"type"`
	Name                string                  `json:"name"`
	ModuleSource        string                  `json:"moduleSource"`
	ModuleSourceVersion string                  `json:"moduleSourceVersion"`
	Variables           map[string]any          `json:"vars"`
	Connections         types.ConnectionTargets `json:"connections"`
	IsShared            bool                    `json:"isShared"`
}

func convertConnections(parsed map[string]yaml.ConnectionTarget) map[string]types.ConnectionTarget {
	result := make(map[string]types.ConnectionTarget)
	for key, conn := range parsed {
		result[key] = types.ConnectionTarget{
			StackId:   conn.StackId,
			StackName: conn.StackName,
			BlockId:   conn.BlockId,
			BlockName: conn.BlockName,
			EnvId:     conn.EnvId,
			EnvName:   conn.EnvName,
		}
	}
	return result
}

func convertBlockConfigurations(parsed map[string]yaml.BlockConfiguration) map[string]BlockConfiguration {
	result := map[string]BlockConfiguration{}
	for blockName, blockValue := range parsed {
		result[blockName] = blockConfigFromYaml(blockName, blockValue, BlockTypeBlock)
	}
	return result
}

func blockConfigFromYaml(name string, value yaml.BlockConfiguration, blockType BlockType) BlockConfiguration {
	// set a default module version if not provided
	moduleVersion := "latest"
	if value.ModuleSourceVersion != nil {
		moduleVersion = *value.ModuleSourceVersion
	}
	return BlockConfiguration{
		Type:                blockType,
		Name:                name,
		ModuleSource:        value.ModuleSource,
		ModuleSourceVersion: moduleVersion,
		Variables:           value.Variables,
		Connections:         convertConnections(value.Connections),
		IsShared:            value.IsShared,
	}
}

func (b BlockConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, repoName, filename string) error {
	yamlPath := fmt.Sprintf("blocks.%s", b.Name)
	contract := fmt.Sprintf("block/*/*")
	return ValidateBlock(ctx, resolver, repoName, filename, yamlPath, contract, b.ModuleSource, b.ModuleSourceVersion, b.Variables, b.Connections, nil, nil)
}

func (b *BlockConfiguration) Normalize(ctx context.Context, resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(ctx, b.Connections, resolver)
}
