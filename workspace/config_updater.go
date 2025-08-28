package workspace

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	_ core.WorkspaceConfigUpdater  = ConfigUpdater{}
	_ core.CapabilityConfigUpdater = CapabilityConfigUpdater{}

	DefaultModuleConstraint = "latest"
)

type ConfigUpdater struct {
	Config *types.WorkspaceConfig
	TemplateVars
}

func (w ConfigUpdater) UpdateSchema(moduleSource, moduleConstraint string, moduleVersion *types.ModuleVersion) {
	if moduleVersion == nil {
		return
	}
	if moduleSource != "" {
		w.Config.Source = moduleSource
	}
	w.Config.SourceConstraint = DefaultModuleConstraint
	if moduleConstraint != "" {
		w.Config.SourceConstraint = moduleConstraint
	}
	w.Config.SourceVersion = moduleVersion.Version
	w.Config.SourceToolName = moduleVersion.ToolName
	if w.Config.Variables == nil {
		w.Config.Variables = types.Variables{}
	}
	if w.Config.Connections == nil {
		w.Config.Connections = types.Connections{}
	}
	ModuleSchema(moduleVersion.Manifest).UpdateSchema(w.Config.Variables, w.Config.Connections)
}

func (w ConfigUpdater) UpdateVariableValue(name string, value any) {
	existing, ok := w.Config.Variables[name]
	if !ok {
		return
	}
	// if the existing value is nil, it means it hasn't been set and we are using the default value
	// if someone tries to set the value to the default value, we don't need to do anything
	// this allows us to avoid some very challenging scenarios in workspace_changes
	// without this, variables will show up as "changed" rows in the UI but the values are the same
	if existing.Value == nil && reflect.DeepEqual(value, existing.Variable.Default) {
		return
	}
	existing.Value = value
	w.Config.Variables[name] = existing
}

func (w ConfigUpdater) UpdateConnectionTarget(name string, desired, effective types.ConnectionTarget) {
	existing, ok := w.Config.Connections[name]
	if !ok {
		return
	}
	existing.DesiredTarget = &desired
	existing.EffectiveTarget = &effective
	existing.OldReference = &effective
	w.Config.Connections[name] = existing
}

func (w ConfigUpdater) AddOrUpdateEnvVariable(name string, value string, sensitive bool) {
	envVar, ok := w.Config.EnvVariables[name]
	// if we find the env variable, just update the value and sensitive flag
	if ok {
		envVar.Value = value
		envVar.Sensitive = sensitive
	} else {
		// otherwise, create a new env variable
		envVar = types.EnvVariable{
			Value:     value,
			Sensitive: sensitive,
		}
	}
	if w.Config.EnvVariables == nil {
		w.Config.EnvVariables = types.EnvVariables{}
	}
	w.Config.EnvVariables[name] = envVar
}

func (w ConfigUpdater) RemoveEnvVariablesNotIn(envVariables map[string]string) {
	for k, _ := range w.Config.EnvVariables {
		if _, ok := envVariables[k]; !ok {
			delete(w.Config.EnvVariables, k)
		}
	}
}

func (w ConfigUpdater) GetCapabilityUpdater(identity core.CapabilityIdentity) core.CapabilityConfigUpdater {
	for i, cur := range w.Config.Capabilities {
		found := identity.Match(core.CapabilityIdentity{
			Name:              cur.Name,
			ModuleSource:      cur.Source,
			ConnectionTargets: cur.Connections.DesiredTargets(),
		})
		if found {
			return CapabilityConfigUpdater{
				WorkspaceConfig: w.Config,
				Index:           i,
			}
		}
	}
	return nil
}

func (w ConfigUpdater) AddCapability(id int64, name string) core.CapabilityConfigUpdater {
	w.Config.Capabilities = append(w.Config.Capabilities, types.CapabilityConfig{Id: id, Name: name})
	ccu := CapabilityConfigUpdater{
		WorkspaceConfig: w.Config,
		Index:           len(w.Config.Capabilities) - 1,
	}
	return ccu
}

