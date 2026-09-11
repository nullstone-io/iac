package config

import (
	"context"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

type EnvConfiguration struct {
	IacContext core.IacContext `json:"iacContext"`

	Events EventConfigurations `json:"events"`

	Applications      map[string]*AppConfiguration              `json:"applications"`
	Datastores        map[string]*DatastoreConfiguration        `json:"datastores"`
	Subdomains        map[string]*SubdomainConfiguration        `json:"subdomains"`
	Domains           map[string]*DomainConfiguration           `json:"domains"`
	Ingresses         map[string]*IngressConfiguration          `json:"ingresses"`
	ClusterNamespaces map[string]*ClusterNamespaceConfiguration `json:"clusterNamespaces"`
	Clusters          map[string]*ClusterConfiguration          `json:"clusters"`
	Networks          map[string]*NetworkConfiguration          `json:"networks"`
	Blocks            map[string]*BlockConfiguration            `json:"blocks"`
}

func ConvertConfiguration(repoUrl, repoName, filename string, isOverrides bool, parsed yaml.EnvConfiguration) *EnvConfiguration {
	result := &EnvConfiguration{
		IacContext: core.IacContext{
			RepoUrl:     repoUrl,
			RepoName:    repoName,
			Filename:    filename,
			IsOverrides: isOverrides,
			Version:     parsed.Version,
		},
	}
	result.Events = convertEventConfigurations(parsed.Events)
	result.Applications = convertAppConfigurations(parsed.Applications)
	result.Blocks = convertBlockConfigurations(parsed.Blocks)
	result.Clusters = convertClusterConfigurations(parsed.Clusters)
	result.ClusterNamespaces = convertClusterNamespaceConfigurations(parsed.ClusterNamespaces)
	result.Datastores = convertDatastoreConfigurations(parsed.Datastores)
	result.Domains = convertDomainConfigurations(parsed.Domains)
	result.Ingresses = convertIngressConfigurations(parsed.Ingresses)
	result.Networks = convertNetworkConfigurations(parsed.Networks)
	result.Subdomains = convertSubdomainConfigurations(parsed.Subdomains)
	return result
}

func (e *EnvConfiguration) Initialize(ctx context.Context, resolver core.InitializeResolver) core.InitializeErrors {
	errs := core.InitializeErrors{}

	for _, app := range e.Applications {
		pc := core.NewObjectPathContextKey("apps", app.Name)
		errs = append(errs, app.Initialize(ctx, resolver, e.IacContext, pc)...)
	}
	for _, block := range e.Blocks {
		pc := core.NewObjectPathContextKey("blocks", block.Name)
		errs = append(errs, block.Initialize(ctx, resolver, e.IacContext, pc)...)
	}
	for _, cluster := range e.Clusters {
		pc := core.NewObjectPathContextKey("clusters", cluster.Name)
		errs = append(errs, cluster.Initialize(ctx, resolver, e.IacContext, pc)...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		pc := core.NewObjectPathContextKey("cluster_namespaces", clusterNamespace.Name)
		errs = append(errs, clusterNamespace.Initialize(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ds := range e.Datastores {
		pc := core.NewObjectPathContextKey("datastores", ds.Name)
		errs = append(errs, ds.Initialize(ctx, resolver, e.IacContext, pc)...)
	}
	for _, domain := range e.Domains {
		pc := core.NewObjectPathContextKey("domains", domain.Name)
		errs = append(errs, domain.Initialize(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ingress := range e.Ingresses {
		pc := core.NewObjectPathContextKey("ingresses", ingress.Name)
		errs = append(errs, ingress.Initialize(ctx, resolver, e.IacContext, pc)...)
	}
	for _, network := range e.Networks {
		pc := core.NewObjectPathContextKey("networks", network.Name)
		errs = append(errs, network.Initialize(ctx, resolver, e.IacContext, pc)...)
	}
	for _, sub := range e.Subdomains {
		pc := core.NewObjectPathContextKey("subdomains", sub.Name)
		errs = append(errs, sub.Initialize(ctx, resolver, e.IacContext, pc)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (e *EnvConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver, finder core.IacFinder) core.ResolveErrors {
	errs := core.ResolveErrors{}

	for name, evt := range e.Events {
		pc := core.NewObjectPathContextKey("events", name)
		errs = append(errs, evt.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}

	for _, app := range e.Applications {
		pc := core.NewObjectPathContextKey("apps", app.Name)
		errs = append(errs, app.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}
	for _, block := range e.Blocks {
		pc := core.NewObjectPathContextKey("blocks", block.Name)
		errs = append(errs, block.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}
	for _, cluster := range e.Clusters {
		pc := core.NewObjectPathContextKey("clusters", cluster.Name)
		errs = append(errs, cluster.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		pc := core.NewObjectPathContextKey("cluster_namespaces", clusterNamespace.Name)
		errs = append(errs, clusterNamespace.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}
	for _, ds := range e.Datastores {
		pc := core.NewObjectPathContextKey("datastores", ds.Name)
		errs = append(errs, ds.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}
	for _, domain := range e.Domains {
		pc := core.NewObjectPathContextKey("domains", domain.Name)
		errs = append(errs, domain.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}
	for _, ingress := range e.Ingresses {
		pc := core.NewObjectPathContextKey("ingresses", ingress.Name)
		errs = append(errs, ingress.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}
	for _, network := range e.Networks {
		pc := core.NewObjectPathContextKey("networks", network.Name)
		errs = append(errs, network.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}
	for _, sub := range e.Subdomains {
		pc := core.NewObjectPathContextKey("subdomains", sub.Name)
		errs = append(errs, sub.Resolve(ctx, resolver, finder, e.IacContext, pc)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (e *EnvConfiguration) Validate() core.ValidateErrors {
	errs := core.ValidateErrors{}

	for name, event := range e.Events {
		pc := core.NewObjectPathContextKey("events", name)
		errs = append(errs, event.Validate(e.IacContext, pc)...)
	}

	for _, app := range e.Applications {
		pc := core.NewObjectPathContextKey("apps", app.Name)
		errs = append(errs, app.Validate(e.IacContext, pc)...)
	}
	for _, block := range e.Blocks {
		pc := core.NewObjectPathContextKey("blocks", block.Name)
		errs = append(errs, block.Validate(e.IacContext, pc)...)
	}
	for _, cluster := range e.Clusters {
		pc := core.NewObjectPathContextKey("clusters", cluster.Name)
		errs = append(errs, cluster.Validate(e.IacContext, pc)...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		pc := core.NewObjectPathContextKey("cluster_namespaces", clusterNamespace.Name)
		errs = append(errs, clusterNamespace.Validate(e.IacContext, pc)...)
	}
	for _, ds := range e.Datastores {
		pc := core.NewObjectPathContextKey("datastores", ds.Name)
		errs = append(errs, ds.Validate(e.IacContext, pc)...)
	}
	for _, domain := range e.Domains {
		pc := core.NewObjectPathContextKey("domains", domain.Name)
		errs = append(errs, domain.Validate(e.IacContext, pc)...)
	}
	for _, ingress := range e.Ingresses {
		pc := core.NewObjectPathContextKey("ingresses", ingress.Name)
		errs = append(errs, ingress.Validate(e.IacContext, pc)...)
	}
	for _, network := range e.Networks {
		pc := core.NewObjectPathContextKey("networks", network.Name)
		errs = append(errs, network.Validate(e.IacContext, pc)...)
	}
	for _, sub := range e.Subdomains {
		pc := core.NewObjectPathContextKey("subdomains", sub.Name)
		errs = append(errs, sub.Validate(e.IacContext, pc)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (e *EnvConfiguration) Normalize(ctx context.Context, resolver core.NormalizeResolver) core.NormalizeErrors {
	errs := core.NormalizeErrors{}

	for _, event := range e.Events {
		pc := core.NewObjectPathContextKey("events", event.Name)
		errs = append(errs, event.Normalize(ctx, pc, resolver)...)
	}

	for _, app := range e.Applications {
		pc := core.NewObjectPathContextKey("apps", app.Name)
		errs = append(errs, app.Normalize(ctx, pc, resolver)...)
	}
	for _, block := range e.Blocks {
		pc := core.NewObjectPathContextKey("blocks", block.Name)
		errs = append(errs, block.Normalize(ctx, pc, resolver)...)
	}
	for _, cluster := range e.Clusters {
		pc := core.NewObjectPathContextKey("clusters", cluster.Name)
		errs = append(errs, cluster.Normalize(ctx, pc, resolver)...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		pc := core.NewObjectPathContextKey("cluster_namespaces", clusterNamespace.Name)
		errs = append(errs, clusterNamespace.Normalize(ctx, pc, resolver)...)
	}
	for _, ds := range e.Datastores {
		pc := core.NewObjectPathContextKey("datastores", ds.Name)
		errs = append(errs, ds.Normalize(ctx, pc, resolver)...)
	}
	for _, domain := range e.Domains {
		pc := core.NewObjectPathContextKey("domains", domain.Name)
		errs = append(errs, domain.Normalize(ctx, pc, resolver)...)
	}
	for _, ingress := range e.Ingresses {
		pc := core.NewObjectPathContextKey("ingresses", ingress.Name)
		errs = append(errs, ingress.Normalize(ctx, pc, resolver)...)
	}
	for _, network := range e.Networks {
		pc := core.NewObjectPathContextKey("networks", network.Name)
		errs = append(errs, network.Normalize(ctx, pc, resolver)...)
	}
	for _, sub := range e.Subdomains {
		pc := core.NewObjectPathContextKey("subdomains", sub.Name)
		errs = append(errs, sub.Normalize(ctx, pc, resolver)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// BlockDefinition is a block as declared in config.yml. Block carries only the fields IaC
// governs; anything IaC does not declare is left at its zero value. It is a struct so that
// declared-only information which the flat types.Block cannot express can be added later.
type BlockDefinition struct {
	Block types.Block
}

// ToBlockDefinitions flattens every block in the file into sync definitions.
// Callers that only need the API shape read Block off each definition.
func (e *EnvConfiguration) ToBlockDefinitions(orgName string, stackId int64) []BlockDefinition {
	defs := make([]BlockDefinition, 0)
	if e == nil {
		return defs
	}

	for _, app := range e.Applications {
		defs = append(defs, app.toBlockDefinition(orgName, stackId))
	}
	for _, ds := range e.Datastores {
		defs = append(defs, ds.toBlockDefinition(orgName, stackId))
	}
	for _, sub := range e.Subdomains {
		defs = append(defs, sub.toBlockDefinition(orgName, stackId))
	}
	for _, d := range e.Domains {
		defs = append(defs, d.toBlockDefinition(orgName, stackId))
	}
	for _, i := range e.Ingresses {
		defs = append(defs, i.toBlockDefinition(orgName, stackId))
	}
	for _, cn := range e.ClusterNamespaces {
		defs = append(defs, cn.toBlockDefinition(orgName, stackId))
	}
	for _, c := range e.Clusters {
		defs = append(defs, c.toBlockDefinition(orgName, stackId))
	}
	for _, n := range e.Networks {
		defs = append(defs, n.toBlockDefinition(orgName, stackId))
	}
	for _, b := range e.Blocks {
		defs = append(defs, b.toBlockDefinition(orgName, stackId))
	}

	return defs
}

func (e *EnvConfiguration) ApplyChangesTo(block types.Block, updater core.WorkspaceConfigUpdater) error {
	var ca core.ChangeApplier
	var ok bool
	switch BlockType(block.Type) {
	case BlockTypeApplication:
		ca, ok = e.Applications[block.Name]
	case BlockTypeDomain:
		ca, ok = e.Domains[block.Name]
	case BlockTypeSubdomain:
		ca, ok = e.Subdomains[block.Name]
	case BlockTypeIngress:
		ca, ok = e.Ingresses[block.Name]
	case BlockTypeDatastore:
		ca, ok = e.Datastores[block.Name]
	case BlockTypeClusterNamespace:
		ca, ok = e.ClusterNamespaces[block.Name]
	case BlockTypeCluster:
		ca, ok = e.Clusters[block.Name]
	case BlockTypeNetwork:
		ca, ok = e.Networks[block.Name]
	default:
		ca, ok = e.Blocks[block.Name]
	}
	if !ok {
		return nil
	}
	return ca.ApplyChangesTo(e.IacContext, updater)
}

func (e *EnvConfiguration) BlockNames() map[string]bool {
	names := map[string]bool{}
	for _, cur := range e.Applications {
		names[cur.Name] = true
	}
	for _, cur := range e.Blocks {
		names[cur.Name] = true
	}
	for _, cur := range e.Clusters {
		names[cur.Name] = true
	}
	for _, cur := range e.ClusterNamespaces {
		names[cur.Name] = true
	}
	for _, cur := range e.Datastores {
		names[cur.Name] = true
	}
	for _, cur := range e.Domains {
		names[cur.Name] = true
	}
	for _, cur := range e.Ingresses {
		names[cur.Name] = true
	}
	for _, cur := range e.Networks {
		names[cur.Name] = true
	}
	for _, cur := range e.Subdomains {
		names[cur.Name] = true
	}
	return names
}

func (e *EnvConfiguration) FindBlockConfigurationByName(name string) *BlockConfiguration {
	ptr := func(cur BlockConfiguration) *BlockConfiguration {
		return &cur
	}

	for _, cur := range e.Applications {
		if cur.Name == name {
			return ptr(cur.BlockConfiguration)
		}
	}
	for _, cur := range e.Blocks {
		if cur.Name == name {
			return ptr(*cur)
		}
	}
	for _, cur := range e.Clusters {
		if cur.Name == name {
			return ptr(cur.BlockConfiguration)
		}
	}
	for _, cur := range e.ClusterNamespaces {
		if cur.Name == name {
			return ptr(cur.BlockConfiguration)
		}
	}
	for _, cur := range e.Datastores {
		if cur.Name == name {
			return ptr(cur.BlockConfiguration)
		}
	}
	for _, cur := range e.Domains {
		if cur.Name == name {
			return ptr(cur.BlockConfiguration)
		}
	}
	for _, cur := range e.Ingresses {
		if cur.Name == name {
			return ptr(cur.BlockConfiguration)
		}
	}
	for _, cur := range e.Networks {
		if cur.Name == name {
			return ptr(cur.BlockConfiguration)
		}
	}
	for _, cur := range e.Subdomains {
		if cur.Name == name {
			return ptr(cur.BlockConfiguration)
		}
	}
	return nil
}
