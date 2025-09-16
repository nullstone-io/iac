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
	Type             BlockType                `json:"type"`
	Category         types.CategoryName       `json:"category"`
	Name             string                   `json:"name"`
	ModuleSource     string                   `json:"moduleSource"`
	ModuleConstraint string                   `json:"moduleConstraint"`
	Variables        VariableConfigurations   `json:"vars"`
	Connections      ConnectionConfigurations `json:"connections"`
	IsShared         bool                     `json:"isShared"`

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

func convertConnections(parsed map[string]yaml.ConnectionConstraint) ConnectionConfigurations {
	result := ConnectionConfigurations{}
	for key, conn := range parsed {
		result[key] = &ConnectionConfiguration{
			DesiredTarget: types.ConnectionTarget{
				StackName: conn.StackName,
				BlockName: conn.BlockName,
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
	moduleConstraint := ""
	if value.ModuleConstraint != nil {
		moduleConstraint = *value.ModuleConstraint
	} else if value.ModuleSource != "" {
		moduleConstraint = "latest"
	}
	return &BlockConfiguration{
		Type:             blockType,
		Category:         blockCategory,
		Name:             name,
		ModuleSource:     value.ModuleSource,
		ModuleConstraint: moduleConstraint,
		Variables:        convertVariables(value.Variables),
		Connections:      convertConnections(value.Connections),
		IsShared:         value.IsShared,
	}
}

func (b *BlockConfiguration) Initialize(ctx context.Context, resolver core.InitializeResolver, ic core.IacContext, pc core.ObjectPathContext) core.InitializeErrors {
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

	errs := core.InitializeErrors{}

	contract := types.ModuleContractName{Category: string(b.Category), Provider: "*", Platform: "*"}
	manifest := config.Manifest{Variables: map[string]config.Variable{}, Connections: map[string]config.Connection{}}
	m, mv, err := core.GetModuleVersion(ctx, resolver, pc, b.ModuleSource, b.ModuleConstraint, contract)
	if err != nil {
		errs = append(errs, *err)
	} else {
		b.Module = m
		b.ModuleVersion = mv
		manifest = mv.Manifest
	}

	errs = append(errs, b.Variables.Initialize(manifest)...)
	errs = append(errs, b.Connections.Initialize(ctx, ic, pc, manifest)...)
	return errs
}

func (b *BlockConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver, finder core.IacFinder, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	errs := core.ResolveErrors{}
	errs = append(errs, b.Connections.Resolve(ctx, resolver, finder, ic, pc)...)
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
		Type:     string(b.Type),
		OrgName:  orgName,
		StackId:  stackId,
		Name:     b.Name,
		IsShared: b.IsShared,
		DnsName:  "",
	}
	return block
}

func (b *BlockConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	updater.UpdateSchema(b.ModuleSource, b.ModuleConstraint, b.ModuleVersion)
	for name, vc := range b.Variables {
		updater.UpdateVariableValue(name, vc.Value)
	}
	for name, cc := range b.Connections {
		updater.UpdateConnectionTarget(name, cc.DesiredTarget, cc.EffectiveTarget)
	}
	return nil
}
