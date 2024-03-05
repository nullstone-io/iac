package overrides

import (
	"context"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
)

type EnvOverrides struct {
	RepoName          string
	Filename          string
	Applications      map[string]AppOverrides
	Subdomains        map[string]SubdomainOverrides
	Datastores        map[string]DatastoreOverrides
	Domains           map[string]DomainOverrides
	Ingresses         map[string]IngressOverrides
	ClusterNamespaces map[string]ClusterNamespaceOverrides
	Clusters          map[string]ClusterOverrides
	Networks          map[string]NetworkOverrides
	Blocks            map[string]BlockOverrides
}

func ConvertOverrides(repoName, filename string, parsed yaml.EnvOverrides) EnvOverrides {
	result := EnvOverrides{RepoName: repoName, Filename: filename}
	result.Applications = convertAppOverrides(parsed.Applications)
	result.Datastores = convertDatastoreOverrides(parsed.Datastores)
	result.Subdomains = convertSubdomainOverrides(parsed.Subdomains)
	result.Domains = convertDomainOverrides(parsed.Domains)
	result.Ingresses = convertIngressOverrides(parsed.Ingresses)
	result.ClusterNamespaces = convertClusterNamespaceOverrides(parsed.ClusterNamespaces)
	result.Clusters = convertClusterOverrides(parsed.Clusters)
	result.Networks = convertNetworkOverrides(parsed.Networks)
	result.Blocks = convertBlockOverrides(parsed.Blocks)
	return result
}

func (c *EnvOverrides) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	ve := errors.ValidationErrors{}
	for _, block := range c.Blocks {
		verrs, err := block.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, network := range c.Networks {
		verrs, err := network.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, cluster := range c.Clusters {
		verrs, err := cluster.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, clusterNamespace := range c.ClusterNamespaces {
		verrs, err := clusterNamespace.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, ingress := range c.Ingresses {
		verrs, err := ingress.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, domain := range c.Domains {
		verrs, err := domain.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, subdomain := range c.Subdomains {
		verrs, err := subdomain.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, datastore := range c.Datastores {
		verrs, err := datastore.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, app := range c.Applications {
		verrs, err := app.Validate(resolver)
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	return ve, nil
}

func (c *EnvOverrides) Normalize(ctx context.Context, resolver *find.ResourceResolver) error {
	for key, block := range c.Blocks {
		if err := block.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.Blocks[key] = block
	}
	for key, network := range c.Networks {
		if err := network.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.Networks[key] = network
	}
	for key, cluster := range c.Clusters {
		if err := cluster.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.Clusters[key] = cluster
	}
	for key, clusterNamespace := range c.ClusterNamespaces {
		if err := clusterNamespace.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.ClusterNamespaces[key] = clusterNamespace
	}
	for key, ingress := range c.Ingresses {
		if err := ingress.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.Ingresses[key] = ingress
	}
	for key, domain := range c.Domains {
		if err := domain.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.Domains[key] = domain
	}
	for key, subdomain := range c.Subdomains {
		if err := subdomain.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.Subdomains[key] = subdomain
	}
	for key, datastore := range c.Datastores {
		if err := datastore.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.Datastores[key] = datastore
	}
	for key, app := range c.Applications {
		if err := app.Normalize(ctx, resolver); err != nil {
			return err
		}
		c.Applications[key] = app
	}
	return nil
}
