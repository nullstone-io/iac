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

func ConvertConfiguration(repoName, filename string, isOverrides bool, parsed yaml.EnvConfiguration) *EnvConfiguration {
	result := &EnvConfiguration{
		IacContext: core.IacContext{
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

func (e *EnvConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver) core.ResolveErrors {
	errs := core.ResolveErrors{}

	for _, app := range e.Applications {
		pc := core.NewObjectPathContextKey("apps", app.Name)
		errs = append(errs, app.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, block := range e.Blocks {
		pc := core.NewObjectPathContextKey("blocks", block.Name)
		errs = append(errs, block.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, cluster := range e.Clusters {
		pc := core.NewObjectPathContextKey("clusters", cluster.Name)
		errs = append(errs, cluster.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		pc := core.NewObjectPathContextKey("cluster_namespaces", clusterNamespace.Name)
		errs = append(errs, clusterNamespace.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ds := range e.Datastores {
		pc := core.NewObjectPathContextKey("networks", ds.Name)
		errs = append(errs, ds.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, domain := range e.Domains {
		pc := core.NewObjectPathContextKey("domains", domain.Name)
		errs = append(errs, domain.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ingress := range e.Ingresses {
		pc := core.NewObjectPathContextKey("ingresses", ingress.Name)
		errs = append(errs, ingress.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, network := range e.Networks {
		pc := core.NewObjectPathContextKey("networks", network.Name)
		errs = append(errs, network.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, sub := range e.Subdomains {
		pc := core.NewObjectPathContextKey("subdomains", sub.Name)
		errs = append(errs, sub.Resolve(ctx, resolver, e.IacContext, pc)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (e *EnvConfiguration) Validate() core.ValidateErrors {
	errs := core.ValidateErrors{}

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
		pc := core.NewObjectPathContextKey("networks", ds.Name)
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

	for name, event := range e.Events {
		pc := core.NewObjectPathContextKey("events", name)
		errs = append(errs, event.Validate(e.IacContext, pc)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (e *EnvConfiguration) Normalize(ctx context.Context, resolver core.ConnectionResolver) core.NormalizeErrors {
	errs := core.NormalizeErrors{}

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
		pc := core.NewObjectPathContextKey("networks", ds.Name)
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

func (e *EnvConfiguration) ToBlocks(orgName string, stackId int64) types.Blocks {
	blocks := make([]types.Block, 0)
	if e == nil {
		return blocks
	}

	for _, app := range e.Applications {
		blocks = append(blocks, app.ToBlock(orgName, stackId))
	}
	for _, ds := range e.Datastores {
		blocks = append(blocks, ds.ToBlock(orgName, stackId))
	}
	for _, sub := range e.Subdomains {
		blocks = append(blocks, sub.ToBlock(orgName, stackId))
	}
	for _, d := range e.Domains {
		blocks = append(blocks, d.ToBlock(orgName, stackId))
	}
	for _, i := range e.Ingresses {
		blocks = append(blocks, i.ToBlock(orgName, stackId))
	}
	for _, cn := range e.ClusterNamespaces {
		blocks = append(blocks, cn.ToBlock(orgName, stackId))
	}
	for _, c := range e.Clusters {
		blocks = append(blocks, c.ToBlock(orgName, stackId))
	}
	for _, n := range e.Networks {
		blocks = append(blocks, n.ToBlock(orgName, stackId))
	}
	for _, b := range e.Blocks {
		blocks = append(blocks, b.ToBlock(orgName, stackId))
	}

	return blocks
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
