package workspace

import (
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"reflect"
)

var (
	_ core.WorkspaceConfigUpdater  = ConfigUpdater{}
	_ core.CapabilityConfigUpdater = CapabilityConfigUpdater{}
)

type ConfigUpdater struct {
	Config *types.WorkspaceConfig
}

func (w ConfigUpdater) UpdateSchema(moduleSource string, moduleVersion *types.ModuleVersion) {
	if moduleVersion == nil {
		return
	}
	if moduleSource != "" {
		w.Config.Source = moduleSource
	}
	w.Config.SourceVersion = moduleVersion.Version
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

func (w ConfigUpdater) UpdateConnectionTarget(name string, value types.ConnectionTarget) {
	existing, ok := w.Config.Connections[name]
	if !ok {
		return
	}
	existing.Reference = &value
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
	w.Config.EnvVariables[name] = envVar
}

func (w ConfigUpdater) RemoveEnvVariablesNotIn(envVariables map[string]string) {
	for k, _ := range w.Config.Variables {
		if _, ok := envVariables[k]; !ok {
			delete(w.Config.Variables, k)
		}
	}
}

func (w ConfigUpdater) GetCapabilityUpdater(identity core.CapabilityIdentity) core.CapabilityConfigUpdater {
	for i, cur := range w.Config.Capabilities {
		found := identity.Match(core.CapabilityIdentity{
			ModuleSource:      cur.Source,
			ConnectionTargets: cur.Connections.Targets(),
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

func (w ConfigUpdater) RemoveCapabilitiesNotIn(identities core.CapabilityIdentities) {
	result := make(types.CapabilityConfigs, 0)
	for _, cur := range w.Config.Capabilities {
		found := identities.Find(core.CapabilityIdentity{
			ModuleSource:      cur.Source,
			ConnectionTargets: cur.Connections.Targets(),
		})
		if found != nil {
			result = append(result, cur)
		}
	}
	w.Config.Capabilities = result
}

type CapabilityConfigUpdater struct {
	WorkspaceConfig *types.WorkspaceConfig
	Index           int
}

func (c CapabilityConfigUpdater) UpdateSchema(moduleSource string, moduleVersion *types.ModuleVersion) {
	c.doOperation(func(cc *types.CapabilityConfig) {
		if moduleVersion == nil {
			return
		}
		if moduleSource != "" {
			cc.Source = moduleSource
		}
		cc.SourceVersion = moduleVersion.Version
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

func (c CapabilityConfigUpdater) UpdateConnectionTarget(name string, value types.ConnectionTarget) {
	c.doOperation(func(cc *types.CapabilityConfig) {
		existingConn, ok := cc.Connections[name]
		if !ok {
			return
		}
		existingConn.Reference = &value
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
