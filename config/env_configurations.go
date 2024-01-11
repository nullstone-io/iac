package config

import (
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/yaml.v3"
)

type EnvConfiguration struct {
	Version           string                                   `yaml:"version" json:"version"`
	Applications      map[string]AppConfiguration              `yaml:"apps,omitempty" json:"apps"`
	Subdomains        map[string]SubdomainConfiguration        `yaml:"subdomains,omitempty" json:"subdomains"`
	Datastores        map[string]DatastoreConfiguration        `yaml:"datastores,omitempty" json:"datastores"`
	Domains           map[string]DomainConfiguration           `yaml:"domains,omitempty" json:"domains"`
	Ingresses         map[string]IngressConfiguration          `yaml:"ingresses,omitempty" json:"ingresses"`
	Networks          map[string]NetworkConfiguration          `yaml:"networks,omitempty" json:"networks"`
	Clusters          map[string]ClusterConfiguration          `yaml:"clusters,omitempty" json:"clusters"`
	ClusterNamespaces map[string]ClusterNamespaceConfiguration `yaml:"cluster_namespaces,omitempty" json:"clusterNamespaces"`
	Blocks            map[string]BlockConfiguration            `yaml:"blocks,omitempty" json:"blocks"`
}

func ParseEnvConfiguration(data []byte) (*EnvConfiguration, error) {
	var r *EnvConfiguration
	err := yaml.Unmarshal(data, &r)
	if err != nil {
		return nil, core.InvalidYamlError("config.yml", err)
	}

	newApps := make(map[string]AppConfiguration)
	for appName, appValue := range r.Applications {
		appValue.Name = appName
		// set a default module version if not provided
		if appValue.ModuleSourceVersion == nil {
			latest := "latest"
			appValue.ModuleSourceVersion = &latest
		}
		newCaps := make([]core.CapabilityConfiguration, len(appValue.Capabilities))
		for i, capValue := range appValue.Capabilities {
			// set a default module version if not provided
			if capValue.ModuleSourceVersion == nil {
				latest := "latest"
				capValue.ModuleSourceVersion = &latest
			}
			newCaps[i] = capValue
		}
		appValue.Capabilities = newCaps
		newApps[appName] = appValue
	}
	r.Applications = newApps

	newSubdomains := make(map[string]SubdomainConfiguration)
	for subdomainName, subdomainValue := range r.Subdomains {
		subdomainValue.Name = subdomainName
		// set a default module version if not provided
		if subdomainValue.ModuleSourceVersion == nil {
			latest := "latest"
			subdomainValue.ModuleSourceVersion = &latest
		}
		newSubdomains[subdomainName] = subdomainValue
	}
	r.Subdomains = newSubdomains

	newDatastores := make(map[string]DatastoreConfiguration)
	for datastoreName, datastoreValue := range r.Datastores {
		datastoreValue.Name = datastoreName
		// set a default module version if not provided
		if datastoreValue.ModuleSourceVersion == nil {
			latest := "latest"
			datastoreValue.ModuleSourceVersion = &latest
		}
		newDatastores[datastoreName] = datastoreValue
	}
	r.Datastores = newDatastores

	newDomains := make(map[string]DomainConfiguration)
	for domainName, domainValue := range r.Domains {
		domainValue.Name = domainName
		// set a default module version if not provided
		if domainValue.ModuleSourceVersion == nil {
			latest := "latest"
			domainValue.ModuleSourceVersion = &latest
		}
		newDomains[domainName] = domainValue
	}
	r.Domains = newDomains

	newIngresses := make(map[string]IngressConfiguration)
	for ingressName, ingressValue := range r.Ingresses {
		ingressValue.Name = ingressName
		// set a default module version if not provided
		if ingressValue.ModuleSourceVersion == nil {
			latest := "latest"
			ingressValue.ModuleSourceVersion = &latest
		}
		newIngresses[ingressName] = ingressValue
	}
	r.Ingresses = newIngresses

	newClusterNamespaces := make(map[string]ClusterNamespaceConfiguration)
	for clusterNamespaceName, clusterNamespaceValue := range r.ClusterNamespaces {
		clusterNamespaceValue.Name = clusterNamespaceName
		// set a default module version if not provided
		if clusterNamespaceValue.ModuleSourceVersion == nil {
			latest := "latest"
			clusterNamespaceValue.ModuleSourceVersion = &latest
		}
		newClusterNamespaces[clusterNamespaceName] = clusterNamespaceValue
	}
	r.ClusterNamespaces = newClusterNamespaces

	newClusters := make(map[string]ClusterConfiguration)
	for clusterName, clusterValue := range r.Clusters {
		clusterValue.Name = clusterName
		// set a default module version if not provided
		if clusterValue.ModuleSourceVersion == nil {
			latest := "latest"
			clusterValue.ModuleSourceVersion = &latest
		}
		newClusters[clusterName] = clusterValue
	}
	r.Clusters = newClusters

	newNetworks := make(map[string]NetworkConfiguration)
	for networkName, networkValue := range r.Networks {
		networkValue.Name = networkName
		// set a default module version if not provided
		if networkValue.ModuleSourceVersion == nil {
			latest := "latest"
			networkValue.ModuleSourceVersion = &latest
		}
		newNetworks[networkName] = networkValue
	}
	r.Networks = newNetworks

	newBlocks := make(map[string]BlockConfiguration)
	for blockName, blockValue := range r.Blocks {
		blockValue.Name = blockName
		// set a default module version if not provided
		if blockValue.ModuleSourceVersion == nil {
			latest := "latest"
			blockValue.ModuleSourceVersion = &latest
		}
		newBlocks[blockName] = blockValue
	}
	r.Blocks = newBlocks

	return r, nil
}

