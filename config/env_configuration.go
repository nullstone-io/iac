package config

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
)

type EnvConfiguration struct {
	IacContext        core.IacContext                           `json:"iacContext"`
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
		},
	}
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

func (e *EnvConfiguration) Resolve(ctx context.Context, resolver core.ModuleVersionResolver) core.ResolveErrors {
	errs := core.ResolveErrors{}

	for _, app := range e.Applications {
		pc := core.NewObjectPathContext("apps", app.Name)
		errs = append(errs, app.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ds := range e.Datastores {
		pc := core.NewObjectPathContext("networks", ds.Name)
		errs = append(errs, ds.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, sub := range e.Subdomains {
		pc := core.NewObjectPathContext("subdomains", sub.Name)
		errs = append(errs, sub.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, domain := range e.Domains {
		pc := core.NewObjectPathContext("domains", domain.Name)
		errs = append(errs, domain.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ingress := range e.Ingresses {
		pc := core.NewObjectPathContext("ingresses", ingress.Name)
		errs = append(errs, ingress.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		pc := core.NewObjectPathContext("cluster_namespaces", clusterNamespace.Name)
		errs = append(errs, clusterNamespace.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, cluster := range e.Clusters {
		pc := core.NewObjectPathContext("clusters", cluster.Name)
		errs = append(errs, cluster.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, network := range e.Networks {
		pc := core.NewObjectPathContext("networks", network.Name)
		errs = append(errs, network.Resolve(ctx, resolver, e.IacContext, pc)...)
	}
	for _, block := range e.Blocks {
		pc := core.NewObjectPathContext("blocks", block.Name)
		errs = append(errs, block.Resolve(ctx, resolver, e.IacContext, pc)...)
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (e *EnvConfiguration) Validate(ctx context.Context, resolver core.ValidateResolver) errors.ValidationErrors {
	ve := errors.ValidationErrors{}

	for _, app := range e.Applications {
		pc := core.NewObjectPathContext("apps", app.Name)
		ve = append(ve, app.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ds := range e.Datastores {
		pc := core.NewObjectPathContext("networks", ds.Name)
		ve = append(ve, ds.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, sub := range e.Subdomains {
		pc := core.NewObjectPathContext("subdomains", sub.Name)
		ve = append(ve, sub.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, domain := range e.Domains {
		pc := core.NewObjectPathContext("domains", domain.Name)
		ve = append(ve, domain.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ingress := range e.Ingresses {
		pc := core.NewObjectPathContext("ingresses", ingress.Name)
		ve = append(ve, ingress.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		pc := core.NewObjectPathContext("cluster_namespaces", clusterNamespace.Name)
		ve = append(ve, clusterNamespace.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, cluster := range e.Clusters {
		pc := core.NewObjectPathContext("clusters", cluster.Name)
		ve = append(ve, cluster.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, network := range e.Networks {
		pc := core.NewObjectPathContext("networks", network.Name)
		ve = append(ve, network.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, block := range e.Blocks {
		pc := core.NewObjectPathContext("blocks", block.Name)
		ve = append(ve, block.Validate(ctx, resolver, e.IacContext, pc)...)
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func (e *EnvConfiguration) Normalize(ctx context.Context, resolver core.ConnectionResolver) error {
	for _, block := range e.Blocks {
		if err := block.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	for _, network := range e.Networks {
		if err := network.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	for _, cluster := range e.Clusters {
		if err := cluster.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		if err := clusterNamespace.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	for _, ingress := range e.Ingresses {
		if err := ingress.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	for _, domain := range e.Domains {
		if err := domain.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	for _, subdomain := range e.Subdomains {
		if err := subdomain.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	for _, datastore := range e.Datastores {
		if err := datastore.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	for _, app := range e.Applications {
		if err := app.Normalize(ctx, resolver); err != nil {
			return err
		}
	}
	return nil
}