func (w ConfigUpdater) RemoveCapabilitiesNotIn(identities core.CapabilityIdentities) {
	result := make(types.CapabilityConfigs, 0)
	for _, cur := range w.Config.Capabilities {
		found := identities.Find(core.CapabilityIdentity{
			Name:              cur.Name,
			ModuleSource:      cur.Source,
			ConnectionTargets: cur.Connections.DesiredTargets(),
		})
		if found != nil {
			// If we found the capability in the IaC file, let's keep it
			result = append(result, cur)
		}
	}
	w.Config.Capabilities = result
}

func (w ConfigUpdater) UpdateDomainName(domainNameTemplate *string) {
	extra := w.Config.Extra

	if domainNameTemplate != nil {
		domainName := w.TemplateVars.ReplaceSpecificVars(*domainNameTemplate, "NULLSTONE_ORG")
		extra.Domain = &types.ExtraDomainConfig{
			DomainNameTemplate: *domainNameTemplate,
			DomainName:         domainName,
			Fqdn:               fmt.Sprintf("%s.", extra.Domain.DomainName),
		}
	}

	w.Config.Extra = extra
}

func (w ConfigUpdater) UpdateSubdomainName(domainNameTemplate, subdomainNameTemplate *string, reservation *types.SubdomainReservation) {
	extra := w.Config.Extra

	if extra.Subdomain == nil {
		extra.Subdomain = &types.ExtraSubdomainConfig{}
	}

	if domainNameTemplate != nil {
		extra.Subdomain.DomainName = w.TemplateVars.ReplaceSpecificVars(*domainNameTemplate, "NULLSTONE_ORG")
	} else if reservation != nil {
		extra.Subdomain.DomainName = reservation.DomainName
	}

	if subdomainNameTemplate != nil {
		extra.Subdomain.SubdomainNameTemplate = *subdomainNameTemplate
		if reservation != nil {
			extra.Subdomain.SubdomainName = reservation.SubdomainName
		} else {
			extra.Subdomain.SubdomainName = strings.TrimSuffix(w.TemplateVars.ReplaceVars(*subdomainNameTemplate), ".")
		}
	}

	if domainNameTemplate != nil || subdomainNameTemplate != nil || reservation != nil {
		extra.Subdomain.Fqdn = extra.Subdomain.CalculateFqdn()
	}

	w.Config.Extra = extra
}

type CapabilityConfigUpdater struct {
	WorkspaceConfig *types.WorkspaceConfig
	Index           int
}

func (c CapabilityConfigUpdater) UpdateSchema(moduleSource, moduleConstraint string, moduleVersion *types.ModuleVersion) {
	c.doOperation(func(cc *types.CapabilityConfig) {
		if moduleVersion == nil {
			return
		}
		if moduleSource != "" {
			cc.Source = moduleSource
		}
		cc.SourceConstraint = DefaultModuleConstraint
		if moduleConstraint != "" {
			cc.SourceConstraint = moduleConstraint
		}
		cc.SourceVersion = moduleVersion.Version
		if cc.Variables == nil {
			cc.Variables = types.Variables{}
		}
		if cc.Connections == nil {
			cc.Connections = types.Connections{}
		}
		ModuleSchema(moduleVersion.Manifest).UpdateSchema(cc.Variables, cc.Connections)
	})
}

func (c CapabilityConfigUpdater) UpdateVariableValue(name string, value any) {
	c.doOperation(func(cc *types.CapabilityConfig) {
		existingVar, ok := cc.Variables[name]
		if !ok {
			return
		}
		existingVar.Value = value
		cc.Variables[name] = existingVar
	})
}

func (c CapabilityConfigUpdater) UpdateConnectionTarget(name string, desired, effective types.ConnectionTarget) {
	c.doOperation(func(cc *types.CapabilityConfig) {
		existingConn, ok := cc.Connections[name]
		if !ok {
			return
		}
		existingConn.DesiredTarget = &desired
		existingConn.EffectiveTarget = &effective
		existingConn.OldReference = &effective
		cc.Connections[name] = existingConn
	})
}

func (c CapabilityConfigUpdater) UpdateNamespace(namespace *string) {
	c.doOperation(func(cc *types.CapabilityConfig) {
		if namespace != nil {
			cc.Namespace = *namespace
		}
	})
}

func (c CapabilityConfigUpdater) doOperation(fn func(cc *types.CapabilityConfig)) {
	capConfig := c.WorkspaceConfig.Capabilities[c.Index]
	fn(&capConfig)
	c.WorkspaceConfig.Capabilities[c.Index] = capConfig
}
