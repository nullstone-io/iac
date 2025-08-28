package config

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
)

var (
	randomVarRegexp     = regexp.MustCompile(`\{\{\s*random\(\)\s*}}`)
	onlyRandomVarRegexp = regexp.MustCompile(`^\s*\{\{\s*random\(\)\s*}}\s*$`)
)

type SubdomainConfiguration struct {
	BlockConfiguration

	DomainNameTemplate    *string                     `json:"domainNameTemplate"`
	SubdomainNameTemplate *string                     `json:"subdomainNameTemplate"`
	Reservation           *types.SubdomainReservation `json:"reservation,omitempty"`
}

func (s *SubdomainConfiguration) ToBlock(orgName string, stackId int64) types.Block {
	block := s.BlockConfiguration.ToBlock(orgName, stackId)
	return block
}

func (s *SubdomainConfiguration) Resolve(ctx context.Context, resolver core.ResolveResolver, ic core.IacContext, pc core.ObjectPathContext) core.ResolveErrors {
	errs := s.BlockConfiguration.Resolve(ctx, resolver, ic, pc)
	errs = append(errs, s.reserve(ctx, resolver, pc)...)
	s.resolveDomain()
	return errs
}

// reserve requests a subdomain from Nullstone
// This is only performed if there is no "domain"/"subdomain" connection in the module definition
// - If `template="{{ random() }}"`, we request a random subdomain
// - Else, we ask for the requested subdomain in SubdomainNameTemplate
func (s *SubdomainConfiguration) reserve(ctx context.Context, resolver core.ResolveResolver, pc core.ObjectPathContext) core.ResolveErrors {
	if s.ModuleVersion == nil {
		return nil
	}
	conns := s.ModuleVersion.Manifest.Connections
	_, hasDomain := conns["domain"]
	_, hasSubdomain := conns["subdomain"]
	if hasDomain || hasSubdomain {
		return nil
	}

	isRandom, errs := s.detectRandom(pc)
	if len(errs) > 0 {
		return errs
	}
	if s.SubdomainNameTemplate == nil {
		return nil
	}

	requested := *s.SubdomainNameTemplate
	if isRandom {
		requested = "random()"
	}
	reservation, err := resolver.ReserveNullstoneSubdomain(ctx, s.BlockConfiguration.Name, requested)
	if err != nil {
		return core.ResolveErrors{core.FailedSubdomainReservationError(pc, requested, err)}
	}
	s.Reservation = reservation
	return nil
}

func (s *SubdomainConfiguration) detectRandom(pc core.ObjectPathContext) (bool, core.ResolveErrors) {
	if s.SubdomainNameTemplate == nil {
		return false, nil
	}
	template := *s.SubdomainNameTemplate
	if !randomVarRegexp.MatchString(template) {
		// If {{ random() }} is not in the subdomain template, we're not going to attempt to reserve a random
		return false, nil
	}
	if !onlyRandomVarRegexp.MatchString(template) {
		// {{ random() }} must be used as a standalone subdomain template
		return true, core.ResolveErrors{core.InvalidRandomSubdomainTemplateError(pc.SubField("dns").SubField("template"), template)}
	}
	return true, nil
}

func (s *SubdomainConfiguration) resolveDomain() {
	conn, ok := s.Connections["domain"]
	if !ok || conn.Block == nil {
		return
	}
	dnsName := conn.Block.DnsName
	s.DomainNameTemplate = &dnsName
	return
}

func (s *SubdomainConfiguration) ApplyChangesTo(ic core.IacContext, updater core.WorkspaceConfigUpdater) error {
	if err := s.BlockConfiguration.ApplyChangesTo(ic, updater); err != nil {
		return err
	}
	updater.UpdateSubdomainName(s.DomainNameTemplate, s.SubdomainNameTemplate, s.Reservation)
	return nil
}

func convertSubdomainConfigurations(parsed map[string]yaml.SubdomainConfiguration) map[string]*SubdomainConfiguration {
	result := make(map[string]*SubdomainConfiguration)
	for name, value := range parsed {
		bc := blockConfigFromYaml(name, value.BlockConfiguration, BlockTypeSubdomain, types.CategorySubdomain)
		subdomainNameTemplate := value.Dns.Template
		if subdomainNameTemplate == "" {
			// If a user is using old syntax, we're going to add `.{{ NULLSTONE_ENV }}` suffix to the template
			subdomainNameTemplate = fmt.Sprintf("%s.{{ NULLSTONE_ENV }}", value.DnsName)
		}
		result[name] = &SubdomainConfiguration{
			BlockConfiguration:    *bc,
			SubdomainNameTemplate: &subdomainNameTemplate,
		}
	}
	return result
}
