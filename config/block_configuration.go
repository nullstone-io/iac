package config

import (
	"context"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
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
	Category            types.CategoryName      `json:"category"`
	Name                string                  `json:"name"`
	ModuleSource        string                  `json:"moduleSource"`
	ModuleSourceVersion string                  `json:"moduleSourceVersion"`
	Variables           map[string]any          `json:"vars"`
	Connections         types.ConnectionTargets `json:"connections"`
	IsShared            bool                    `json:"isShared"`

	Module        *types.Module        `json:"module"`
	ModuleVersion *types.ModuleVersion `json:"moduleVersion"`
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
		result[blockName] = blockConfigFromYaml(blockName, blockValue, BlockTypeBlock, types.CategoryBlock)
	}
	return result
}

func blockConfigFromYaml(name string, value yaml.BlockConfiguration, blockType BlockType, blockCategory types.CategoryName) BlockConfiguration {
	// set a default module version if not provided
	moduleVersion := ""
	if value.ModuleSourceVersion != nil {
		moduleVersion = *value.ModuleSourceVersion
	} else if value.ModuleSource != "" {
		moduleVersion = "latest"
	}
	return BlockConfiguration{
		Type:                blockType,
		Category:            blockCategory,
		Name:                name,
		ModuleSource:        value.ModuleSource,
		ModuleSourceVersion: moduleVersion,
		Variables:           value.Variables,
		Connections:         convertConnections(value.Connections),
		IsShared:            value.IsShared,
	}
}

func (b *BlockConfiguration) Validate(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.YamlPathContext) errors.ValidationErrors {
	if ic.IsOverrides && b.ModuleSource == "" {
		// TODO: Add support for validating variables and connections in an overrides file that has no module source
		return nil
	}

	contract := types.ModuleContractName{Category: string(b.Category), Provider: "*", Platform: "*"}
	m, mv, err := ResolveModule(ctx, resolver, ic, pc, b.ModuleSource, b.ModuleSourceVersion, contract)
	if err != nil {
		return errors.ValidationErrors{*err}
	}
	b.Module = m
	b.ModuleVersion = mv

	ve := errors.ValidationErrors{}
	ve = append(ve, b.ValidateVariables(ic, pc)...)
	ve = append(ve, b.ValidateConnections(ctx, resolver, ic, pc, b.Connections)...)

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func (b *BlockConfiguration) ValidateVariables(ic core.IacContext, pc core.YamlPathContext) errors.ValidationErrors {
	moduleName := fmt.Sprintf("%s/%s@%s", b.Module.OrgName, b.Module.Name, b.ModuleVersion.Version)
	return ValidateVariables(ic, pc, b.Variables, b.ModuleVersion.Manifest.Variables, moduleName)
}

// ValidateConnections performs validation on all IaC connections by matching them against connections in the module
func (b *BlockConfiguration) ValidateConnections(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext,
	pc core.YamlPathContext, connections types.ConnectionTargets) errors.ValidationErrors {
	moduleName := fmt.Sprintf("%s/%s@%s", b.Module.OrgName, b.Module.Name, b.ModuleVersion.Version)
	return ValidateConnections(ctx, resolver, ic, pc, connections, b.ModuleVersion.Manifest.Connections, moduleName)
}

func (b *BlockConfiguration) Normalize(ctx context.Context, resolver core.ConnectionResolver) error {
	return NormalizeConnectionTargets(ctx, b.Connections, resolver)
}
