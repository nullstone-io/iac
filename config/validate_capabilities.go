package config

import (
	"fmt"
	"github.com/BSick7/go-api/errors"
	"github.com/nullstone-io/iac/core"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// ValidateCapabilities performs validation on a all IaC capabilities within an application
func ValidateCapabilities(resolver *find.ResourceResolver, path string, capabilities CapabilityConfigurations, subcategory types.SubcategoryName) (errors.ValidationErrors, error) {
	ve := errors.ValidationErrors{}
	for i, iacCap := range capabilities {
		capPath := fmt.Sprintf("%s.capabilities[%d]", path, i)
		verrs, err := core.ValidateCapability(resolver, capPath, iacCap, string(subcategory))
		if err != nil {
			return ve, err
		}
		ve = append(ve, verrs...)
	}
	return ve, nil
}
