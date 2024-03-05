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
	Type                BlockType
	Name                string
	ModuleSource        string
	ModuleSourceVersion string
	Variables           map[string]any
	Connections         types.ConnectionTargets
	IsShared            bool
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
		// set a default module version if not provided
		moduleVersion := "latest"
		if blockValue.ModuleSourceVersion != nil {
			moduleVersion = *blockValue.ModuleSourceVersion
		}
		block := BlockConfiguration{
			Type:                BlockTypeBlock,
			Name:                blockName,
			ModuleSource:        blockValue.ModuleSource,
			ModuleSourceVersion: moduleVersion,
			Variables:           blockValue.Variables,
			Connections:         convertConnections(blockValue.Connections),
			IsShared:            blockValue.IsShared,
		}
		result[blockName] = block
	}
	return result
}

func (b BlockConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver, repoName, filename string) error {
	yamlPath := fmt.Sprintf("blocks.%s", b.Name)
	contract := fmt.Sprintf("block/*/*")
	return ValidateBlock(ctx, resolver, repoName, filename, yamlPath, contract, b.ModuleSource, b.ModuleSourceVersion, b.Variables, b.Connections, nil, nil)
}

func (b *BlockConfiguration) Normalize(ctx context.Context, resolver *find.ResourceResolver) error {
	return core.NormalizeConnectionTargets(ctx, b.Connections, resolver)
}
