package config

import (
	"context"
	errs "errors"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type EnvConfiguration struct {
	RepoName          string                                   `json:"repoName"`
	Filename          string                                   `json:"filename"`
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

func ConvertConfiguration(repoName, filename string, parsed yaml.EnvConfiguration) EnvConfiguration {
	result := EnvConfiguration{RepoName: repoName, Filename: filename}
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

func (e EnvConfiguration) Validate(ctx context.Context, resolver *find.ResourceResolver) error {
	ve := errors.ValidationErrors{}
	for _, app := range e.Applications {
		err := app.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, ds := range e.Datastores {
		err := ds.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, sub := range e.Subdomains {
		err := sub.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, domain := range e.Domains {
		err := domain.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, ingress := range e.Ingresses {
		err := ingress.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		err := clusterNamespace.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, cluster := range e.Clusters {
		err := cluster.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, network := range e.Networks {
		err := network.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, block := range e.Blocks {
		err := block.Validate(ctx, resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			} else {
				return err
			}
		}
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}

func (e *EnvConfiguration) Normalize(ctx context.Context, resolver *find.ResourceResolver) error {
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
