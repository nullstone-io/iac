package config

import (
	errs "errors"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type EnvConfiguration struct {
	RepoName          string
	Filename          string
	Applications      map[string]AppConfiguration
	Datastores        map[string]DatastoreConfiguration
	Subdomains        map[string]SubdomainConfiguration
	Domains           map[string]DomainConfiguration
	Ingresses         map[string]IngressConfiguration
	ClusterNamespaces map[string]ClusterNamespaceConfiguration
	Clusters          map[string]ClusterConfiguration
	Networks          map[string]NetworkConfiguration
	Blocks            map[string]BlockConfiguration
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

func (e EnvConfiguration) Validate(resolver *find.ResourceResolver) error {
	ve := errors.ValidationErrors{}
	for _, app := range e.Applications {
		err := app.Validate(resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, ds := range e.Datastores {
		err := ds.Validate(resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, sub := range e.Subdomains {
		err := sub.Validate(resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, domain := range e.Domains {
		err := domain.Validate(resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, ingress := range e.Ingresses {
		err := ingress.Validate(resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		err := clusterNamespace.Validate(resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, cluster := range e.Clusters {
		err := cluster.Validate(resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, network := range e.Networks {
		err := network.Validate(resolver, e.RepoName, e.Filename)
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}
	for _, block := range e.Blocks {
		err := block.Validate(resolver, e.RepoName, e.Filename)
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

func (e *EnvConfiguration) Normalize(resolver *find.ResourceResolver) error {
	for key, block := range e.Blocks {
		if err := block.Normalize(resolver); err != nil {
			return err
		}
		e.Blocks[key] = block
	}
	for key, network := range e.Networks {
		if err := network.Normalize(resolver); err != nil {
			return err
		}
		e.Networks[key] = network
	}
	for key, cluster := range e.Clusters {
		if err := cluster.Normalize(resolver); err != nil {
			return err
		}
		e.Clusters[key] = cluster
	}
	for key, clusterNamespace := range e.ClusterNamespaces {
		if err := clusterNamespace.Normalize(resolver); err != nil {
			return err
		}
		e.ClusterNamespaces[key] = clusterNamespace
	}
	for key, ingress := range e.Ingresses {
		if err := ingress.Normalize(resolver); err != nil {
			return err
		}
		e.Ingresses[key] = ingress
	}
	for key, domain := range e.Domains {
		if err := domain.Normalize(resolver); err != nil {
			return err
		}
		e.Domains[key] = domain
	}
	for key, subdomain := range e.Subdomains {
		if err := subdomain.Normalize(resolver); err != nil {
			return err
		}
		e.Subdomains[key] = subdomain
	}
	for key, datastore := range e.Datastores {
		if err := datastore.Normalize(resolver); err != nil {
			return err
		}
		e.Datastores[key] = datastore
	}
	for key, app := range e.Applications {
		if err := app.Normalize(resolver); err != nil {
			return err
		}
		e.Applications[key] = app
	}
	return nil
}