func (e EnvConfiguration) Validate(resolver *find.ResourceResolver) (errors.ValidationErrors, error) {
	ve := errors.ValidationErrors{}
	for _, block := range e.Blocks {
		verrs, err := block.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, network := range e.Networks {
		verrs, err := network.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, cluster := range e.Clusters {
		verrs, err := cluster.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, clusterNamespace := range e.ClusterNamespaces {
		verrs, err := clusterNamespace.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, ingress := range e.Ingresses {
		verrs, err := ingress.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, domain := range e.Domains {
		verrs, err := domain.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, subdomain := range e.Subdomains {
		verrs, err := subdomain.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, datastore := range e.Datastores {
		verrs, err := datastore.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	for _, app := range e.Applications {
		verrs, err := app.Validate(resolver, e.blocksFromConfig())
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	return ve, nil
}

func (e EnvConfiguration) blocksFromConfig() []core.BlockConfiguration {
	result := make([]core.BlockConfiguration, 0)
	for name, block := range e.Blocks {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: block.ModuleSource, ModuleSourceVersion: block.ModuleSourceVersion})
	}
	for name, network := range e.Networks {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: network.ModuleSource, ModuleSourceVersion: network.ModuleSourceVersion})
	}
	for name, cluster := range e.Clusters {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: cluster.ModuleSource, ModuleSourceVersion: cluster.ModuleSourceVersion})
	}
	for name, clusterNamespace := range e.ClusterNamespaces {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: clusterNamespace.ModuleSource, ModuleSourceVersion: clusterNamespace.ModuleSourceVersion})
	}
	for name, ingress := range e.Ingresses {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: ingress.ModuleSource, ModuleSourceVersion: ingress.ModuleSourceVersion})
	}
	for name, domain := range e.Domains {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: domain.ModuleSource, ModuleSourceVersion: domain.ModuleSourceVersion})
	}
	for name, sub := range e.Subdomains {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: sub.ModuleSource, ModuleSourceVersion: sub.ModuleSourceVersion})
	}
	for name, ds := range e.Datastores {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: ds.ModuleSource, ModuleSourceVersion: ds.ModuleSourceVersion})
	}
	for name, app := range e.Applications {
		result = append(result, core.BlockConfiguration{Name: name, ModuleSource: app.ModuleSource, ModuleSourceVersion: app.ModuleSourceVersion})
	}
	return result
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
