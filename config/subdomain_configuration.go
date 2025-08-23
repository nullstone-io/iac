package config

import (
	"context"
	"github.com/nullstone-io/iac/core"
	"github.com/nullstone-io/iac/yaml"
	"gopkg.in/nullstone-io/go-api-client.v0/types"
	"regexp"
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
	errs = append(errs, s.reserveRandom(ctx, resolver, pc)...)
	s.resolveDomain(ctx, resolver)
	return errs
}

func (s *SubdomainConfiguration) reserveRandom(ctx context.Context, resolver core.ResolveResolver, pc core.ObjectPathContext) core.ResolveErrors {
	if s.SubdomainNameTemplate != nil {
		template := *s.SubdomainNameTemplate
		if !randomVarRegexp.MatchString(template) {
			// If {{ random() }} is not in the subdomain template, we're not going to attempt to reserve a random
			return nil
		}
		if !onlyRandomVarRegexp.MatchString(template) {
			// {{ random() }} must be used as a standalone subdomain template
			return core.ResolveErrors{core.InvalidRandomSubdomainTemplateError(pc.SubField("dns").SubField("template"), template)}
		}
	}

	reservation, err := resolver.ReserveNullstoneSubdomain(ctx, s.BlockConfiguration.Name, "random()")
	if err != nil {
		return core.ResolveErrors{core.FailedSubdomainReservationError(pc, "random()", err)}
	}
	s.Reservation = reservation
	return nil
}

func (s *SubdomainConfiguration) resolveDomain(ctx context.Context, resolver core.ResolveResolver) {
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
			subdomainNameTemplate = value.DnsName
		}
		result[name] = &SubdomainConfiguration{BlockConfiguration: *bc, SubdomainNameTemplate: &subdomainNameTemplate}
	}
	return result
}
