package config

import (
	errs "errors"
	"fmt"
	"github.com/BSick7/go-api/errors"
	"gopkg.in/nullstone-io/go-api-client.v0/find"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

// ValidateCapabilities performs validation on a all IaC capabilities within an application
func ValidateCapabilities(resolver *find.ResourceResolver, configBlocks []BlockConfiguration, path string, capabilities CapabilityConfigurations, subcategory types.SubcategoryName) error {
	ve := errors.ValidationErrors{}
	for i, iacCap := range capabilities {
		capPath := fmt.Sprintf("%s.capabilities[%d]", path, i)
		err := ValidateCapability(resolver, configBlocks, capPath, iacCap, string(subcategory))
		if err != nil {
			var verrs errors.ValidationErrors
			if errs.As(err, &verrs) {
				ve = append(ve, verrs...)
			}
		}
	}

	if len(ve) > 0 {
		return ve
	}
	return nil
}
