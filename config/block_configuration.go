package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"github.com/nullstone-io/module/config"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ core.ChangeApplier = &BlockConfiguration{}
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
	Type                BlockType                `json:"type"`
	Category            types.CategoryName       `json:"category"`
	Name                string                   `json:"name"`
	ModuleSource        string                   `json:"moduleSource"`
	ModuleSourceVersion string                   `json:"moduleSourceVersion"`
	Variables           VariableConfigurations   `json:"vars"`
	Connections         ConnectionConfigurations `json:"connections"`
	IsShared            bool                     `json:"isShared"`

	// These fields are populated via Resolve()
	Module        *types.Module        `json:"module"`
	ModuleVersion *types.ModuleVersion `json:"moduleVersion"`
}

func convertVariables(parsed map[string]any) VariableConfigurations {
	result := VariableConfigurations{}
	for key, value := range parsed {
		result[key] = &VariableConfiguration{Value: value}
	}
	return result
}

func convertConnections(parsed map[string]yaml.ConnectionTarget) ConnectionConfigurations {
	result := ConnectionConfigurations{}
	for key, conn := range parsed {
		result[key] = &ConnectionConfiguration{
			Target: types.ConnectionTarget{
				StackId:   conn.StackId,
				StackName: conn.StackName,
				BlockId:   conn.BlockId,
				BlockName: conn.BlockName,
				EnvId:     conn.EnvId,
				EnvName:   conn.EnvName,
			},
		}
	}
	return result
}

func convertBlockConfigurations(parsed map[string]yaml.BlockConfiguration) map[string]*BlockConfiguration {
	result := map[string]*BlockConfiguration{}
	for blockName, blockValue := range parsed {
		result[blockName] = blockConfigFromYaml(blockName, blockValue, BlockTypeBlock, types.CategoryBlock)
	}
	return result
}

func blockConfigFromYaml(name string, value yaml.BlockConfiguration, blockType BlockType, blockCategory types.CategoryName) *BlockConfiguration {
	// set a default module version if not provided
	moduleVersion := ""
	if value.ModuleSourceVersion != nil {
		moduleVersion = *value.ModuleSourceVersion
	} else if value.ModuleSource != "" {
		moduleVersion = "latest"
	}
	return &BlockConfiguration{
		Type:                blockType,
		Category:            blockCategory,
		Name:                name,
		ModuleSource:        value.ModuleSource,
		ModuleSourceVersion: moduleVersion,
		Variables:           convertVariables(value.Variables),
		Connections:         convertConnections(value.Connections),
		IsShared:            value.IsShared,
	}
}

func (b *BlockConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	if b.Variables == nil {
		b.Variables = VariableConfigurations{}
	}
	if b.Connections == nil {
		b.Connections = ConnectionConfigurations{}
	}
	if ic.IsOverrides && b.ModuleSource == "" {
		// TODO: Add support for validating variables and connections in an overrides file that has no module source
		return nil
	}

	errs := core.ResolveErrors{}

	contract := types.ModuleContractName{Category: string(b.Category), Provider: "*", Platform: "*"}
	manifest := config.Manifest{Variables: map[string]config.Variable{}, Connections: map[string]config.Connection{}}
	m, mv, err := core.ResolveModule(ctx, resolver, pc, b.ModuleSource, b.ModuleSourceVersion, contract)
	if err != nil {
		errs = append(errs, *err)
	} else {
		b.Module = m
		b.ModuleVersion = mv
		manifest = mv.Manifest
	}

	errs = append(errs, b.Variables.Resolve(manifest)...)
	errs = append(errs, b.Connections.Resolve(ctx, resolver, pc, manifest)...)
	return errs
}

func (b *BlockConfiguration) Validate(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	if b.Module == nil {
		// TODO: Add support for validating variables and connections in an overrides file that has no module source
		return nil
	}

	moduleName := fmt.Sprintf("%s/%s@%s", b.Module.OrgName, b.Module.Name, b.ModuleVersion.Version)

	errs := core.ValidateErrors{}
	errs = append(errs, b.Variables.Validate(pc, moduleName)...)
	errs = append(errs, b.Connections.Validate(pc, moduleName)...)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (b *BlockConfiguration) Normalize(ctx context.Context, pc core.ObjectPathContext, resolver core.ConnectionResolver) core.NormalizeErrors {
	return b.Connections.Normalize(ctx, pc, resolver)
}

func (b *BlockConfiguration) ToBlock(orgName string, stackId int64) types.Block {
	block := types.Block{
		Type:                string(b.Type),
		OrgName:             orgName,
		StackId:             stackId,
		Name:                b.Name,
		IsShared:            b.IsShared,
		DnsName:             "",
		ModuleSource:        b.ModuleSource,
		ModuleSourceVersion: b.ModuleSourceVersion,
		Connections:         b.Connections.Targets(),
	}
	for k, conn := range block.Connections {
		if conn.StackId == 0 && conn.StackName == "" {
			conn.StackId = stackId
		}
		block.Connections[k] = conn
	}
	return block
}

func (b *BlockConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	updater.UpdateSchema(b.ModuleSource, b.ModuleVersion)
	for name, vc := range b.Variables {
		updater.UpdateVariableValue(name, vc.Value)
	}
	for name, cc := range b.Connections {
		updater.UpdateConnectionTarget(name, cc.Target)
	}
	return nil
}
