package services

import (
	"fmt"
	"os"
)

const (
	ServiceDomainEnvVar = "SERVICE_DOMAIN"
)

func QualifySubdomain(service string) string {
	if domain := os.Getenv(ServiceDomainEnvVar); domain != "" {
		return fmt.Sprintf("%s.%s", service, domain)
	}
	return service
}
