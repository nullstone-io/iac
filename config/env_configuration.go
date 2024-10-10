package config

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
)

type EnvConfiguration struct {
	IacContext        core.IacContext                          `json:"iacContext"`
	Applications      map[string]AppConfiguration              `json:"applications"`
	Datastores        map[string]DatastoreConfiguration        `json:"datastores"`
	Subdomains        map[string]SubdomainConfiguration        `json:"subdomains"`
	Domains           map[string]DomainConfiguration           `json:"domains"`
	Ingresses         map[string]IngressConfiguration          `json:"ingresses"`
	ClusterNamespaces map[string]ClusterNamespaceConfiguration `json:"clusterNamespaces"`
	Clusters          map[string]ClusterConfiguration          `json:"clusters"`
	Networks          map[string]NetworkConfiguration          `json:"networks"`
	Blocks            map[string]BlockConfiguration            `json:"blocks"`
}

func ConvertConfiguration(repoName, filename string, isOverrides bool, parsed yaml.EnvConfiguration) EnvConfiguration {
	result := EnvConfiguration{
		IacContext: core.IacContext{
			RepoName:    repoName,
			Filename:    filename,
			IsOverrides: isOverrides,
		},
	}
	result.Applications = convertAppConfigurations(parsed.Applications)
	result.Datastores = convertDatastoreConfigurations(parsed.Datastores)
	result.Subdomains = convertSubdomainConfigurations(parsed.Subdomains)
	result.Domains = convertDomainConfigurations(parsed.Domains)
	result.Ingresses = convertIngressConfigurations(parsed.Ingresses)
	result.ClusterNamespaces = convertClusterNamespaceConfigurations(parsed.ClusterNamespaces)
	result.Clusters = convertClusterConfigurations(parsed.Clusters)
	result.Networks = convertNetworkConfigurations(parsed.Networks)
	result.Blocks = convertBlockConfigurations(parsed.Blocks)
	return result
}

func (e *EnvConfiguration) Validate(ctx context.Context, resolver core.ValidateResolver) errors.ValidationErrors {
	ve := errors.ValidationErrors{}
	for _, app := range e.Applications {
		pc := core.NewYamlPathContext("apps", app.Name)
		ve = append(ve, app.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ds := range e.Datastores {
		pc := core.NewYamlPathContext("networks", ds.Name)
		ve = append(ve, ds.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, sub := range e.Subdomains {
		pc := core.NewYamlPathContext("subdomains", sub.Name)
		ve = append(ve, sub.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, domain := range e.Domains {
		pc := core.NewYamlPathContext("domains", domain.Name)
		ve = append(ve, domain.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, ingress := range e.Ingresses {
		pc := core.NewYamlPathContext("ingresses", ingress.Name)
		ve = append(ve, ingress.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		pc := core.NewYamlPathContext("cluster_namespaces", clusterNamespace.Name)
		ve = append(ve, clusterNamespace.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, cluster := range e.Clusters {
		pc := core.NewYamlPathContext("clusters", cluster.Name)
		ve = append(ve, cluster.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, network := range e.Networks {
		pc := core.NewYamlPathContext("networks", network.Name)
		ve = append(ve, network.Validate(ctx, resolver, e.IacContext, pc)...)
	}
	for _, block := range e.Blocks {
		pc := core.NewYamlPathContext("blocks", block.Name)
		ve = append(ve, block.Validate(ctx, resolver, e.IacContext, pc)...)
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func (e *EnvConfiguration) Normalize(ctx context.Context, resolver core.ConnectionResolver) error {
	for key, block := range e.Blocks {
		if err := block.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.Blocks[key] = block
	}
	for key, network := range e.Networks {
		if err := network.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.Networks[key] = network
	}
	for key, cluster := range e.Clusters {
		if err := cluster.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.Clusters[key] = cluster
	}
	for key, clusterNamespace := range e.ClusterNamespaces {
		if err := clusterNamespace.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.ClusterNamespaces[key] = clusterNamespace
	}
	for key, ingress := range e.Ingresses {
		if err := ingress.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.Ingresses[key] = ingress
	}
	for key, domain := range e.Domains {
		if err := domain.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.Domains[key] = domain
	}
	for key, subdomain := range e.Subdomains {
		if err := subdomain.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.Subdomains[key] = subdomain
	}
	for key, datastore := range e.Datastores {
		if err := datastore.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.Datastores[key] = datastore
	}
	for key, app := range e.Applications {
		if err := app.Normalize(ctx, resolver); err != nil {
			return err
		}
		e.Applications[key] = app
	}
	return nil
}
