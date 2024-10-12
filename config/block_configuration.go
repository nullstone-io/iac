package config

import (
	"context"
	"fmt"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
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
		Variables:           value.Variables,
		Connections:         convertConnections(value.Connections),
		IsShared:            value.IsShared,
	}
}

func (b *BlockConfiguration) Resolve(ctx context.Context, resolver core.ModuleVersionResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	if ic.IsOverrides && b.ModuleSource == "" {
		// TODO: Add support for validating variables and connections in an overrides file that has no module source
		return nil
	}

	contract := types.ModuleContractName{Category: string(b.Category), Provider: "*", Platform: "*"}
	m, mv, err := core.ResolveModule(ctx, resolver, pc, b.ModuleSource, b.ModuleSourceVersion, contract)
	if err != nil {
		return core.ResolveErrors{*err}
	}
	b.Module = m
	b.ModuleVersion = mv
	return nil
}

func (b *BlockConfiguration) Validate(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	if b.Module == nil {
		// TODO: Add support for validating variables and connections in an overrides file that has no module source
		return nil
	}

	errs := core.ValidateErrors{}
	errs = append(errs, b.ValidateVariables(ic, pc)...)
	errs = append(errs, b.ValidateConnections(ctx, resolver, ic, pc, b.Connections)...)
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (b *BlockConfiguration) ValidateVariables(ic core.IacContext, pc core.ObjectPathContext) core.ValidateErrors {
	moduleName := fmt.Sprintf("%s/%s@%s", b.Module.OrgName, b.Module.Name, b.ModuleVersion.Version)
	return core.ValidateVariables(pc, b.Variables, b.ModuleVersion.Manifest.Variables, moduleName)
}

// ValidateConnections performs validation on all IaC connections by matching them against connections in the module
func (b *BlockConfiguration) ValidateConnections(ctx context.Context, resolver core.ValidateResolver, ic core.IacContext,
	pc core.ObjectPathContext, connections types.ConnectionTargets) core.ValidateErrors {
	moduleName := fmt.Sprintf("%s/%s@%s", b.Module.OrgName, b.Module.Name, b.ModuleVersion.Version)
	return core.ValidateConnections(ctx, resolver, pc, connections, b.ModuleVersion.Manifest.Connections, moduleName)
}

func (b *BlockConfiguration) Normalize(ctx context.Context, resolver core.ConnectionResolver) error {
	return NormalizeConnectionTargets(ctx, b.Connections, resolver)
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
		Connections:         b.Connections,
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
	for name, value := range b.Variables {
		updater.UpdateVariableValue(name, value)
	}
	for name, value := range b.Connections {
		updater.UpdateConnectionTarget(name, value)
	}
	return nil
}
